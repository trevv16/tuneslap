import { test, expect } from '@playwright/test'
import { AccountPage } from '../pages/account.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockUserResponse } from '../fixtures/test-data'

test.describe('Account Settings', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should display profile section', async () => {
    await accountPage.goto()
    await accountPage.expectProfileSectionVisible()
  })

  test('should display theme section', async () => {
    await accountPage.goto()
    await accountPage.expectThemeSectionVisible()
  })

  test('should display theme toggle', async () => {
    await accountPage.goto()
    await accountPage.expectThemeToggleVisible()
  })

  test('should display change password section', async () => {
    await accountPage.goto()
    await accountPage.expectChangePasswordSectionVisible()
  })

  test('should display delete account section', async () => {
    await accountPage.goto()
    await accountPage.expectDeleteAccountSectionVisible()
  })
})

test.describe('Theme Toggle', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should toggle theme when clicked', async () => {
    await accountPage.goto()

    // Get initial theme state
    const initialHtml = accountPage.page.locator('html')
    const wasDark = await initialHtml.evaluate((el) => el.classList.contains('dark'))

    // Toggle theme
    await accountPage.toggleTheme()

    // Verify theme changed
    if (wasDark) {
      await accountPage.expectLightTheme()
    } else {
      await accountPage.expectDarkTheme()
    }
  })
})

test.describe('Change Password Form', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should have password form fields', async () => {
    await accountPage.goto()

    await expect(accountPage.currentPasswordInput).toBeVisible()
    await expect(accountPage.newPasswordInput).toBeVisible()
    await expect(accountPage.confirmPasswordInput).toBeVisible()
  })
})

test.describe('Delete Account', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should have delete account button', async () => {
    await accountPage.goto()

    await expect(accountPage.deleteAccountButton).toBeVisible()
  })
})

test.describe('Sign Out', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should have sign out button visible', async ({ page }) => {
    await accountPage.goto()

    // Sign out button is in the navbar, may need to open a menu
    const signOutVisible = await accountPage.signOutButton.isVisible()
    
    if (!signOutVisible) {
      // Try opening user menu first
      const userMenu = page.getByRole('button', { name: /menu|user|account/i })
      if (await userMenu.isVisible()) {
        await userMenu.click()
      }
    }

    await expect(accountPage.signOutButton).toBeVisible()
  })
})
