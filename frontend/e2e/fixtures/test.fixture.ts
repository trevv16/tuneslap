import { test as base, expect, type Page, type BrowserContext } from '@playwright/test'
import { createApiMocks, type ApiMocks } from './api.fixture'
import { setAuthToken, STORAGE_STATE } from './auth.fixture'
import { mockUserResponse } from './test-data'

/**
 * Extended test fixtures combining all helpers.
 */
export const test = base.extend<{
  // API mocking utilities
  apiMocks: ApiMocks
  // An authenticated page with mocked APIs
  authedPage: Page
}>({
  // Provide API mocking utilities for the current page
  apiMocks: async ({ page }, use) => {
    const mocks = createApiMocks(page)
    await use(mocks)
  },

  // Provide an authenticated page with common API mocks set up
  authedPage: async ({ browser }, use) => {
    // Create a new context with stored auth state if available
    let context: BrowserContext

    try {
      context = await browser.newContext({
        storageState: STORAGE_STATE,
      })
    } catch {
      // If storage state doesn't exist, create a fresh context
      context = await browser.newContext()
    }

    const page = await context.newPage()

    // Set up auth token
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const mocks = createApiMocks(page)
    await mocks.auth.me(mockUserResponse)

    await use(page)
    await context.close()
  },
})

// Re-export expect for convenience
export { expect }

// Re-export test data
export * from './test-data'

// Re-export auth helpers
export { login, loginWithMock, signup, logout, isAuthenticated, setAuthToken } from './auth.fixture'

// Re-export API mock creator
export { createApiMocks, type ApiMocks } from './api.fixture'
