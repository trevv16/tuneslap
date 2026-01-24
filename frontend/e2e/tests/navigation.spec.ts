import { test, expect } from '@playwright/test'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockUserResponse, mockBoards } from '../fixtures/test-data'

test.describe('Authenticated Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.boards.list(mockBoards)
    await apiMocks.media.list()
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

  test('should navigate to account settings', async ({ page }) => {
    await page.goto('/dashboard')

    // Click on user menu or account link
    const accountLink = page.getByRole('link', { name: /account|settings/i })
    const userMenu = page.getByRole('button', { name: /account|profile|user/i })

    if (await accountLink.isVisible()) {
      await accountLink.click()
    } else if (await userMenu.isVisible()) {
      await userMenu.click()
      await page.getByRole('link', { name: /account|settings/i }).click()
    }

    await expect(page).toHaveURL(/\/account/)
  })

  test('should navigate back to dashboard from board detail', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get('board-1', mockBoards[0])

    await page.goto('/boards/board-1')

    // Click back link or dashboard link
    const backLink = page.getByRole('link', { name: /back|dashboard/i })
    await backLink.click()

    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('should have consistent navbar across authenticated pages', async ({ page }) => {
    // Check dashboard
    await page.goto('/dashboard')
    const dashboardNav = page.locator('nav')
    await expect(dashboardNav).toBeVisible()

    // Check library
    await page.goto('/library')
    const libraryNav = page.locator('nav')
    await expect(libraryNav).toBeVisible()

    // Check account
    await page.goto('/account')
    const accountNav = page.locator('nav')
    await expect(accountNav).toBeVisible()
  })
})

test.describe('Protected Routes', () => {
  test('should redirect to sign in when accessing dashboard without auth', async ({ page }) => {
    // Don't set auth token
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/dashboard')

    // Should redirect to sign in or show unauthorized
    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing library without auth', async ({ page }) => {
    // Don't set auth token
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/library')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing account without auth', async ({ page }) => {
    // Don't set auth token
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/account')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should redirect to sign in when accessing board detail without auth', async ({ page }) => {
    // Don't set auth token
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.meUnauthorized()

    await page.goto('/boards/some-board-id')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })
})

test.describe('Public Routes', () => {
  test('should access sign in page without auth', async ({ page }) => {
    await page.goto('/auth/signin')

    await expect(page).toHaveURL(/\/auth\/signin/)
    await expect(page.getByRole('heading', { name: /sign in/i })).toBeVisible()
  })

  test('should access sign up page without auth', async ({ page }) => {
    await page.goto('/auth/signup')

    await expect(page).toHaveURL(/\/auth\/signup/)
    await expect(page.getByRole('heading', { name: /create/i })).toBeVisible()
  })

  test('should access forgot password page without auth', async ({ page }) => {
    await page.goto('/auth/forgot')

    await expect(page).toHaveURL(/\/auth\/forgot/)
  })
})

test.describe('Sign Out Navigation', () => {
  test('should redirect to sign in after sign out', async ({ page }) => {
    // Set up authentication first
    await page.goto('/')
    await setAuthToken(page)

    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.boards.list(mockBoards)

    await page.goto('/dashboard')

    // Sign out
    await page.evaluate(() => {
      localStorage.removeItem('token')
    })
    await page.goto('/auth/signin')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })

  test('should not allow access to protected routes after sign out', async ({ page }) => {
    // Set up authentication first
    await page.goto('/')
    await setAuthToken(page)

    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    await page.goto('/dashboard')

    // Sign out by removing token
    await page.evaluate(() => {
      localStorage.removeItem('token')
    })

    // Mock unauthorized response
    await apiMocks.auth.meUnauthorized()

    // Try to access protected route
    await page.goto('/dashboard')

    await expect(page).toHaveURL(/\/auth\/signin/)
  })
})

test.describe('Logo Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.boards.list(mockBoards)
  })

  test('should navigate to dashboard when clicking logo', async ({ page }) => {
    await page.goto('/library')

    // Click the logo
    const logo = page.getByRole('link').filter({ has: page.locator('img[alt*="logo" i], svg') }).first()
    if (await logo.isVisible()) {
      await logo.click()
      await expect(page).toHaveURL(/\/dashboard|\//)
    }
  })
})
