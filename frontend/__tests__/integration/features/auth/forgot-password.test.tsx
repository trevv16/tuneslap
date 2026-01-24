// Integration tests for ForgotClient component
// Tests the forgot password form UI (note: form submission not yet implemented in component)
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
    prefetch: jest.fn(),
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

    expect(screen.getByRole('heading', { name: /forgot your password/i })).toBeInTheDocument()
    expect(screen.getByText(/enter your email address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /send reset link/i })).toBeInTheDocument()
  })

  it('has link back to sign in page', () => {
    renderWithProviders(<ForgotClient />)

    const signinLink = screen.getByRole('link', { name: /sign in/i })
    expect(signinLink).toHaveAttribute('href', '/auth/signin')
  })

  it('has link to home page via logo', () => {
    renderWithProviders(<ForgotClient />)

    const homeLinks = screen.getAllByRole('link')
    const homeLink = homeLinks.find(link => link.getAttribute('href') === '/')
    expect(homeLink).toBeInTheDocument()
  })

  it('shows helpful description text', () => {
    renderWithProviders(<ForgotClient />)

    expect(screen.getByText(/send you a link to reset your password/i)).toBeInTheDocument()
  })

  it('renders email input with correct attributes', () => {
    renderWithProviders(<ForgotClient />)

    const emailInput = screen.getByLabelText(/email address/i)
    expect(emailInput).toHaveAttribute('type', 'email')
    expect(emailInput).toHaveAttribute('autocomplete', 'email')
    expect(emailInput).toBeRequired()
  })
})
