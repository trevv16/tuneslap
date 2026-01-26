import type { Page, Locator } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Page object for the account settings page.
 * Uses data-testid attributes for reliable element selection.
 */
export class AccountPage {
  readonly page: Page

  // Sections
  readonly profileSection: Locator
  readonly themeSection: Locator
  readonly changePasswordSection: Locator
  readonly deleteAccountSection: Locator

  // Profile form
  readonly nameInput: Locator
  readonly emailInput: Locator
  readonly saveProfileButton: Locator

  // Theme
  readonly themeToggle: Locator

  // Change password form
  readonly currentPasswordInput: Locator
  readonly newPasswordInput: Locator
  readonly confirmPasswordInput: Locator

  // Delete account
  readonly deleteAccountButton: Locator

  // Sign out (in Navbar dropdown)
  readonly userMenuTrigger: Locator
  readonly signOutButton: Locator

  constructor(page: Page) {
    this.page = page

    // Sections by data-testid
    this.profileSection = page.getByTestId('profile-section')
    this.themeSection = page.getByTestId('theme-section')
    this.changePasswordSection = page.getByTestId('change-password-section')
    this.deleteAccountSection = page.getByTestId('delete-account-section')

    // Profile form elements
    this.nameInput = page.getByLabel(/^name$/i)
    this.emailInput = page.getByLabel(/email/i)
    this.saveProfileButton = this.profileSection.getByRole('button', { name: /save/i })

    // Theme toggle
    this.themeToggle = page.getByTestId('theme-toggle')

    // Change password form
    this.currentPasswordInput = page.getByLabel(/current password/i)
    this.newPasswordInput = page.getByLabel(/^new password$/i)
    this.confirmPasswordInput = page.getByLabel(/confirm password/i)

    // Delete account button
    this.deleteAccountButton = page.getByRole('button', { name: /yes, delete my account/i })

    // Sign out in navbar (dropdown)
    this.userMenuTrigger = page.getByTestId('user-menu-trigger')
    this.signOutButton = page.getByTestId('sign-out-button')
  }

  // Navigation
  async goto(): Promise<void> {
    await this.page.goto('/account')
    await this.page.waitForLoadState('networkidle')
  }

  // Theme actions
  async toggleTheme(): Promise<void> {
    await this.themeToggle.click()
  }

  // Password actions
  async fillChangePasswordForm(
    currentPassword: string,
    newPassword: string,
    confirmPassword: string
  ): Promise<void> {
    await this.currentPasswordInput.fill(currentPassword)
    await this.newPasswordInput.fill(newPassword)
    await this.confirmPasswordInput.fill(confirmPassword)
  }

  // Sign out (opens user menu first, then clicks sign out)
  async signOut(): Promise<void> {
    await this.userMenuTrigger.click()
    await this.signOutButton.click()
  }

  async openUserMenu(): Promise<void> {
    await this.userMenuTrigger.click()
  }

  // Assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.page).toHaveURL(/\/account/)
    // Check for account-specific elements
    await expect(this.profileSection).toBeVisible()
  }

  async expectProfileSectionVisible(): Promise<void> {
    await expect(this.profileSection).toBeVisible()
  }

  async expectThemeSectionVisible(): Promise<void> {
    await expect(this.themeSection).toBeVisible()
  }

  async expectThemeToggleVisible(): Promise<void> {
    await expect(this.themeToggle).toBeVisible()
  }

  async expectChangePasswordSectionVisible(): Promise<void> {
    await expect(this.changePasswordSection).toBeVisible()
  }

  async expectDeleteAccountSectionVisible(): Promise<void> {
    await expect(this.deleteAccountSection).toBeVisible()
  }

  async expectDarkTheme(): Promise<void> {
    await expect(this.page.locator('html')).toHaveClass(/dark/)
  }

  async expectLightTheme(): Promise<void> {
    await expect(this.page.locator('html')).not.toHaveClass(/dark/)
  }

  async expectRedirectToSignIn(): Promise<void> {
    await expect(this.page).toHaveURL(/\/auth\/signin/)
  }
}
