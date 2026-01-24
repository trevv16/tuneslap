import '@testing-library/jest-dom'

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
  length: 0,
  key: jest.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
})

// Mock ResizeObserver
class ResizeObserverMock {
  observe = jest.fn()
  unobserve = jest.fn()
  disconnect = jest.fn()
}
Object.defineProperty(window, 'ResizeObserver', { value: ResizeObserverMock })

// Mock IntersectionObserver
class IntersectionObserverMock {
  observe = jest.fn()
  unobserve = jest.fn()
  disconnect = jest.fn()
  root = null
  rootMargin = ''
  thresholds = []
}
Object.defineProperty(window, 'IntersectionObserver', { value: IntersectionObserverMock })

// Mock fetch
global.fetch = jest.fn()

// Mock AudioContext
class AudioContextMock {
  createBufferSource = jest.fn(() => ({
    buffer: null,
    connect: jest.fn(),
    start: jest.fn(),
    stop: jest.fn(),
    disconnect: jest.fn(),
  }))
  decodeAudioData = jest.fn()
  destination = {}
  close = jest.fn(() => Promise.resolve())
}
Object.defineProperty(window, 'AudioContext', { value: AudioContextMock, writable: true, configurable: true })

// Mock Image
class ImageMock {
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  src = ''
}
Object.defineProperty(window, 'Image', { value: ImageMock, writable: true, configurable: true })

// Mock Audio
class AudioMock {
  oncanplaythrough: (() => void) | null = null
  onerror: (() => void) | null = null
  src = ''
  load = jest.fn()
}
Object.defineProperty(window, 'Audio', { value: AudioMock, writable: true, configurable: true })

// Reset mocks before each test
beforeEach(() => {
  jest.clearAllMocks()
  localStorageMock.getItem.mockReset()
  localStorageMock.setItem.mockReset()
  localStorageMock.removeItem.mockReset()
})
