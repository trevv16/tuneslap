// Integration tests for AccountClient component
// Tests the account settings page with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  mockUser,
} from '../../setup/test-utils'
import AccountClient from '@/app/(authenticated)/account/AccountClient'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/account',
  useSearchParams: () => new URLSearchParams(),
}))

// Mock next-themes
jest.mock('next-themes', () => ({
  useTheme: () => ({
    theme: 'light',
    setTheme: jest.fn(),
  }),
}))

describe('AccountClient Integration', () => {
  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
  })

  it('renders all account sections', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    // All sections should be visible
    expect(screen.getByText('Personal Information')).toBeInTheDocument()
    expect(screen.getByText('Theme Settings')).toBeInTheDocument()
    expect(screen.getByText('Change Password')).toBeInTheDocument()
    expect(screen.getByText('Delete Account')).toBeInTheDocument()
  })

  it('renders user profile section with user data', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    // User profile form should have user data
    await waitFor(() => {
      const nameInput = screen.getByLabelText(/name/i)
      expect(nameInput).toHaveValue(mockUser.name)
    })

    // Email should be displayed (disabled)
    const emailInput = screen.getByLabelText(/email address/i)
    expect(emailInput).toHaveValue(mockUser.email)
    expect(emailInput).toBeDisabled()
  })

  it('shows user initials in avatar', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    // Avatar should show initials (T from Test, U from User)
    expect(screen.getByText('TU')).toBeInTheDocument()
  })

  it('has change avatar button', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByRole('button', { name: /change avatar/i })).toBeInTheDocument()
  })

  it('renders theme toggle section', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByText('Theme Settings')).toBeInTheDocument()
    expect(screen.getByText(/choose your preferred theme/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /light mode/i })).toBeInTheDocument()
  })

  it('shows current theme status', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByText(/currently using light theme/i)).toBeInTheDocument()
  })

  it('renders change password section', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByText('Change Password')).toBeInTheDocument()
    expect(screen.getByText(/update your password/i)).toBeInTheDocument()
  })

  it('renders delete account section with warning', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByText('Delete Account')).toBeInTheDocument()
    // Should show some warning text about deletion
    expect(screen.getByText(/delete/i)).toBeInTheDocument()
  })

  it('has save button in profile section', async () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
  })

  it('validates profile name is required', async () => {
    const user = userEvent.setup()
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    const nameInput = screen.getByLabelText(/name/i)
    const saveButton = screen.getByRole('button', { name: /save/i })

    // Clear the name and submit
    await user.clear(nameInput)
    await user.click(saveButton)

    await waitFor(() => {
      expect(screen.getByText(/name is required|name must be at least/i)).toBeInTheDocument()
    })
  })

  it('shows loading state when updating profile', async () => {
    const user = userEvent.setup()
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    const nameInput = screen.getByLabelText(/name/i)
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.type(nameInput, ' Updated')
    await user.click(saveButton)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /saving/i })).toBeInTheDocument()
    })
  })

  it('shows success toast after updating profile', async () => {
    const user = userEvent.setup()
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    const nameInput = screen.getByLabelText(/name/i)
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.type(nameInput, ' Updated')
    await user.click(saveButton)

    await waitFor(() => {
      expect(screen.getByText(/profile updated successfully/i)).toBeInTheDocument()
    })
  })

  it('has accessibility heading for screen readers', () => {
    renderWithProviders(<AccountClient />, {
      authState: { isAuthenticated: true, user: mockUser },
    })

    expect(screen.getByRole('heading', { name: /account settings/i })).toBeInTheDocument()
  })
})
