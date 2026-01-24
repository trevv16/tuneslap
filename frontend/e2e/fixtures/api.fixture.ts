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
        await page.route('**/api/v1/users/me', async (route) => {
          await fulfillJson(route, user, options)
        })
      },

      /**
       * Mock /me to return unauthorized.
       */
      meUnauthorized: async () => {
        await page.route('**/api/v1/users/me', async (route) => {
          await fulfillJson(route, mockUnauthorizedResponse, { status: 401 })
        })
      },
    },

    boards: {
      /**
       * Mock the boards list endpoint.
       */
      list: async (boards = mockBoards, options?: MockResponseOptions) => {
        await page.route('**/api/v1/boards', async (route) => {
          if (route.request().method() === 'GET') {
            await fulfillJson(route, boards, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock empty boards list.
       */
      listEmpty: async (options?: MockResponseOptions) => {
        await page.route('**/api/v1/boards', async (route) => {
          if (route.request().method() === 'GET') {
            await fulfillJson(route, [], options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock a single board by ID.
       */
      get: async (
        id: string,
        board = mockBoard,
        options?: MockResponseOptions
      ) => {
        await page.route(`**/api/v1/boards/${id}`, async (route) => {
          if (route.request().method() === 'GET') {
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
        await page.route(`**/api/v1/boards/${id}`, async (route) => {
          if (route.request().method() === 'GET') {
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
        response = mockBoard,
        options?: MockResponseOptions
      ) => {
        await page.route('**/api/v1/boards', async (route) => {
          if (route.request().method() === 'POST') {
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
        response = mockBoard,
        options?: MockResponseOptions
      ) => {
        await page.route(`**/api/v1/boards/${id}`, async (route) => {
          if (route.request().method() === 'PATCH') {
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
        await page.route(`**/api/v1/boards/${id}`, async (route) => {
          if (route.request().method() === 'DELETE') {
            await fulfillJson(route, { success: true }, options)
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
        await page.route('**/api/v1/media**', async (route) => {
          if (route.request().method() === 'GET') {
            await fulfillJson(route, media, options)
          } else {
            await route.continue()
          }
        })
      },

      /**
       * Mock empty media list.
       */
      listEmpty: async (options?: MockResponseOptions) => {
        await page.route('**/api/v1/media**', async (route) => {
          if (route.request().method() === 'GET') {
            await fulfillJson(route, [], options)
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
        await page.route('**/api/v1/media**', async (route) => {
          if (route.request().method() === 'GET') {
            await fulfillJson(route, filtered, options)
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
        response: unknown,
        options?: MockResponseOptions
      ) => {
        await page.route(`**/api/v1/boards/${boardId}/keys`, async (route) => {
          if (route.request().method() === 'POST') {
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
      await page.route('**/api/v1/**', async (route) => {
        await fulfillJson(
          route,
          { ...mockErrorResponse, message },
          { status: 500 }
        )
      })
    },

    /**
     * Mock all API calls to return unauthorized.
     */
    allUnauthorized: async () => {
      await page.route('**/api/v1/**', async (route) => {
        await fulfillJson(route, mockUnauthorizedResponse, { status: 401 })
      })
    },
  }
}

export type ApiMocks = ReturnType<typeof createApiMocks>
