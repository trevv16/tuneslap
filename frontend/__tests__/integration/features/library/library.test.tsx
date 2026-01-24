// Integration tests for LibraryClient component
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  mockAudioMedia,
} from '../../setup/test-utils'
import LibraryClient from '@/app/(authenticated)/library/LibraryClient'

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
  usePathname: () => '/library',
  useSearchParams: () => new URLSearchParams(),
}))

describe('LibraryClient Integration', () => {
  beforeEach(() => {
    localStorage.setItem('tuneslap_api_token', 'mock-token')
  })

  it('renders library page with add file button', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add file/i })).toBeInTheDocument()
    })
  })

  it('shows tabs for filtering media types', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByRole('tab', { name: /all/i })).toBeInTheDocument()
      expect(screen.getByRole('tab', { name: /audio/i })).toBeInTheDocument()
      expect(screen.getByRole('tab', { name: /images/i })).toBeInTheDocument()
    })
  })

  it('renders media gallery when media exists', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.fileName)).toBeInTheDocument()
    }, { timeout: 3000 })
  })

  it('opens create modal when add file button is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add file/i })).toBeInTheDocument()
    })

    const addButton = screen.getByRole('button', { name: /add file/i })
    await user.click(addButton)

    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument()
      expect(screen.getByText(/add new media/i)).toBeInTheDocument()
    })
  })
})
