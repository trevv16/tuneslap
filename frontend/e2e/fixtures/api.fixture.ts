import type { Page, Route } from '@playwright/test'
import {
  mockBoards,
  mockBoard,
  mockEmptyBoard,
  mockMedia,
  mockSigninResponse,
  mockSignupResponse,
  mockUserResponse,
  mockErrorResponse,
  mockUnauthorizedResponse,
} from './test-data'

type MockResponseOptions = {
  status?: number
  delay?: number
}

/**
 * Helper to create a JSON response for route fulfillment.
 */
async function fulfillJson(
  route: Route,
  data: unknown,
  options: MockResponseOptions = {}
): Promise<void> {
  const { status = 200, delay = 0 } = options

  if (delay > 0) {
    await new Promise((resolve) => setTimeout(resolve, delay))
  }

  await route.fulfill({
    status,
    contentType: 'application/json',
    body: JSON.stringify(data),
  })
}

/**
 * Create API mocking utilities for a page.
 * Uses test data from fixtures as default responses.
 */
export function createApiMocks(page: Page) {
  return {
    auth: {
      /**
       * Mock the signin endpoint.
       */
      signin: async (
        response = mockSigninResponse,
        options?: MockResponseOptions
      ) => {
        await page.route('**/api/v1/auth/signin', async (route) => {
          await fulfillJson(route, response, options)
        })
      },

      /**
       * Mock signin to return an error.
       */
      signinError: async (
        message = 'Invalid credentials',
        options?: MockResponseOptions
      ) => {
        await page.route('**/api/v1/auth/signin', async (route) => {
          await fulfillJson(
            route,
            { success: false, message, error: 'AUTH_ERROR' },
            { status: 401, ...options }
          )
        })
      },

      /**
       * Mock the signup endpoint.
       */
      signup: async (
        response = mockSignupResponse,
        options?: MockResponseOptions
      ) => {
        await page.route('**/api/v1/auth/signup', async (route) => {
          await fulfillJson(route, response, options)
        })
      },

      /**
       * Mock the /me endpoint for current user.
       */
      me: async (user = mockUserResponse, options?: MockResponseOptions) => {
        await page.route('**/users/me', async (route) => {
          // Wrap user in the expected API response format
          const response = {
            success: true,
            message: 'Success',
            data: user,
          }
          await fulfillJson(route, response, options)
        })
      },

      /**
       * Mock /me to return unauthorized.
       */
      meUnauthorized: async () => {
        await page.route('**/users/me', async (route) => {
          await fulfillJson(route, mockUnauthorizedResponse, { status: 401 })
        })
      },
    },

    boards: {
      /**
       * Mock the boards list endpoint.
       * Uses pattern to only match API calls, not page navigation.
       */
      list: async (boards = mockBoards, options?: MockResponseOptions) => {
        await page.route('**/api/**/boards', async (route) => {
          if (route.request().method() === 'GET') {
            // Wrap in API response format
            const response = {
              success: true,
              message: 'Success',
              data: {
                boards,
                pagination: { page: 1, limit: 20, total: boards.length, totalPages: 1 },
              },
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock empty boards list.
       */
      listEmpty: async (options?: MockResponseOptions) => {
        await page.route('**/api/**/boards', async (route) => {
          if (route.request().method() === 'GET') {
            const response = {
              success: true,
              message: 'Success',
              data: {
                boards: [],
                pagination: { page: 1, limit: 20, total: 0, totalPages: 0 },
              },
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock a single board by ID.
       * Returns BoardResponse directly (not wrapped).
       */
      get: async (
        id: string,
        board = mockBoard,
        options?: MockResponseOptions
      ) => {
        await page.route(`**/boards/${id}`, async (route) => {
          const url = route.request().url()
          // Only mock API calls (contain /api/ or api.tuneslap.com), not page navigation
          if (route.request().method() === 'GET' && (url.includes('/api/') || url.includes('api.tuneslap.com'))) {
            // Return board directly - API client expects BoardResponse, not wrapped
            await fulfillJson(route, board, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock getting an empty board (no keys).
       */
      getEmpty: async (id: string, options?: MockResponseOptions) => {
        await page.route(`**/boards/${id}`, async (route) => {
          const url = route.request().url()
          // Only mock API calls, not page navigation
          if (route.request().method() === 'GET' && (url.includes('/api/') || url.includes('api.tuneslap.com'))) {
            // Return board directly - API client expects BoardResponse, not wrapped
            await fulfillJson(route, { ...mockEmptyBoard, id }, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock board creation.
       */
      create: async (
        board = mockBoard,
        options?: MockResponseOptions
      ) => {
        await page.route('**/api/**/boards', async (route) => {
          if (route.request().method() === 'POST') {
            const response = {
              success: true,
              message: 'Board created successfully',
              data: board,
            }
            await fulfillJson(route, response, { status: 201, ...options })
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock board update.
       */
      update: async (
        id: string,
        board = mockBoard,
        options?: MockResponseOptions
      ) => {
        await page.route(`**/api/**/boards/${id}`, async (route) => {
          if (route.request().method() === 'PATCH') {
            const response = {
              success: true,
              message: 'Board updated successfully',
              data: board,
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock board deletion.
       */
      delete: async (id: string, options?: MockResponseOptions) => {
        await page.route(`**/api/**/boards/${id}`, async (route) => {
          if (route.request().method() === 'DELETE') {
            const response = {
              success: true,
              message: 'Board deleted successfully',
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },
    },

    media: {
      /**
       * Mock the media list endpoint.
       */
      list: async (media = mockMedia, options?: MockResponseOptions) => {
        await page.route('**/media**', async (route) => {
          if (route.request().method() === 'GET') {
            const response = {
              success: true,
              message: 'Success',
              data: {
                media,
                pagination: { page: 1, limit: 20, total: media.length, totalPages: 1 },
              },
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock empty media list.
       */
      listEmpty: async (options?: MockResponseOptions) => {
        await page.route('**/media**', async (route) => {
          if (route.request().method() === 'GET') {
            const response = {
              success: true,
              message: 'Success',
              data: {
                media: [],
                pagination: { page: 1, limit: 20, total: 0, totalPages: 0 },
              },
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock media by type filter.
       */
      listByType: async (
        type: 'audio' | 'image',
        options?: MockResponseOptions
      ) => {
        const filtered = mockMedia.filter((m) => m.mediaType === type)
        await page.route('**/media**', async (route) => {
          if (route.request().method() === 'GET') {
            const response = {
              success: true,
              message: 'Success',
              data: {
                media: filtered,
                pagination: { page: 1, limit: 20, total: filtered.length, totalPages: 1 },
              },
            }
            await fulfillJson(route, response, options)
          } else {
            await route.continue()
          }
        })
      },
    },

    keys: {
      /**
       * Mock key creation.
       */
      create: async (
        boardId: string,
        key: unknown,
        options?: MockResponseOptions
      ) => {
        await page.route(`**/api/**/boards/${boardId}/keys`, async (route) => {
          if (route.request().method() === 'POST') {
            const response = {
              success: true,
              message: 'Key created successfully',
              data: key,
            }
            await fulfillJson(route, response, { status: 201, ...options })
          } else {
            await route.continue()
          }
        })
      },
    },

    /**
     * Mock all API calls to return an error.
     */
    allError: async (message = 'Server error') => {
      await page.route('**/*', async (route) => {
        const url = route.request().url()
        if (url.includes('/api/') || url.includes('/users/') || url.includes('/boards/') || url.includes('/media/')) {
          await fulfillJson(
            route,
            { ...mockErrorResponse, message },
            { status: 500 }
          )
        } else {
          await route.continue()
        }
      })
    },

    /**
     * Mock all API calls to return unauthorized.
     */
    allUnauthorized: async () => {
      await page.route('**/*', async (route) => {
        const url = route.request().url()
        if (url.includes('/api/') || url.includes('/users/') || url.includes('/boards/') || url.includes('/media/')) {
          await fulfillJson(route, mockUnauthorizedResponse, { status: 401 })
        } else {
          await route.continue()
        }
      })
    },
  }
}

export type ApiMocks = ReturnType<typeof createApiMocks>
