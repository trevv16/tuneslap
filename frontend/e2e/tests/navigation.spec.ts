import { test, expect } from '@playwright/test'
import { login } from '../fixtures/auth.fixture'

test.describe('Authenticated Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Real authentication against the backend
    await login(page)
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
  // These tests verify the backend correctly redirects unauthenticated users
  test('should redirect to sign in when accessing dashboard without auth', async ({ page }) => {
    // Don't login - just try to access protected route directly
    await page.goto('/dashboard')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing library without auth', async ({ page }) => {
    // Don't login - just try to access protected route directly
    await page.goto('/library')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing account without auth', async ({ page }) => {
    // Don't login - just try to access protected route directly
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
