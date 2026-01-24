import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import CreateBoardForm from './CreateBoardForm'
import { createWrapper } from '@/__mocks__/providers/allProviders'

// Mock the hooks and toast
jest.mock('@/hooks/boards', () => ({
  useCreateBoard: jest.fn(),
}))

jest.mock('sonner', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
  },
}))

import { useCreateBoard } from '@/hooks/boards'
import { toast } from 'sonner'

describe('CreateBoardForm', () => {
  const mockSetOpen = jest.fn()
  const mockMutateAsync = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    ;(useCreateBoard as jest.Mock).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    })
  })

  it('should render form fields', () => {
    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    expect(screen.getByLabelText(/Name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/Description/i)).toBeInTheDocument()
    expect(screen.getByText(/Layout/i)).toBeInTheDocument()
  })

  it('should render submit and cancel buttons', () => {
    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    expect(screen.getByRole('button', { name: /Create board/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /Cancel/i })).toBeInTheDocument()
  })

  it('should show validation errors for empty fields', async () => {
    const user = userEvent.setup()

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    await user.click(screen.getByRole('button', { name: /Create board/i }))

    await waitFor(() => {
      expect(screen.getByText(/Name is required/i)).toBeInTheDocument()
      expect(screen.getByText(/Description is required/i)).toBeInTheDocument()
    })
  })

  it('should call mutateAsync with form data on submit', async () => {
    const user = userEvent.setup()
    mockMutateAsync.mockResolvedValue({ success: true })

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    await user.type(screen.getByLabelText(/Name/i), 'My New Board')
    await user.type(screen.getByLabelText(/Description/i), 'A great soundboard')
    await user.click(screen.getByRole('button', { name: /Create board/i }))

    await waitFor(() => {
      expect(mockMutateAsync).toHaveBeenCalledWith({
        name: 'My New Board',
        description: 'A great soundboard',
        imageUrl: '',
        layout: 'grid',
      })
    })
  })

  it('should show success toast on successful creation', async () => {
    const user = userEvent.setup()
    mockMutateAsync.mockResolvedValue({ success: true })

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    await user.type(screen.getByLabelText(/Name/i), 'My New Board')
    await user.type(screen.getByLabelText(/Description/i), 'A great soundboard')
    await user.click(screen.getByRole('button', { name: /Create board/i }))

    await waitFor(() => {
      expect(toast.success).toHaveBeenCalledWith('Board created successfully!')
    })
  })

  it('should close dialog on successful creation', async () => {
    const user = userEvent.setup()
    mockMutateAsync.mockResolvedValue({ success: true })

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    await user.type(screen.getByLabelText(/Name/i), 'My New Board')
    await user.type(screen.getByLabelText(/Description/i), 'A great soundboard')
    await user.click(screen.getByRole('button', { name: /Create board/i }))

    await waitFor(() => {
      expect(mockSetOpen).toHaveBeenCalledWith(false)
    })
  })

  it('should show error toast on failed creation', async () => {
    const user = userEvent.setup()
    mockMutateAsync.mockRejectedValue(new Error('Failed'))

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    await user.type(screen.getByLabelText(/Name/i), 'My New Board')
    await user.type(screen.getByLabelText(/Description/i), 'A great soundboard')
    await user.click(screen.getByRole('button', { name: /Create board/i }))

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Failed to create board. Please try again.')
    })
  })

  it('should close dialog when cancel is clicked', async () => {
    const user = userEvent.setup()

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    await user.click(screen.getByRole('button', { name: /Cancel/i }))

    expect(mockSetOpen).toHaveBeenCalledWith(false)
  })

  it('should disable inputs when pending', () => {
    ;(useCreateBoard as jest.Mock).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: true,
    })

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    expect(screen.getByLabelText(/Name/i)).toBeDisabled()
    expect(screen.getByLabelText(/Description/i)).toBeDisabled()
    expect(screen.getByRole('button', { name: /Creating/i })).toBeDisabled()
    expect(screen.getByRole('button', { name: /Cancel/i })).toBeDisabled()
  })

  it('should show "Creating..." text when pending', () => {
    ;(useCreateBoard as jest.Mock).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: true,
    })

    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    expect(screen.getByRole('button', { name: /Creating/i })).toBeInTheDocument()
  })

  it('should have default layout value of grid', () => {
    render(<CreateBoardForm setOpen={mockSetOpen} />, {
      wrapper: createWrapper(),
    })

    // The select trigger should show "Grid" as the default
    expect(screen.getByRole('combobox')).toHaveTextContent('Grid')
  })
})
