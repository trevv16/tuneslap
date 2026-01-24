import type { Page, Locator } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Page object for the account settings page.
 */
export class AccountPage {
  readonly page: Page

  // Navigation
  readonly accountMenu: Locator
  readonly settingsLink: Locator

  // Profile section
  readonly profileSection: Locator
  readonly userName: Locator
  readonly userEmail: Locator
  readonly editProfileButton: Locator

  // Theme section
  readonly themeSection: Locator
  readonly themeToggle: Locator
  readonly lightThemeOption: Locator
  readonly darkThemeOption: Locator
  readonly systemThemeOption: Locator

  // Change password section
  readonly changePasswordSection: Locator
  readonly currentPasswordInput: Locator
  readonly newPasswordInput: Locator
  readonly confirmPasswordInput: Locator
  readonly changePasswordButton: Locator

  // Delete account section
  readonly deleteAccountSection: Locator
  readonly deleteAccountButton: Locator
  readonly deleteConfirmInput: Locator
  readonly deleteConfirmButton: Locator

  // Sign out
  readonly signOutButton: Locator

  constructor(page: Page) {
    this.page = page

    // Navigation
    this.accountMenu = page.getByRole('button', { name: /account|profile|user/i })
    this.settingsLink = page.getByRole('link', { name: /settings|account/i })

    // Profile section
    this.profileSection = page.locator('[data-testid="profile-section"]').or(
      page.locator('section').filter({ hasText: /profile/i })
    )
    this.userName = page.locator('[data-testid="user-name"]').or(
      page.getByText(/test user/i)
    )
    this.userEmail = page.locator('[data-testid="user-email"]').or(
      page.getByText(/@example\.com/)
    )
    this.editProfileButton = page.getByRole('button', { name: /edit profile/i })

    // Theme section
    this.themeSection = page.locator('[data-testid="theme-section"]').or(
      page.locator('section').filter({ hasText: /theme|appearance/i })
    )
    this.themeToggle = page.getByRole('button', { name: /theme/i }).or(
      page.locator('[data-testid="theme-toggle"]')
    )
    this.lightThemeOption = page.getByRole('menuitem', { name: /light/i }).or(
      page.getByText(/light/i)
    )
    this.darkThemeOption = page.getByRole('menuitem', { name: /dark/i }).or(
      page.getByText(/dark/i)
    )
    this.systemThemeOption = page.getByRole('menuitem', { name: /system/i }).or(
      page.getByText(/system/i)
    )

    // Change password section
    this.changePasswordSection = page.locator('[data-testid="change-password-section"]').or(
      page.locator('section').filter({ hasText: /change password/i })
    )
    this.currentPasswordInput = page.getByLabel(/current password/i)
    this.newPasswordInput = page.getByLabel(/new password/i)
    this.confirmPasswordInput = page.getByLabel(/confirm password/i)
    this.changePasswordButton = page.getByRole('button', { name: /change password|update password/i })

    // Delete account section
    this.deleteAccountSection = page.locator('[data-testid="delete-account-section"]').or(
      page.locator('section').filter({ hasText: /delete account/i })
    )
    this.deleteAccountButton = page.getByRole('button', { name: /delete account/i })
    this.deleteConfirmInput = page.getByPlaceholder(/delete|confirm/i)
    this.deleteConfirmButton = page.getByRole('button', { name: /permanently delete/i })

    // Sign out
    this.signOutButton = page.getByRole('button', { name: /sign out|logout/i })
  }

  // Navigation
  async goto(): Promise<void> {
    await this.page.goto('/account')
  }

  async openAccountMenu(): Promise<void> {
    await this.accountMenu.click()
  }

  async navigateToSettings(): Promise<void> {
    await this.openAccountMenu()
    await this.settingsLink.click()
  }

  // Theme actions
  async openThemeMenu(): Promise<void> {
    await this.themeToggle.click()
  }

  async selectLightTheme(): Promise<void> {
    await this.openThemeMenu()
    await this.lightThemeOption.click()
  }

  async selectDarkTheme(): Promise<void> {
    await this.openThemeMenu()
    await this.darkThemeOption.click()
  }

  async selectSystemTheme(): Promise<void> {
    await this.openThemeMenu()
    await this.systemThemeOption.click()
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

  async submitChangePassword(): Promise<void> {
    await this.changePasswordButton.click()
  }

  // Delete account actions
  async initiateDeleteAccount(): Promise<void> {
    await this.deleteAccountButton.click()
  }

  async confirmDeleteAccount(confirmText: string): Promise<void> {
    await this.deleteConfirmInput.fill(confirmText)
    await this.deleteConfirmButton.click()
  }

  // Sign out
  async signOut(): Promise<void> {
    await this.signOutButton.click()
  }

  // Assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.page).toHaveURL(/\/account/)
  }

  async expectProfileSectionVisible(): Promise<void> {
    await expect(this.profileSection).toBeVisible()
  }

  async expectUserInfo(name: string, email: string): Promise<void> {
    await expect(this.page.getByText(name)).toBeVisible()
    await expect(this.page.getByText(email)).toBeVisible()
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
