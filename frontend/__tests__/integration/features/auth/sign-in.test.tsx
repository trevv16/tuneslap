// Integration tests for SignInClient component
// Tests the sign-in form with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  server,
  errorHandlers,
} from '../../setup/test-utils'
import SignInClient from '@/app/auth/signin/SignInClient'

// Mock next/navigation
const mockPush = jest.fn()
const mockReplace = jest.fn()
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
    replace: mockReplace,
    back: jest.fn(),
    forward: jest.fn(),
    prefetch: jest.fn(),
    refresh: jest.fn(),
  }),
  usePathname: () => '/auth/signin',
  useSearchParams: () => new URLSearchParams(),
}))

describe('SignInClient Integration', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear()
  })

  it('renders the sign in form with all fields', () => {
    renderWithProviders(<SignInClient />)

    expect(screen.getByRole('heading', { name: /sign in to your account/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /forgot password/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /start a 14 day free trial/i })).toBeInTheDocument()
  })

  it('shows validation error for invalid email format', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'invalid-email')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/please enter a valid email address/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for empty password', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/password is required/i)).toBeInTheDocument()
    })
  })

  it('shows loading state while submitting', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    // Button should show loading state
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /signing in/i })).toBeInTheDocument()
    })
  })

  it('handles successful sign in', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    // Wait for success toast
    await waitFor(() => {
      expect(screen.getByText(/sign in successful/i)).toBeInTheDocument()
    })
  })

  it('handles API error on sign in', async () => {
    // Override with error handler
    server.use(errorHandlers.signinError)

    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'wrongpassword')
    await user.click(submitButton)

    // Wait for error toast
    await waitFor(() => {
      expect(screen.getByText(/sign in failed/i)).toBeInTheDocument()
    })
  })

  it('disables form inputs while submitting', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    // Inputs should be disabled during submission
    await waitFor(() => {
      expect(emailInput).toBeDisabled()
      expect(passwordInput).toBeDisabled()
    })
  })

  it('has link to forgot password page', () => {
    renderWithProviders(<SignInClient />)

    const forgotLink = screen.getByRole('link', { name: /forgot password/i })
    expect(forgotLink).toHaveAttribute('href', '/auth/forgot')
  })

  it('has link to sign up page', () => {
    renderWithProviders(<SignInClient />)

    const signupLink = screen.getByRole('link', { name: /start a 14 day free trial/i })
    expect(signupLink).toHaveAttribute('href', '/auth/signup')
  })
})
