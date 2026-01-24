// Integration tests for SignUpClient component
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  server,
  errorHandlers,
} from '../../setup/test-utils'
import SignupClient from '@/app/auth/signup/SignUpClient'

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
  usePathname: () => '/auth/signup',
  useSearchParams: () => new URLSearchParams(),
}))

describe('SignUpClient Integration', () => {
  beforeEach(() => {
    localStorage.clear()
    mockPush.mockClear()
  })

  it('renders the sign up form with all fields', () => {
    renderWithProviders(<SignupClient />)

    expect(screen.getByRole('heading', { name: /create your account/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign up/i })).toBeInTheDocument()
  })

  it('has sign in link with correct href', () => {
    renderWithProviders(<SignupClient />)

    const signinLink = screen.getByRole('link', { name: /sign in/i })
    expect(signinLink).toHaveAttribute('href', '/auth/signin')
  })

  it('shows success toast on successful sign up', async () => {
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

  it('redirects to sign in page on successful registration', async () => {
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
      expect(mockPush).toHaveBeenCalledWith('/auth/signin')
    })
  })
})
