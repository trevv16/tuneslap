import { render, screen } from '@testing-library/react'
import SoundKey from './SoundKey'
import { mockKey } from '@/__mocks__/data/fixtures'

// Mock the useAudio hook
jest.mock('@/hooks/useAudio', () => ({
  useAudio: jest.fn(() => ({
    play: jest.fn(),
    stop: jest.fn(),
    bindPressHandlers: jest.fn(),
  })),
}))

import { useAudio } from '@/hooks/useAudio'

describe('SoundKey', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    const mockUseAudio = useAudio as jest.Mock
    mockUseAudio.mockReturnValue({
      play: jest.fn(),
      stop: jest.fn(),
      bindPressHandlers: jest.fn(),
    })
  })

  it('should render key name', () => {
    render(<SoundKey boardKey={mockKey} />)

    expect(screen.getByText(mockKey.name)).toBeInTheDocument()
  })

  it('should render hotkey', () => {
    render(<SoundKey boardKey={mockKey} />)

    expect(screen.getByText(mockKey.hotKey ?? '')).toBeInTheDocument()
  })

  it('should render key image when imageUrl is provided', () => {
    render(<SoundKey boardKey={mockKey} />)

    const image = screen.getByAltText('Key Image')
    expect(image).toBeInTheDocument()
    expect(image).toHaveAttribute('src', mockKey.imageUrl)
  })

  it('should render default image when imageUrl is empty', () => {
    const keyWithNoImage = { ...mockKey, imageUrl: '' }

    render(<SoundKey boardKey={keyWithNoImage} />)

    const image = screen.getByAltText('Key Image')
    expect(image).toHaveAttribute('src', '/defaultKey.png')
  })

  it('should render default image when imageUrl is undefined', () => {
    const keyWithNoImage = { ...mockKey, imageUrl: undefined }

    render(<SoundKey boardKey={keyWithNoImage} />)

    const image = screen.getByAltText('Key Image')
    expect(image).toHaveAttribute('src', '/defaultKey.png')
  })

  it('should initialize useAudio with correct parameters', () => {
    render(<SoundKey boardKey={mockKey} />)

    expect(useAudio).toHaveBeenCalledWith(
      mockKey.audioUrl,
      mockKey.hotKey
    )
  })

  it('should handle empty audioUrl', () => {
    const keyWithNoAudio = { ...mockKey, audioUrl: undefined }

    render(<SoundKey boardKey={keyWithNoAudio} />)

    expect(useAudio).toHaveBeenCalledWith('', keyWithNoAudio.hotKey)
  })

  it('should handle key with different hotkey', () => {
    const keyWithDifferentHotKey = { ...mockKey, hotKey: 'Z' }

    render(<SoundKey boardKey={keyWithDifferentHotKey} />)

    expect(useAudio).toHaveBeenCalledWith(keyWithDifferentHotKey.audioUrl, 'Z')
  })

  it('should render play button with screen reader text', () => {
    render(<SoundKey boardKey={mockKey} />)

    const button = screen.getByRole('button')
    expect(button).toBeInTheDocument()
    expect(screen.getByText(`Play ${mockKey.name}`)).toBeInTheDocument()
  })

  it('should call bindPressHandlers on mount', () => {
    const mockBindPressHandlers = jest.fn()
    const mockUseAudio = useAudio as jest.Mock
    mockUseAudio.mockReturnValue({
      play: jest.fn(),
      stop: jest.fn(),
      bindPressHandlers: mockBindPressHandlers,
    })

    render(<SoundKey boardKey={mockKey} />)

    expect(mockBindPressHandlers).toHaveBeenCalled()
  })
})
