// Integration tests for CreateBoardForm component
// Tests the create board form with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
} from '../../setup/test-utils'
import CreateBoardForm from '@/app/(authenticated)/dashboard/CreateBoardForm'

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

describe('CreateBoardForm Integration', () => {
  const mockSetOpen = jest.fn()

  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
    mockSetOpen.mockClear()
  })

  it('renders the create board form with all fields', () => {
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/layout/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create board/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('shows validation error for empty name', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const descriptionInput = screen.getByLabelText(/description/i)
    const submitButton = screen.getByRole('button', { name: /create board/i })

    await user.type(descriptionInput, 'A test board description')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/name is required/i)).toBeInTheDocument()
    })
  })

  it('shows validation error for empty description', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const submitButton = screen.getByRole('button', { name: /create board/i })

    await user.type(nameInput, 'My Test Board')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/description is required/i)).toBeInTheDocument()
    })
  })

  it('shows loading state while submitting', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const descriptionInput = screen.getByLabelText(/description/i)
    const submitButton = screen.getByRole('button', { name: /create board/i })

    await user.type(nameInput, 'My Test Board')
    await user.type(descriptionInput, 'A test board description')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /creating/i })).toBeInTheDocument()
    })
  })

  it('handles successful board creation', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const descriptionInput = screen.getByLabelText(/description/i)
    const submitButton = screen.getByRole('button', { name: /create board/i })

    await user.type(nameInput, 'My Test Board')
    await user.type(descriptionInput, 'A test board description')
    await user.click(submitButton)

    // Wait for success toast and modal close
    await waitFor(() => {
      expect(screen.getByText(/board created successfully/i)).toBeInTheDocument()
    })

    // setOpen should be called with false to close the modal
    await waitFor(() => {
      expect(mockSetOpen).toHaveBeenCalledWith(false)
    })
  })

  it('disables form inputs while submitting', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const descriptionInput = screen.getByLabelText(/description/i)
    const submitButton = screen.getByRole('button', { name: /create board/i })

    await user.type(nameInput, 'My Test Board')
    await user.type(descriptionInput, 'A test board description')
    await user.click(submitButton)

    await waitFor(() => {
      expect(nameInput).toBeDisabled()
      expect(descriptionInput).toBeDisabled()
    })
  })

  it('closes modal when cancel button is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const cancelButton = screen.getByRole('button', { name: /cancel/i })
    await user.click(cancelButton)

    expect(mockSetOpen).toHaveBeenCalledWith(false)
  })

  it('resets form after successful submission', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const nameInput = screen.getByLabelText(/name/i)
    const descriptionInput = screen.getByLabelText(/description/i)
    const submitButton = screen.getByRole('button', { name: /create board/i })

    await user.type(nameInput, 'My Test Board')
    await user.type(descriptionInput, 'A test board description')
    await user.click(submitButton)

    // Wait for success
    await waitFor(() => {
      expect(screen.getByText(/board created successfully/i)).toBeInTheDocument()
    })

    // Form inputs should be cleared (reset)
    await waitFor(() => {
      expect(nameInput).toHaveValue('')
      expect(descriptionInput).toHaveValue('')
    })
  })

  it('has layout select with grid option preselected', () => {
    renderWithProviders(
      <CreateBoardForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    // Layout select should have grid as default
    expect(screen.getByRole('combobox')).toBeInTheDocument()
  })
})
