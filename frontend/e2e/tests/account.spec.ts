import { test, expect } from '@playwright/test'
import { AccountPage } from '../pages/account.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockUserResponse, testUser } from '../fixtures/test-data'

test.describe('Account Settings', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should display user profile section', async ({ page }) => {
    await accountPage.goto()

    await accountPage.expectPageLoaded()
    await accountPage.expectProfileSectionVisible()
  })

  test('should display user name and email', async ({ page }) => {
    await accountPage.goto()

    await accountPage.expectUserInfo(testUser.name, testUser.email)
  })

  test('should display theme section', async ({ page }) => {
    await accountPage.goto()

    await accountPage.expectThemeSectionVisible()
  })

  test('should display theme toggle', async ({ page }) => {
    await accountPage.goto()

    await accountPage.expectThemeToggleVisible()
  })

  test('should display change password section', async ({ page }) => {
    await accountPage.goto()

    await accountPage.expectChangePasswordSectionVisible()
  })

  test('should display delete account section with warning', async ({ page }) => {
    await accountPage.goto()

    await accountPage.expectDeleteAccountSectionVisible()
  })
})

test.describe('Theme Toggle', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should switch to dark theme', async ({ page }) => {
    await accountPage.goto()

    await accountPage.selectDarkTheme()

    await accountPage.expectDarkTheme()
  })

  test('should switch to light theme', async ({ page }) => {
    // First set to dark theme
    await page.emulateMedia({ colorScheme: 'dark' })

    await accountPage.goto()

    await accountPage.selectLightTheme()

    await accountPage.expectLightTheme()
  })

  test('should persist theme preference', async ({ page, context }) => {
    await accountPage.goto()

    await accountPage.selectDarkTheme()

    // Reload the page
    await page.reload()

    // Theme should still be dark
    await accountPage.expectDarkTheme()
  })
})

test.describe('Change Password', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should have change password form fields', async ({ page }) => {
    await accountPage.goto()

    await expect(accountPage.currentPasswordInput).toBeVisible()
    await expect(accountPage.newPasswordInput).toBeVisible()
    await expect(accountPage.confirmPasswordInput).toBeVisible()
    await expect(accountPage.changePasswordButton).toBeVisible()
  })

  test('should submit change password form', async ({ page }) => {
    // Mock the change password API
    await page.route('**/api/v1/users/me/password', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, message: 'Password changed' }),
      })
    })

    await accountPage.goto()

    await accountPage.fillChangePasswordForm(
      'currentPassword123',
      'newPassword456',
      'newPassword456'
    )
    await accountPage.submitChangePassword()

    // Should show success message or the form should reset
    await expect(page.getByText(/success|changed/i)).toBeVisible()
  })
})

test.describe('Delete Account', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should have delete account button', async ({ page }) => {
    await accountPage.goto()

    await expect(accountPage.deleteAccountButton).toBeVisible()
  })

  test('should show confirmation dialog on delete click', async ({ page }) => {
    await accountPage.goto()

    await accountPage.initiateDeleteAccount()

    // Should show confirmation dialog or input
    await expect(
      accountPage.deleteConfirmInput.or(page.getByText(/confirm|are you sure/i))
    ).toBeVisible()
  })
})

test.describe('Sign Out', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should have sign out button', async ({ page }) => {
    await accountPage.goto()

    await expect(accountPage.signOutButton).toBeVisible()
  })

  test('should sign out and redirect to sign in', async ({ page }) => {
    await accountPage.goto()

    await accountPage.signOut()

    await accountPage.expectRedirectToSignIn()
  })

  test('should clear token on sign out', async ({ page }) => {
    await accountPage.goto()

    await accountPage.signOut()

    // Token should be removed from localStorage
    const token = await page.evaluate(() => localStorage.getItem('token'))
    expect(token).toBeNull()
  })
})
