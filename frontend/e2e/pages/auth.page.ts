import type { Page, Locator } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Page object for authentication pages (signin, signup, forgot password).
 */
export class AuthPage {
  readonly page: Page

  // Sign In locators
  readonly emailInput: Locator
  readonly passwordInput: Locator
  readonly signInButton: Locator
  readonly forgotPasswordLink: Locator
  readonly signUpLink: Locator

  // Sign Up locators
  readonly nameInput: Locator
  readonly signUpButton: Locator
  readonly signInLink: Locator

  // Forgot Password locators
  readonly forgotEmailInput: Locator
  readonly resetButton: Locator

  // Error messages
  readonly errorMessage: Locator
  readonly fieldErrors: Locator

  constructor(page: Page) {
    this.page = page

    // Sign In
    this.emailInput = page.getByLabel('Email address')
    this.passwordInput = page.getByLabel('Password')
    this.signInButton = page.getByRole('button', { name: 'Sign in' })
    this.forgotPasswordLink = page.getByRole('link', { name: 'Forgot password?' })
    this.signUpLink = page.getByRole('link', { name: 'Start a 14 day free trial' })

    // Sign Up
    this.nameInput = page.getByLabel('Name')
    this.signUpButton = page.getByRole('button', { name: 'Sign up' })
    this.signInLink = page.getByRole('link', { name: 'Sign in' })

    // Forgot Password
    this.forgotEmailInput = page.getByLabel('Email address')
    this.resetButton = page.getByRole('button', { name: /send|reset/i })

    // Errors
    this.errorMessage = page.locator('.text-destructive')
    this.fieldErrors = page.locator('p.text-destructive')
  }

  // Navigation methods
  async gotoSignIn(): Promise<void> {
    await this.page.goto('/auth/signin')
  }

  async gotoSignUp(): Promise<void> {
    await this.page.goto('/auth/signup')
  }

  async gotoForgotPassword(): Promise<void> {
    await this.page.goto('/auth/forgot')
  }

  // Sign In actions
  async fillSignInForm(email: string, password: string): Promise<void> {
    await this.emailInput.fill(email)
    await this.passwordInput.fill(password)
  }

  async submitSignIn(): Promise<void> {
    await this.signInButton.click()
  }

  async signIn(email: string, password: string): Promise<void> {
    await this.fillSignInForm(email, password)
    await this.submitSignIn()
  }

  // Sign Up actions
  async fillSignUpForm(name: string, email: string, password: string): Promise<void> {
    await this.nameInput.fill(name)
    await this.emailInput.fill(email)
    await this.passwordInput.fill(password)
  }

  async submitSignUp(): Promise<void> {
    await this.signUpButton.click()
  }

  async signUp(name: string, email: string, password: string): Promise<void> {
    await this.fillSignUpForm(name, email, password)
    await this.submitSignUp()
  }

  // Forgot Password actions
  async fillForgotPasswordForm(email: string): Promise<void> {
    await this.forgotEmailInput.fill(email)
  }

  async submitForgotPassword(): Promise<void> {
    await this.resetButton.click()
  }

  // Assertions
  async expectSignInFormVisible(): Promise<void> {
    await expect(this.page.getByRole('heading', { name: 'Sign in to your account' })).toBeVisible()
    await expect(this.emailInput).toBeVisible()
    await expect(this.passwordInput).toBeVisible()
    await expect(this.signInButton).toBeVisible()
  }

  async expectSignUpFormVisible(): Promise<void> {
    await expect(this.page.getByRole('heading', { name: 'Create your account' })).toBeVisible()
    await expect(this.nameInput).toBeVisible()
    await expect(this.emailInput).toBeVisible()
    await expect(this.passwordInput).toBeVisible()
    await expect(this.signUpButton).toBeVisible()
  }

  async expectForgotPasswordFormVisible(): Promise<void> {
    await expect(this.forgotEmailInput).toBeVisible()
    await expect(this.resetButton).toBeVisible()
  }

  async expectFieldError(message: string): Promise<void> {
    await expect(this.fieldErrors.filter({ hasText: message })).toBeVisible()
  }

  async expectNoFieldErrors(): Promise<void> {
    await expect(this.fieldErrors).toHaveCount(0)
  }

  async expectRedirectToDashboard(): Promise<void> {
    await expect(this.page).toHaveURL(/\/dashboard/)
  }

  async expectRedirectToSignIn(): Promise<void> {
    await expect(this.page).toHaveURL(/\/auth\/signin/)
  }

  async expectLoadingState(): Promise<void> {
    await expect(this.signInButton).toBeDisabled()
  }

  async expectButtonEnabled(): Promise<void> {
    await expect(this.signInButton).toBeEnabled()
  }
}
