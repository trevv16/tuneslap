import { renderHook, act } from '@testing-library/react'
import { useAudio } from './useAudio'

// Mock sonner toast
jest.mock('sonner', () => ({
  toast: {
    error: jest.fn(),
  },
}))

describe('useAudio', () => {
  let mockAudioContext: {
    createBufferSource: jest.Mock
    decodeAudioData: jest.Mock
    destination: object
    close: jest.Mock
  }
  let mockSourceNode: {
    buffer: AudioBuffer | null
    connect: jest.Mock
    start: jest.Mock
    stop: jest.Mock
    disconnect: jest.Mock
  }
  let mockFetchResponse: {
    arrayBuffer: jest.Mock
  }

  beforeEach(() => {
    jest.clearAllMocks()

    mockSourceNode = {
      buffer: null,
      connect: jest.fn(),
      start: jest.fn(),
      stop: jest.fn(),
      disconnect: jest.fn(),
    }

    mockAudioContext = {
      createBufferSource: jest.fn(() => mockSourceNode),
      decodeAudioData: jest.fn().mockResolvedValue({
        duration: 1.0,
        length: 44100,
      }),
      destination: {},
      close: jest.fn().mockResolvedValue(undefined),
    }

    mockFetchResponse = {
      arrayBuffer: jest.fn().mockResolvedValue(new ArrayBuffer(100)),
    }

    ;(window.AudioContext as unknown as jest.Mock) = jest.fn(() => mockAudioContext)
    ;(global.fetch as jest.Mock).mockResolvedValue(mockFetchResponse)
  })

  it('should load audio on mount', async () => {
    const audioUrl = 'https://example.com/audio.mp3'

    renderHook(() => useAudio(audioUrl))

    // Wait for audio to load
    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    expect(global.fetch).toHaveBeenCalledWith(audioUrl)
    expect(mockAudioContext.decodeAudioData).toHaveBeenCalled()
  })

  it('should not load when URL is empty', async () => {
    renderHook(() => useAudio(''))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    expect(global.fetch).not.toHaveBeenCalled()
  })

  it('should provide play function', async () => {
    const audioUrl = 'https://example.com/audio.mp3'

    const { result } = renderHook(() => useAudio(audioUrl))

    // Wait for audio to load
    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    act(() => {
      result.current.play()
    })

    expect(mockAudioContext.createBufferSource).toHaveBeenCalled()
    expect(mockSourceNode.connect).toHaveBeenCalledWith(mockAudioContext.destination)
    expect(mockSourceNode.start).toHaveBeenCalledWith(0)
  })

  it('should provide stop function', async () => {
    const audioUrl = 'https://example.com/audio.mp3'

    const { result } = renderHook(() => useAudio(audioUrl))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    // Play first to create a source
    act(() => {
      result.current.play()
    })

    // Then stop
    act(() => {
      result.current.stop()
    })

    expect(mockSourceNode.stop).toHaveBeenCalled()
    expect(mockSourceNode.disconnect).toHaveBeenCalled()
  })

  it('should call onError callback when fetch fails', async () => {
    const onError = jest.fn()
    const error = new Error('Network error')
    ;(global.fetch as jest.Mock).mockRejectedValue(error)

    renderHook(() => useAudio('https://example.com/audio.mp3', undefined, onError))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    expect(onError).toHaveBeenCalledWith('Network error')
  })

  it('should handle keyboard events when hotKey is provided', async () => {
    const audioUrl = 'https://example.com/audio.mp3'
    const hotKey = 'a'

    renderHook(() => useAudio(audioUrl, hotKey))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    // Simulate keydown
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'a' }))
    })

    expect(mockSourceNode.start).toHaveBeenCalled()

    // Simulate keyup
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keyup', { key: 'a' }))
    })

    expect(mockSourceNode.stop).toHaveBeenCalled()
  })

  it('should not trigger on wrong key', async () => {
    const audioUrl = 'https://example.com/audio.mp3'
    const hotKey = 'a'

    renderHook(() => useAudio(audioUrl, hotKey))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    // Simulate wrong key
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'b' }))
    })

    expect(mockSourceNode.start).not.toHaveBeenCalled()
  })

  it('should provide bindPressHandlers function', async () => {
    const audioUrl = 'https://example.com/audio.mp3'

    const { result } = renderHook(() => useAudio(audioUrl))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    const mockElement = document.createElement('button')

    act(() => {
      result.current.bindPressHandlers(mockElement)
    })

    // Simulate mousedown
    act(() => {
      mockElement.onmousedown?.(new MouseEvent('mousedown') as unknown as MouseEvent)
    })

    expect(mockSourceNode.start).toHaveBeenCalled()

    // Simulate mouseup
    act(() => {
      mockElement.onmouseup?.(new MouseEvent('mouseup') as unknown as MouseEvent)
    })

    expect(mockSourceNode.stop).toHaveBeenCalled()
  })

  it('should clean up AudioContext on unmount', async () => {
    const audioUrl = 'https://example.com/audio.mp3'

    const { unmount } = renderHook(() => useAudio(audioUrl))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    unmount()

    expect(mockAudioContext.close).toHaveBeenCalled()
  })

  it('should handle case-insensitive hotkey matching', async () => {
    const audioUrl = 'https://example.com/audio.mp3'
    const hotKey = 'A' // Uppercase

    renderHook(() => useAudio(audioUrl, hotKey))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    // Press lowercase 'a'
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'a' }))
    })

    expect(mockSourceNode.start).toHaveBeenCalled()
  })

  it('should not retrigger while key is held down', async () => {
    const audioUrl = 'https://example.com/audio.mp3'
    const hotKey = 'a'

    renderHook(() => useAudio(audioUrl, hotKey))

    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 0))
    })

    // First keydown
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'a' }))
    })

    // Second keydown (held)
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'a' }))
    })

    // Should only have started once
    expect(mockSourceNode.start).toHaveBeenCalledTimes(1)
  })
})
