// Integration tests for BoardDetailClient component
// Tests the board detail page with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  mockBoard,
  mockKeys,
} from '../../setup/test-utils'
import BoardDetailClient from '@/app/(authenticated)/boards/[boardId]/BoardDetailClient'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/boards/board-1',
  useSearchParams: () => new URLSearchParams(),
  useParams: () => ({ boardId: 'board-1' }),
}))

// Mock next/link
jest.mock('next/link', () => {
  const Link = ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  )
  Link.displayName = 'MockLink'
  return Link
})

describe('BoardDetailClient Integration', () => {
  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
  })

  it('renders board name in header', async () => {
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })
  })

  it('shows edit and add key buttons', async () => {
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByRole('link', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add key/i })).toBeInTheDocument()
  })

  it('edit link points to edit page', async () => {
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    const editLink = screen.getByRole('link', { name: /edit/i })
    expect(editLink).toHaveAttribute('href', '/boards/board-1/edit')
  })

  it('displays keys grid when keys exist', async () => {
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockKeys[0].name)).toBeInTheDocument()
    })

    // All keys should be visible
    mockKeys.forEach(key => {
      expect(screen.getByText(key.name)).toBeInTheDocument()
    })
  })

  it('displays hotkeys on key cards', async () => {
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockKeys[0].name)).toBeInTheDocument()
    })

    // Hotkeys should be displayed
    expect(screen.getByText('A')).toBeInTheDocument()
    expect(screen.getByText('B')).toBeInTheDocument()
    expect(screen.getByText('C')).toBeInTheDocument()
  })

  it('opens add key sheet when add key button is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    const addKeyButton = screen.getByRole('button', { name: /add key/i })
    await user.click(addKeyButton)

    // Sheet should open with title
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument()
      expect(screen.getByText(/create a new key/i)).toBeInTheDocument()
    })
  })

  it('shows key form fields in add key sheet', async () => {
    const user = userEvent.setup()
    renderWithProviders(<BoardDetailClient />, {
      authState: { isAuthenticated: true },
    })

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    const addKeyButton = screen.getByRole('button', { name: /add key/i })
    await user.click(addKeyButton)

    await waitFor(() => {
      expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/hotkey/i)).toBeInTheDocument()
      expect(screen.getByText(/audio/i)).toBeInTheDocument()
    })
  })
})
