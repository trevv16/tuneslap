import { test, expect } from '@playwright/test'
import { AuthPage } from '../pages/auth.page'
import { E2E_TEST_USER } from '../fixtures/auth.fixture'

test.describe('Sign In', () => {
  let authPage: AuthPage

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page)
  })

  test('should display sign in form correctly', async ({ page }) => {
    await authPage.gotoSignIn()
    await authPage.expectSignInFormVisible()

    // Check for navigation links
    await expect(authPage.forgotPasswordLink).toBeVisible()
    await expect(authPage.signUpLink).toBeVisible()
  })

  test('should show validation error for empty password', async ({ page }) => {
    await authPage.gotoSignIn()

    await authPage.fillSignInForm('test@example.com', '')
    await authPage.submitSignIn()

    await authPage.expectFieldError('required')
  })

  test('should navigate to forgot password page', async ({ page }) => {
    await authPage.gotoSignIn()

    await authPage.forgotPasswordLink.click()

    await expect(page).toHaveURL(/\/auth\/forgot/)
  })

  test('should navigate to sign up page', async ({ page }) => {
    await authPage.gotoSignIn()

    await authPage.signUpLink.click()

    await expect(page).toHaveURL(/\/auth\/signup/)
  })

  test('should sign in successfully with valid credentials', async ({ page }) => {
    // Use E2E test user credentials for real authentication
    await authPage.gotoSignIn()
    await authPage.signIn(E2E_TEST_USER.email, E2E_TEST_USER.password)

    await authPage.expectRedirectToDashboard()
  })

  test('should show error for invalid credentials', async ({ page }) => {
    // Test with wrong credentials against real backend
    await authPage.gotoSignIn()
    await authPage.signIn('wrong@example.com', 'wrongpassword')

    // Should stay on sign in page or show error toast
    await expect(page).toHaveURL(/\/auth\/signin/)
  })
})

test.describe('Sign Up', () => {
  let authPage: AuthPage

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page)
  })

  test('should display sign up form correctly', async ({ page }) => {
    await authPage.gotoSignUp()
    await authPage.expectSignUpFormVisible()

    // Check for sign in link
    await expect(authPage.signInLink).toBeVisible()
  })

  test('should show validation error for short name', async ({ page }) => {
    await authPage.gotoSignUp()

    await authPage.fillSignUpForm('Ab', 'test@example.com', 'password123')
    await authPage.submitSignUp()

    await authPage.expectFieldError('at least 3')
  })

  test('should show validation error for short password', async ({ page }) => {
    await authPage.gotoSignUp()

    await authPage.fillSignUpForm('Test User', 'test@example.com', 'short')
    await authPage.submitSignUp()

    await authPage.expectFieldError('at least 8')
  })

  test('should sign up successfully and redirect to sign in', async ({ page }) => {
    // Use a unique email for each test run to avoid conflicts
    const uniqueEmail = `e2e-test-${Date.now()}@tuneslap.test`

    await authPage.gotoSignUp()
    await authPage.signUp('E2E New User', uniqueEmail, 'password123456')

    await authPage.expectRedirectToSignIn()
  })

  test('should navigate to sign in page', async ({ page }) => {
    await authPage.gotoSignUp()

    await authPage.signInLink.click()

    await expect(page).toHaveURL(/\/auth\/signin/)
  })
})

test.describe('Forgot Password', () => {
  let authPage: AuthPage

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page)
  })

  test('should display forgot password form', async ({ page }) => {
    await authPage.gotoForgotPassword()
    await authPage.expectForgotPasswordFormVisible()
  })
})
