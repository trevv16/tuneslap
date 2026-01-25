import { test as base, expect, type Page } from '@playwright/test'
import { testUser, mockSigninResponse, mockUserResponse } from './test-data'

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
 * Log in a user via the UI.
 */
export async function login(
  page: Page,
  email: string = testUser.email,
  password: string = testUser.password
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
 * Log in a user with mocked API response.
 * Use this for faster tests that don't need real authentication.
 */
export async function loginWithMock(page: Page): Promise<void> {
  // Mock the signin API
  await page.route('**/api/v1/auth/signin', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockSigninResponse),
    })
  })

  // Mock the /me endpoint for user data
  await page.route('**/api/v1/users/me', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockUserResponse),
    })
  })

  await page.goto('/auth/signin')

  // Fill in credentials
  await page.getByLabel('Email address').fill(testUser.email)
  await page.getByLabel('Password').fill(testUser.password)

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

export { expect }
