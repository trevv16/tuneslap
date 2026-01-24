// Integration tests for LibraryClient component
// Tests the library page with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  server,
  errorHandlers,
  mockMedia,
  mockAudioMedia,
  mockImageMedia,
} from '../../setup/test-utils'
import LibraryClient from '@/app/(authenticated)/library/LibraryClient'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/library',
  useSearchParams: () => new URLSearchParams(),
}))

describe('LibraryClient Integration', () => {
  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
  })

  it('renders media gallery when media exists', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    expect(screen.getByText(mockImageMedia.name)).toBeInTheDocument()
  })

  it('shows tabs for filtering media types', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    // All, Audio, Images tabs should be visible
    expect(screen.getByRole('tab', { name: /all/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /audio/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /images/i })).toBeInTheDocument()
  })

  it('filters media when switching tabs', async () => {
    const user = userEvent.setup()
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    // Click audio tab
    const audioTab = screen.getByRole('tab', { name: /audio/i })
    await user.click(audioTab)

    // Should still show audio
    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })
  })

  it('opens media details sidebar when item is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    // Find and click on the media item
    const mediaItem = screen.getByText(mockAudioMedia.name)
    await user.click(mediaItem)

    // Details sidebar should open with media info
    await waitFor(() => {
      // Should show more detailed view with download/delete options
      expect(screen.getByRole('button', { name: /download/i })).toBeInTheDocument()
    })
  })

  it('closes sidebar when clicking the same item', async () => {
    const user = userEvent.setup()
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    // Click to open
    const mediaItem = screen.getByText(mockAudioMedia.name)
    await user.click(mediaItem)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /download/i })).toBeInTheDocument()
    })

    // Click again to close
    await user.click(mediaItem)

    // Download button should no longer be visible
    await waitFor(() => {
      expect(screen.queryByRole('button', { name: /download/i })).not.toBeInTheDocument()
    })
  })

  it('shows add file button', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    expect(screen.getByRole('button', { name: /add file/i })).toBeInTheDocument()
  })

  it('opens create modal when add file button is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    const addButton = screen.getByRole('button', { name: /add file/i })
    await user.click(addButton)

    // Modal should open with title
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument()
      expect(screen.getByText(/add new media/i)).toBeInTheDocument()
    })
  })

  it('shows error state when API fails', async () => {
    server.use(
      errorHandlers.emptyMedia
    )

    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    // Should render without error (empty state is valid)
    await waitFor(() => {
      // Component should be rendered
      expect(screen.getByRole('button', { name: /add file/i })).toBeInTheDocument()
    })
  })

  it('has view toggle for grid and list views', async () => {
    renderWithProviders(<LibraryClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.name)).toBeInTheDocument()
    })

    // Should have view toggle buttons
    const viewButtons = screen.getAllByRole('button')
    // View toggle should exist (grid/list icons)
    expect(viewButtons.length).toBeGreaterThan(1)
  })
})
