import { test as base, expect, type Page, type BrowserContext } from '@playwright/test'
import { login, STORAGE_STATE } from './auth.fixture'

/**
 * Extended test fixtures for true E2E testing.
 * No API mocking - tests run against the real backend.
 */
export const test = base.extend<{
  // An authenticated page ready for testing
  authedPage: Page
}>({
  // Provide an authenticated page with real backend authentication
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

    // Real login against the backend
    await login(page)

    await use(page)
    await context.close()
  },
})

// Re-export expect for convenience
export { expect }

// Re-export test data (kept for assertions, not for mocking)
export * from './test-data'

// Re-export auth helpers (real authentication only)
export { login, signup, logout, isAuthenticated, setAuthToken, E2E_TEST_USER } from './auth.fixture'
