import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EmptyState from './EmptyState'

describe('EmptyState', () => {
  const defaultProps = {
    title: 'No items found',
    description: 'Get started by creating your first item.',
    buttonText: 'Create Item',
    buttonOnClick: jest.fn(),
  }

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('should render title', () => {
    render(<EmptyState {...defaultProps} />)

    expect(screen.getByText('No items found')).toBeInTheDocument()
  })

  it('should render description', () => {
    render(<EmptyState {...defaultProps} />)

    expect(screen.getByText('Get started by creating your first item.')).toBeInTheDocument()
  })

  it('should render button with correct text', () => {
    render(<EmptyState {...defaultProps} />)

    expect(screen.getByRole('button', { name: /Create Item/i })).toBeInTheDocument()
  })

  it('should call buttonOnClick when button is clicked', async () => {
    const user = userEvent.setup()
    const onClickMock = jest.fn()

    render(<EmptyState {...defaultProps} buttonOnClick={onClickMock} />)

    await user.click(screen.getByRole('button', { name: /Create Item/i }))

    expect(onClickMock).toHaveBeenCalledTimes(1)
  })

  it('should render the folder icon', () => {
    render(<EmptyState {...defaultProps} />)

    // The FolderPlus icon should be rendered
    const icon = document.querySelector('svg')
    expect(icon).toBeInTheDocument()
  })

  it('should render with different props', () => {
    render(
      <EmptyState
        title="No boards yet"
        description="Create a new soundboard to get started."
        buttonText="New Board"
        buttonOnClick={jest.fn()}
      />
    )

    expect(screen.getByText('No boards yet')).toBeInTheDocument()
    expect(screen.getByText('Create a new soundboard to get started.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /New Board/i })).toBeInTheDocument()
  })
})
