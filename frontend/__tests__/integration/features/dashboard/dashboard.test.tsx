// Integration tests for DashboardClient component
// Tests the dashboard with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  server,
  errorHandlers,
  mockBoards,
} from '../../setup/test-utils'
import DashboardClient from '@/app/(authenticated)/dashboard/DashboardClient'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/dashboard',
  useSearchParams: () => new URLSearchParams(),
}))

// Mock next/link
jest.mock('next/link', () => {
  const Link = ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  )
  Link.displayName = 'MockLink'
  return Link
})

describe('DashboardClient Integration', () => {
  const mockOnOpenCreateModal = jest.fn()

  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
    mockOnOpenCreateModal.mockClear()
  })

  it('renders boards list when boards exist', async () => {
    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    // Wait for boards to load
    await waitFor(() => {
      expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    })

    // Both boards should be visible
    expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    expect(screen.getByText(mockBoards[1].name)).toBeInTheDocument()
  })

  it('shows empty state when no boards exist', async () => {
    server.use(errorHandlers.emptyBoards)

    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(/create your first soundboard/i)).toBeInTheDocument()
    })

    expect(screen.getByText(/get started by creating a new board/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new board/i })).toBeInTheDocument()
  })

  it('calls onOpenCreateModal when empty state button is clicked', async () => {
    server.use(errorHandlers.emptyBoards)
    const user = userEvent.setup()

    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /new board/i })).toBeInTheDocument()
    })

    await user.click(screen.getByRole('button', { name: /new board/i }))

    expect(mockOnOpenCreateModal).toHaveBeenCalledTimes(1)
  })

  it('displays board cards with correct information', async () => {
    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    })

    // Check first board details
    expect(screen.getByText(mockBoards[0].description!)).toBeInTheDocument()
    
    // Check second board (has no description)
    expect(screen.getByText(mockBoards[1].name)).toBeInTheDocument()
    expect(screen.getByText(/no description/i)).toBeInTheDocument()
  })

  it('board cards link to correct board detail page', async () => {
    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    })

    // Find the link to the first board
    const boardLinks = screen.getAllByRole('link')
    const firstBoardLink = boardLinks.find(
      link => link.getAttribute('href') === `/boards/${mockBoards[0].id}`
    )
    
    expect(firstBoardLink).toBeInTheDocument()
  })

  it('shows key and collaborator counts on board cards', async () => {
    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    })

    // First board has 3 keys and 1 collaborator
    expect(screen.getByText('3')).toBeInTheDocument()
    expect(screen.getByText('1')).toBeInTheDocument()
  })

  it('shows layout badge on board cards', async () => {
    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    })

    // First board has 'grid' layout, second has 'list'
    expect(screen.getByText('GRID')).toBeInTheDocument()
    expect(screen.getByText('LIST')).toBeInTheDocument()
  })
})
