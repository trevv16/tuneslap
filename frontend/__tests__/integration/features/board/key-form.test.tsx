// Integration tests for KeyForm component
// Tests the key creation/edit form with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
  mockKeys,
  mockAudioMedia,
} from '../../setup/test-utils'
import KeyForm from '@/app/(authenticated)/boards/[boardId]/KeyForm'

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
}))

describe('KeyForm Integration', () => {
  const mockOnSubmit = jest.fn()
  const mockOnCancel = jest.fn()

  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
    mockOnSubmit.mockClear()
    mockOnCancel.mockClear()
  })

  it('renders the key form with all fields', () => {
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/hotkey/i)).toBeInTheDocument()
    expect(screen.getByText(/description/i)).toBeInTheDocument()
    expect(screen.getByText(/audio/i)).toBeInTheDocument()
    expect(screen.getByText(/image/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add key/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('shows validation error for name too short', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const hotkeyInput = screen.getByLabelText(/hotkey/i)
    const submitButton = screen.getByRole('button', { name: /add key/i })

    await user.type(nameInput, 'AB')
    await user.type(hotkeyInput, 'Z')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/name must be at least 3 characters/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for missing hotkey', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const submitButton = screen.getByRole('button', { name: /add key/i })

    await user.type(nameInput, 'My New Key')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/hotkey is required/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for missing audio', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const hotkeyInput = screen.getByLabelText(/hotkey/i)
    const submitButton = screen.getByRole('button', { name: /add key/i })

    await user.type(nameInput, 'My New Key')
    await user.type(hotkeyInput, 'Z')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/audio is required/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for duplicate hotkey', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={mockKeys}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const hotkeyInput = screen.getByLabelText(/hotkey/i)
    const submitButton = screen.getByRole('button', { name: /add key/i })

    await user.type(nameInput, 'My New Key')
    await user.type(hotkeyInput, 'A') // Already used by mockKeys[0]
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/this hotkey is already used/i)).toBeInTheDocument()
    })
  })

  it('displays available audio files for selection', async () => {
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    // Wait for media to load
    await waitFor(() => {
      expect(screen.getByText(mockAudioMedia.fileName)).toBeInTheDocument()
    })
  })

  it('displays available images for selection', async () => {
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    // Wait for media to load
    await waitFor(() => {
      // Image media should be visible (as img alt text)
      const images = screen.getAllByRole('img')
      expect(images.length).toBeGreaterThan(0)
    })
  })

  it('calls onCancel when cancel button is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    const cancelButton = screen.getByRole('button', { name: /cancel/i })
    await user.click(cancelButton)

    expect(mockOnCancel).toHaveBeenCalledTimes(1)
  })

  it('disables form when isSubmitting is true', () => {
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="add"
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
        isSubmitting={true}
      />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByLabelText(/name/i)).toBeDisabled()
    expect(screen.getByLabelText(/hotkey/i)).toBeDisabled()
    expect(screen.getByRole('button', { name: /adding/i })).toBeDisabled()
  })

  it('shows edit mode button text when mode is edit', () => {
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={[]}
        mode="edit"
        initialData={mockKeys[0]}
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByRole('button', { name: /save changes/i })).toBeInTheDocument()
  })

  it('pre-fills form with initial data in edit mode', () => {
    renderWithProviders(
      <KeyForm
        boardId="board-1"
        existingKeys={mockKeys}
        mode="edit"
        initialData={mockKeys[0]}
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByLabelText(/name/i)).toHaveValue(mockKeys[0].name)
    expect(screen.getByLabelText(/hotkey/i)).toHaveValue(mockKeys[0].hotKey)
  })
})
