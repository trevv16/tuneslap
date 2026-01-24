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
  thresholds: number[] = []
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

// Mock Image as a jest.fn() so tests can mockImplementation
const createImageMock = jest.fn(() => {
  const img = {
    onload: null as (() => void) | null,
    onerror: null as (() => void) | null,
    src: '',
  }
  return img
})
Object.defineProperty(window, 'Image', { value: createImageMock, writable: true, configurable: true })

// Mock Audio as a jest.fn() so tests can mockImplementation
const createAudioMock = jest.fn(() => {
  const audio = {
    oncanplaythrough: null as (() => void) | null,
    onerror: null as (() => void) | null,
    src: '',
    load: jest.fn(),
  }
  return audio
})
Object.defineProperty(window, 'Audio', { value: createAudioMock, writable: true, configurable: true })

// Reset mocks before each test
beforeEach(() => {
  jest.clearAllMocks()
  localStorageMock.getItem.mockReset()
  localStorageMock.setItem.mockReset()
  localStorageMock.removeItem.mockReset()
})
