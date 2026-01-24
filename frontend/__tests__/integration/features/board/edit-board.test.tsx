// Integration tests for EditBoardClient component
// Tests the edit board page with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  mockBoard,
  mockKeys,
} from '../../setup/test-utils'
import EditBoardClient from '@/app/(authenticated)/boards/[boardId]/edit/EditBoardClient'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/boards/board-1/edit',
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

describe('EditBoardClient Integration', () => {
  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
  })

  it('renders board header with board info', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText(mockBoard.description!)).toBeInTheDocument()
  })

  it('shows key and collaborator counts in header', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText(`${mockKeys.length} keys`)).toBeInTheDocument()
    expect(screen.getByText(`${mockBoard.collaborators?.length || 0} collaborators`)).toBeInTheDocument()
  })

  it('shows layout badge in header', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText('GRID')).toBeInTheDocument()
  })

  it('has view board link', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    const viewBoardLink = screen.getByRole('link', { name: /view board/i })
    expect(viewBoardLink).toHaveAttribute('href', `/boards/${mockBoard.id}`)
  })

  it('renders board details section', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText('Board Details')).toBeInTheDocument()
    expect(screen.getByText(/update your board name/i)).toBeInTheDocument()
  })

  it('renders collaborators section', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText('Collaborators')).toBeInTheDocument()
    expect(screen.getByText(/manage who has access/i)).toBeInTheDocument()
  })

  it('renders keys section', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText('Keys')).toBeInTheDocument()
    expect(screen.getByText(/configure the sound keys/i)).toBeInTheDocument()
  })

  it('renders danger zone section', async () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoard.name)).toBeInTheDocument()
    })

    expect(screen.getByText('Danger Zone')).toBeInTheDocument()
    expect(screen.getByText(/irreversible actions/i)).toBeInTheDocument()
  })

  it('shows loading skeletons while data is loading', () => {
    renderWithProviders(
      <EditBoardClient boardId="board-1" />,
      { authState: { isAuthenticated: true } }
    )

    // Initially should show skeleton loaders
    // The skeleton components should be present before data loads
    expect(document.querySelector('.animate-pulse')).toBeInTheDocument()
  })
})
