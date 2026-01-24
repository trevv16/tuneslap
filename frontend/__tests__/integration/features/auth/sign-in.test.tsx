// Integration tests for SignInClient component
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
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
    replace: jest.fn(),
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
    localStorage.clear()
    mockPush.mockClear()
  })

  it('renders the sign in form with all fields', () => {
    renderWithProviders(<SignInClient />)

    expect(screen.getByRole('heading', { name: /sign in to your account/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('has forgot password link with correct href', () => {
    renderWithProviders(<SignInClient />)

    const forgotLink = screen.getByRole('link', { name: /forgot password/i })
    expect(forgotLink).toHaveAttribute('href', '/auth/forgot')
  })

  it('has sign up link with correct href', () => {
    renderWithProviders(<SignInClient />)

    const signupLink = screen.getByRole('link', { name: /start a 14 day free trial/i })
    expect(signupLink).toHaveAttribute('href', '/auth/signup')
  })

  it('shows success toast on successful sign in', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/sign in successful/i)).toBeInTheDocument()
    })
  })

  it('shows error toast on API failure', async () => {
    server.use(errorHandlers.signinError)

    const user = userEvent.setup()
    renderWithProviders(<SignInClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign in/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'wrongpassword')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/sign in failed/i)).toBeInTheDocument()
    })
  })
})
