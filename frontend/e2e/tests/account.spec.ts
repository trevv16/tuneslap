import { test, expect } from '@playwright/test'
import { AccountPage } from '../pages/account.page'
import { login } from '../fixtures/auth.fixture'

test.describe('Account Settings', () => {
  let accountPage: AccountPage

  test.beforeEach(async ({ page }) => {
    accountPage = new AccountPage(page)

    // Real authentication against the backend
    await login(page)
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

    // Real authentication against the backend
    await login(page)
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

    // Real authentication against the backend
    await login(page)
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

    // Real authentication against the backend
    await login(page)
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

    // Real authentication against the backend
    await login(page)
  })

  test('should have sign out button visible', async () => {
    await accountPage.goto()

    // Sign out button is in the user dropdown menu - need to open it first
    await accountPage.openUserMenu()
    await expect(accountPage.signOutButton).toBeVisible()
  })
})
