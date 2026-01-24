import {
  preloadImage,
  preloadAudio,
  preloadImages,
  preloadAudios,
  extractMediaUrls,
} from './preloadUtils'
import { mockBoard, mockKeys } from '@/__mocks__/data/fixtures'

describe('preloadImage', () => {
  let mockImage: { onload: (() => void) | null; onerror: (() => void) | null; src: string }
  const ImageMock = window.Image as jest.Mock

  beforeEach(() => {
    mockImage = { onload: null, onerror: null, src: '' }
    ImageMock.mockImplementation(() => mockImage)
    jest.useFakeTimers()
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('should resolve on successful image load', async () => {
    const url = 'https://example.com/image.png'
    const promise = preloadImage(url)

    // Simulate image load
    mockImage.onload?.()

    await expect(promise).resolves.toBeUndefined()
    expect(mockImage.src).toBe(url)
  })

  it('should reject on image load error', async () => {
    const url = 'https://example.com/broken.png'
    const promise = preloadImage(url)

    // Simulate image error
    mockImage.onerror?.()

    await expect(promise).rejects.toThrow(`Failed to preload image: ${url}`)
  })

  it('should reject on timeout', async () => {
    const url = 'https://example.com/slow.png'
    const promise = preloadImage(url, { timeout: 1000 })

    // Advance timers past timeout
    jest.advanceTimersByTime(1001)

    await expect(promise).rejects.toThrow(`Image preload timeout: ${url}`)
  })

  it('should resolve immediately for empty URL', async () => {
    const promise = preloadImage('')
    await expect(promise).resolves.toBeUndefined()
  })

  it('should call onSuccess callback on successful load', async () => {
    const onSuccess = jest.fn()
    const promise = preloadImage('https://example.com/image.png', { onSuccess })

    mockImage.onload?.()
    await promise

    expect(onSuccess).toHaveBeenCalled()
  })

  it('should call onError callback on error', async () => {
    const onError = jest.fn()
    const url = 'https://example.com/broken.png'
    const promise = preloadImage(url, { onError })

    mockImage.onerror?.()

    await expect(promise).rejects.toThrow()
    expect(onError).toHaveBeenCalled()
  })
})

describe('preloadAudio', () => {
  let mockAudio: { oncanplaythrough: (() => void) | null; onerror: (() => void) | null; src: string; load: jest.Mock }
  const AudioMock = window.Audio as jest.Mock

  beforeEach(() => {
    mockAudio = { oncanplaythrough: null, onerror: null, src: '', load: jest.fn() }
    AudioMock.mockImplementation(() => mockAudio)
    jest.useFakeTimers()
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('should resolve on successful audio load', async () => {
    const url = 'https://example.com/audio.mp3'
    const promise = preloadAudio(url)

    // Simulate audio load
    mockAudio.oncanplaythrough?.()

    await expect(promise).resolves.toBeUndefined()
    expect(mockAudio.src).toBe(url)
    expect(mockAudio.load).toHaveBeenCalled()
  })

  it('should reject on audio load error', async () => {
    const url = 'https://example.com/broken.mp3'
    const promise = preloadAudio(url)

    // Simulate audio error
    mockAudio.onerror?.()

    await expect(promise).rejects.toThrow(`Failed to preload audio: ${url}`)
  })

  it('should reject on timeout', async () => {
    const url = 'https://example.com/slow.mp3'
    const promise = preloadAudio(url, { timeout: 1000 })

    // Advance timers past timeout
    jest.advanceTimersByTime(1001)

    await expect(promise).rejects.toThrow(`Audio preload timeout: ${url}`)
  })

  it('should resolve immediately for empty URL', async () => {
    const promise = preloadAudio('')
    await expect(promise).resolves.toBeUndefined()
  })
})

describe('preloadImages', () => {
  type MockImageType = { onload: (() => void) | null; onerror: (() => void) | null; src: string }
  let mockImages: MockImageType[]
  const ImageMock = window.Image as jest.Mock

  beforeEach(() => {
    mockImages = []
    ImageMock.mockImplementation(() => {
      const img: MockImageType = { onload: null, onerror: null, src: '' }
      mockImages.push(img)
      return img
    })
    jest.useFakeTimers()
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('should preload multiple images', async () => {
    const urls = ['https://example.com/1.png', 'https://example.com/2.png']
    const promise = preloadImages(urls)

    // Simulate all images loading
    mockImages.forEach(img => img.onload?.())

    await expect(promise).resolves.toBeUndefined()
    expect(mockImages).toHaveLength(2)
  })

  it('should filter out empty URLs', async () => {
    const urls = ['https://example.com/1.png', '', 'https://example.com/2.png']
    const promise = preloadImages(urls)

    mockImages.forEach(img => img.onload?.())

    await expect(promise).resolves.toBeUndefined()
    expect(mockImages).toHaveLength(2)
  })

  it('should resolve even if some images fail', async () => {
    const urls = ['https://example.com/1.png', 'https://example.com/broken.png']
    const consoleSpy = jest.spyOn(console, 'warn').mockImplementation()
    const promise = preloadImages(urls)

    // First image loads, second fails
    mockImages[0]?.onload?.()
    mockImages[1]?.onerror?.()

    await expect(promise).resolves.toBeUndefined()
    consoleSpy.mockRestore()
  })

  it('should resolve immediately for empty array', async () => {
    await expect(preloadImages([])).resolves.toBeUndefined()
  })
})

describe('preloadAudios', () => {
  type MockAudioType = { oncanplaythrough: (() => void) | null; onerror: (() => void) | null; src: string; load: jest.Mock }
  let mockAudios: MockAudioType[]
  const AudioMock = window.Audio as jest.Mock

  beforeEach(() => {
    mockAudios = []
    AudioMock.mockImplementation(() => {
      const audio: MockAudioType = { oncanplaythrough: null, onerror: null, src: '', load: jest.fn() }
      mockAudios.push(audio)
      return audio
    })
    jest.useFakeTimers()
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('should preload multiple audio files', async () => {
    const urls = ['https://example.com/1.mp3', 'https://example.com/2.mp3']
    const promise = preloadAudios(urls)

    mockAudios.forEach(audio => audio.oncanplaythrough?.())

    await expect(promise).resolves.toBeUndefined()
    expect(mockAudios).toHaveLength(2)
  })

  it('should resolve immediately for empty array', async () => {
    await expect(preloadAudios([])).resolves.toBeUndefined()
  })
})

describe('extractMediaUrls', () => {
  it('should extract image and audio URLs from board', () => {
    const result = extractMediaUrls(mockBoard)

    expect(result.images).toContain(mockBoard.imageUrl)
    expect(result.images.length).toBeGreaterThan(0)
    expect(result.audios.length).toBeGreaterThan(0)
  })

  it('should include board image URL', () => {
    const result = extractMediaUrls(mockBoard)

    expect(result.images).toContain(mockBoard.imageUrl)
  })

  it('should extract media from keys', () => {
    const result = extractMediaUrls(mockBoard)

    mockKeys.forEach(key => {
      if (key.imageUrl) {
        expect(result.images).toContain(key.imageUrl)
      }
      if (key.audioUrl) {
        expect(result.audios).toContain(key.audioUrl)
      }
    })
  })

  it('should handle board with no keys', () => {
    const boardWithNoKeys = { ...mockBoard, keys: [] }
    const result = extractMediaUrls(boardWithNoKeys)

    expect(result.images).toHaveLength(1) // Just the board image
    expect(result.audios).toHaveLength(0)
  })

  it('should handle board with no image', () => {
    const boardWithNoImage = { ...mockBoard, imageUrl: undefined }
    const result = extractMediaUrls(boardWithNoImage)

    // Should not include undefined board image - all URLs should be strings
    result.images.forEach(url => {
      expect(typeof url).toBe('string')
      expect(url.length).toBeGreaterThan(0)
    })
  })

  it('should filter out empty string URLs', () => {
    const boardWithEmptyUrls = {
      ...mockBoard,
      imageUrl: '',
      keys: [{ ...mockKeys[0], imageUrl: '', audioUrl: '' }],
    }
    const result = extractMediaUrls(boardWithEmptyUrls)

    result.images.forEach(url => {
      expect(url).not.toBe('')
    })
    result.audios.forEach(url => {
      expect(url).not.toBe('')
    })
  })
})
