import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SoundBoard from './SoundBoard'
import { mockKeys } from '@/__mocks__/data/fixtures'

// Mock the SoundKey component since it has its own tests
jest.mock('./SoundKey', () => {
  return function MockSoundKey({ boardKey }: { boardKey: { id: string; name: string } }) {
    return <div data-testid={`sound-key-${boardKey.id}`}>{boardKey.name}</div>
  }
})

describe('SoundBoard', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('when keys array is empty', () => {
    it('should render empty state', () => {
      render(<SoundBoard keys={[]} />)

      expect(screen.getByText('No keys yet')).toBeInTheDocument()
      expect(screen.getByText(/Add keys to your soundboard/i)).toBeInTheDocument()
    })

    it('should render add key button when onAddKey is provided', () => {
      render(<SoundBoard keys={[]} onAddKey={jest.fn()} />)

      expect(screen.getByRole('button', { name: /Add Your First Key/i })).toBeInTheDocument()
    })

    it('should not render add key button when onAddKey is not provided', () => {
      render(<SoundBoard keys={[]} />)

      expect(screen.queryByRole('button', { name: /Add Your First Key/i })).not.toBeInTheDocument()
    })

    it('should call onAddKey when button is clicked', async () => {
      const user = userEvent.setup()
      const onAddKeyMock = jest.fn()

      render(<SoundBoard keys={[]} onAddKey={onAddKeyMock} />)

      await user.click(screen.getByRole('button', { name: /Add Your First Key/i }))

      expect(onAddKeyMock).toHaveBeenCalledTimes(1)
    })

    it('should render keyboard icon in empty state', () => {
      render(<SoundBoard keys={[]} />)

      // Check that SVG icon is present
      const icon = document.querySelector('svg')
      expect(icon).toBeInTheDocument()
    })
  })

  describe('when keys array has items', () => {
    it('should render keys grid', () => {
      render(<SoundBoard keys={mockKeys} />)

      // Should not show empty state
      expect(screen.queryByText('No keys yet')).not.toBeInTheDocument()

      // Should render SoundKey for each key
      mockKeys.forEach(key => {
        expect(screen.getByTestId(`sound-key-${key.id}`)).toBeInTheDocument()
      })
    })

    it('should render correct number of keys', () => {
      render(<SoundBoard keys={mockKeys} />)

      const keyElements = screen.getAllByTestId(/sound-key-/)
      expect(keyElements).toHaveLength(mockKeys.length)
    })

    it('should render keys within a list', () => {
      render(<SoundBoard keys={mockKeys} />)

      const list = screen.getByRole('list')
      expect(list).toBeInTheDocument()
    })

    it('should render each key as a list item', () => {
      render(<SoundBoard keys={mockKeys} />)

      const listItems = screen.getAllByRole('listitem')
      expect(listItems).toHaveLength(mockKeys.length)
    })
  })

  describe('with single key', () => {
    it('should render single key correctly', () => {
      const singleKey = [mockKeys[0]]

      render(<SoundBoard keys={singleKey} />)

      expect(screen.getByTestId(`sound-key-${mockKeys[0].id}`)).toBeInTheDocument()
      expect(screen.queryByText('No keys yet')).not.toBeInTheDocument()
    })
  })
})
