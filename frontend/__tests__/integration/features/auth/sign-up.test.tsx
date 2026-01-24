// Integration tests for SignupClient component
// Tests the sign-up form with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  server,
  errorHandlers,
} from '../../setup/test-utils'
import SignupClient from '@/app/auth/signup/SignupClient'

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
    replace: jest.fn(),
    back: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/auth/signup',
  useSearchParams: () => new URLSearchParams(),
}))

describe('SignupClient Integration', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders the sign up form with all fields', () => {
    renderWithProviders(<SignupClient />)

    expect(screen.getByRole('heading', { name: /create your account/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign up/i })).toBeInTheDocument()
  })

  it('shows validation error for name too short', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignupClient />)

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign up/i })

    await user.type(nameInput, 'AB')
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/name must be at least 3 characters/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for invalid email', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignupClient />)

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign up/i })

    await user.type(nameInput, 'Test User')
    await user.type(emailInput, 'invalid-email')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/invalid email address/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for password too short', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignupClient />)

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign up/i })

    await user.type(nameInput, 'Test User')
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'short')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/password must be at least 8 characters/i)).toBeInTheDocument()
    })
  })

  it('shows loading state while submitting', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignupClient />)

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign up/i })

    await user.type(nameInput, 'Test User')
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /creating account/i })).toBeInTheDocument()
    })
  })

  it('handles successful sign up and shows success toast', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SignupClient />)

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign up/i })

    await user.type(nameInput, 'Test User')
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/account created successfully/i)).toBeInTheDocument()
    })
  })

  it('handles duplicate email error from API', async () => {
    server.use(errorHandlers.signupError)

    const user = userEvent.setup()
    renderWithProviders(<SignupClient />)

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email address/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /sign up/i })

    await user.type(nameInput, 'Test User')
    await user.type(emailInput, 'existing@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/email already exists|failed to create account/i)).toBeInTheDocument()
    })
  })

  it('has link to sign in page', () => {
    renderWithProviders(<SignupClient />)

    const signinLink = screen.getByRole('link', { name: /sign in/i })
    expect(signinLink).toHaveAttribute('href', '/auth/signin')
  })

  it('has link to home page via logo', () => {
    renderWithProviders(<SignupClient />)

    // Logo should link to home
    const homeLinks = screen.getAllByRole('link')
    const homeLink = homeLinks.find(link => link.getAttribute('href') === '/')
    expect(homeLink).toBeInTheDocument()
  })
})
