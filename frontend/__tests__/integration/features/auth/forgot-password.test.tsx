// Integration tests for ForgotClient component
import {
  renderWithProviders,
  screen,
} from '../../setup/test-utils'
import ForgotClient from '@/app/auth/forgot/ForgotClient'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    forward: jest.fn(),
    prefetch: jest.fn(),
    refresh: jest.fn(),
  }),
  usePathname: () => '/auth/forgot',
  useSearchParams: () => new URLSearchParams(),
}))

describe('ForgotClient Integration', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders the forgot password form', () => {
    renderWithProviders(<ForgotClient />)

    // Check heading - matches actual component text
    expect(screen.getByRole('heading', { name: /forgot your password/i })).toBeInTheDocument()
  })

  it('shows helpful description text', () => {
    renderWithProviders(<ForgotClient />)

    // Description text from the actual component
    expect(screen.getByText(/enter your email address and we'll send you a link to reset your password/i)).toBeInTheDocument()
  })

  it('renders email input with correct attributes', () => {
    renderWithProviders(<ForgotClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    expect(emailInput).toBeInTheDocument()
    expect(emailInput).toHaveAttribute('type', 'email')
    expect(emailInput).toHaveAttribute('name', 'email')
    expect(emailInput).toBeRequired()
  })

  it('renders submit button', () => {
    renderWithProviders(<ForgotClient />)

    // Button text from actual component: 'Send reset link'
    expect(screen.getByRole('button', { name: /send reset link/i })).toBeInTheDocument()
  })

  it('has sign in link with correct href', () => {
    renderWithProviders(<ForgotClient />)

    const signinLink = screen.getByRole('link', { name: /sign in/i })
    expect(signinLink).toHaveAttribute('href', '/auth/signin')
  })

  it('has logo link to home page', () => {
    renderWithProviders(<ForgotClient />)

    // The Logo component is wrapped in a Link to "/"
    const homeLinks = screen.getAllByRole('link')
    const homeLink = homeLinks.find(link => link.getAttribute('href') === '/')
    expect(homeLink).toBeInTheDocument()
  })
})
