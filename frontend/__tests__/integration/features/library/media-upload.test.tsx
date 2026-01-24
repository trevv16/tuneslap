// Integration tests for CreateMediaForm component
// Tests the media upload form with real hooks and MSW network mocking
import {
  renderWithProviders,
  screen,
  waitFor,
  userEvent,
} from '../../setup/test-utils'
import CreateMediaForm from '@/app/(authenticated)/library/components/CreateMediaForm'

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

// Mock fetch for file upload to GCS
const mockFetch = jest.fn()
global.fetch = mockFetch

describe('CreateMediaForm Integration', () => {
  const mockSetOpen = jest.fn()

  beforeEach(() => {
    localStorage.setItem('auth_token', 'mock-token')
    mockSetOpen.mockClear()
    mockFetch.mockClear()
    // Default successful fetch response for upload
    mockFetch.mockResolvedValue({
      ok: true,
      status: 200,
    })
  })

  it('renders the upload form with all elements', () => {
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByText(/file/i)).toBeInTheDocument()
    expect(screen.getByText(/drag and drop/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/file name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /upload media/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('shows validation error when no file is selected', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const fileNameInput = screen.getByLabelText(/file name/i)
    const submitButton = screen.getByRole('button', { name: /upload media/i })

    await user.type(fileNameInput, 'my-file')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/please select a file/i)).toBeInTheDocument()
    })
  })

  it('shows file preview when file is selected', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const file = new File(['audio content'], 'test.mp3', { type: 'audio/mpeg' })
    
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    await user.upload(input, file)

    await waitFor(() => {
      expect(screen.getByText('test.mp3')).toBeInTheDocument()
    })
  })

  it('auto-populates file name from selected file', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const file = new File(['audio content'], 'my-audio-file.mp3', { type: 'audio/mpeg' })
    
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    await user.upload(input, file)

    await waitFor(() => {
      const fileNameInput = screen.getByLabelText(/file name/i)
      expect(fileNameInput).toHaveValue('my-audio-file')
    })
  })

  it('shows file extension preview', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const file = new File(['audio content'], 'test.mp3', { type: 'audio/mpeg' })
    
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    await user.upload(input, file)

    await waitFor(() => {
      expect(screen.getByText('.mp3')).toBeInTheDocument()
    })
  })

  it('allows removing selected file', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const file = new File(['audio content'], 'test.mp3', { type: 'audio/mpeg' })
    
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    await user.upload(input, file)

    await waitFor(() => {
      expect(screen.getByText('test.mp3')).toBeInTheDocument()
    })

    const removeButton = screen.getByRole('button', { name: /remove file/i })
    await user.click(removeButton)

    await waitFor(() => {
      expect(screen.queryByText('test.mp3')).not.toBeInTheDocument()
      expect(screen.getByText(/drag and drop/i)).toBeInTheDocument()
    })
  })

  it('closes form when cancel button is clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const cancelButton = screen.getByRole('button', { name: /cancel/i })
    await user.click(cancelButton)

    expect(mockSetOpen).toHaveBeenCalledWith(false)
  })

  it('accepts image files', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    const file = new File(['image content'], 'test.png', { type: 'image/png' })
    
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    await user.upload(input, file)

    await waitFor(() => {
      expect(screen.getByText('test.png')).toBeInTheDocument()
      expect(screen.getByText('.png')).toBeInTheDocument()
    })
  })

  it('shows file size in preview', async () => {
    const user = userEvent.setup()
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    // Create a file with some content
    const content = 'a'.repeat(1024 * 1024) // 1MB
    const file = new File([content], 'test.mp3', { type: 'audio/mpeg' })
    
    const input = document.querySelector('input[type="file"]') as HTMLInputElement
    await user.upload(input, file)

    await waitFor(() => {
      // Should show size in MB format
      expect(screen.getByText(/1\.00 MB/i)).toBeInTheDocument()
    })
  })

  it('shows demo banner with limitations', () => {
    renderWithProviders(
      <CreateMediaForm setOpen={mockSetOpen} />,
      { authState: { isAuthenticated: true } }
    )

    expect(screen.getByText(/demo mode/i)).toBeInTheDocument()
    expect(screen.getByText(/max file size/i)).toBeInTheDocument()
  })
})
