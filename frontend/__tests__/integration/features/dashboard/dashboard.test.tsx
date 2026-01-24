// Integration tests for DashboardClient component
import {
  renderWithProviders,
  screen,
  waitFor,
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
    forward: jest.fn(),
    prefetch: jest.fn(),
    refresh: jest.fn(),
  }),
  usePathname: () => '/dashboard',
  useSearchParams: () => new URLSearchParams(),
}))

describe('DashboardClient Integration', () => {
  const mockOnOpenCreateModal = jest.fn()

  beforeEach(() => {
    localStorage.setItem('tuneslap_api_token', 'mock-token')
    mockOnOpenCreateModal.mockClear()
  })

  it('renders boards list when boards exist', async () => {
    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByText(mockBoards[0].name)).toBeInTheDocument()
    }, { timeout: 3000 })
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
  })

  it('empty state has new board button', async () => {
    server.use(errorHandlers.emptyBoards)

    renderWithProviders(
      <DashboardClient onOpenCreateModal={mockOnOpenCreateModal} />,
      { authState: { isAuthenticated: true } }
    )

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /new board/i })).toBeInTheDocument()
    })
  })
})
