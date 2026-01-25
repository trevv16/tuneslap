import { test, expect } from '@playwright/test'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockUserResponse, mockBoards } from '../fixtures/test-data'

test.describe('Authenticated Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.boards.list(mockBoards)
    await apiMocks.media.list()

    // Navigate to a page without SSR API calls and set auth token
    await page.goto('/auth/signin')
    await setAuthToken(page)
  })

  test('should navigate from dashboard to library', async ({ page }) => {
    await page.goto('/dashboard')

    await page.getByRole('link', { name: /library/i }).click()

    await expect(page).toHaveURL(/\/library/)
  })

  test('should navigate from library to dashboard', async ({ page }) => {
    await page.goto('/library')

    await page.getByRole('link', { name: /dashboard|boards/i }).click()

    await expect(page).toHaveURL(/\/dashboard/)
  })
})

test.describe('Protected Routes', () => {
  test('should redirect to sign in when accessing dashboard without auth', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/dashboard')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing library without auth', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/library')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing account without auth', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/account')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })
})

test.describe('Public Routes', () => {
  test('should access sign in page without auth', async ({ page }) => {
    await page.goto('/auth/signin')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should access sign up page without auth', async ({ page }) => {
    await page.goto('/auth/signup')

    await expect(page).toHaveURL(/\/auth\/signup/)
  })
})
