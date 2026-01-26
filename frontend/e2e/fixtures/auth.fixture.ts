import { test as base, expect, type Page } from '@playwright/test'

// E2E Test User credentials - must match server/config/demo.go
// These are used for real authentication against the backend
export const E2E_TEST_USER = {
  email: 'e2e-test@tuneslap.test',
  password: 'e2e-test-password-123',
  name: 'E2E Test User',
}

// Storage state file path for authenticated sessions
export const STORAGE_STATE = 'e2e/.auth/user.json'

/**
 * Extended test fixture with authentication helpers.
 */
export const test = base.extend<{
  authenticatedPage: Page
}>({
  // Provide a page that's already authenticated
  authenticatedPage: async ({ browser }, use) => {
    const context = await browser.newContext({
      storageState: STORAGE_STATE,
    })
    const page = await context.newPage()
    await use(page)
    await context.close()
  },
})

/**
 * Log in a user via the UI against the real backend.
 * Default credentials are the E2E test user seeded by the backend.
 */
export async function login(
  page: Page,
  email: string = E2E_TEST_USER.email,
  password: string = E2E_TEST_USER.password
): Promise<void> {
  await page.goto('/auth/signin')

  // Fill in credentials
  await page.getByLabel('Email address').fill(email)
  await page.getByLabel('Password').fill(password)

  // Submit form
  await page.getByRole('button', { name: 'Sign in' }).click()

  // Wait for navigation to dashboard
  await expect(page).toHaveURL(/\/dashboard/)
}

/**
 * Sign up a new user via the UI.
 */
export async function signup(
  page: Page,
  name: string,
  email: string,
  password: string
): Promise<void> {
  await page.goto('/auth/signup')

  // Fill in registration form
  await page.getByLabel('Name').fill(name)
  await page.getByLabel('Email address').fill(email)
  await page.getByLabel('Password').fill(password)

  // Submit form
  await page.getByRole('button', { name: 'Sign up' }).click()

  // Wait for redirect to signin page
  await expect(page).toHaveURL(/\/auth\/signin/)
}

/**
 * Log out the current user.
 */
export async function logout(page: Page): Promise<void> {
  // Click the user menu or account dropdown
  // Then click sign out
  // This will depend on the actual UI implementation
  await page.evaluate(() => {
    localStorage.removeItem('tuneslap_api_token')
  })
  await page.goto('/auth/signin')
}

/**
 * Check if user is authenticated by verifying token exists.
 */
export async function isAuthenticated(page: Page): Promise<boolean> {
  const token = await page.evaluate(() => {
    return localStorage.getItem('tuneslap_api_token')
  })
  return token !== null
}

/**
 * Set authentication token directly (for setup).
 */
export async function setAuthToken(
  page: Page,
  token: string = 'mock-jwt-token-for-testing'
): Promise<void> {
  await page.evaluate((t) => {
    localStorage.setItem('tuneslap_api_token', t)
  }, token)
}
