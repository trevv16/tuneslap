import '@testing-library/jest-dom'

// Set test API URL for MSW to intercept
process.env.NEXT_PUBLIC_API_URL = 'http://localhost:8082/api/v1'

// Polyfill for Fetch API globals required by MSW v2
// jsdom doesn't provide native fetch globals, but Node 18+ does
import { TextDecoder, TextEncoder } from 'util'
import { Blob } from 'buffer'
import { ReadableStream, TransformStream, WritableStream } from 'stream/web'
import { MessageChannel, MessagePort } from 'worker_threads'
import { BroadcastChannel } from 'worker_threads'

Object.assign(global, { 
  TextDecoder, 
  TextEncoder,
  Blob,
  ReadableStream,
  TransformStream,
  WritableStream,
  MessageChannel,
  MessagePort,
  BroadcastChannel,
})

// Import fetch globals from undici (bundled with Node 18+)
// eslint-disable-next-line @typescript-eslint/no-require-imports
const { fetch, Request, Response, Headers, FormData } = require('undici')
Object.assign(globalThis, { fetch, Request, Response, Headers, FormData })

// Mock localStorage with actual storage functionality
class LocalStorageMock implements Storage {
  private store: Record<string, string> = {}
  
  get length(): number {
    return Object.keys(this.store).length
  }
  
  clear(): void {
    this.store = {}
  }
  
  getItem(key: string): string | null {
    return this.store[key] ?? null
  }
  
  key(index: number): string | null {
    const keys = Object.keys(this.store)
    return keys[index] ?? null
  }
  
  removeItem(key: string): void {
    delete this.store[key]
  }
  
  setItem(key: string, value: string): void {
    this.store[key] = value
  }
}

const localStorageMock = new LocalStorageMock()
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

// Avoid mocking fetch globally; MSW v2 relies on undici fetch above.
// Tests needing fetch assertion can override in their own setup.

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
  localStorageMock.clear()
})

// MSW server setup for integration tests
// Only load if the server file exists (avoids errors for unit tests)
let server: { listen: (opts?: object) => void; resetHandlers: () => void; close: () => void } | undefined

try {
  // Dynamic import to avoid errors when running unit tests
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const mswSetup = require('./__tests__/integration/setup/msw-server')
  server = mswSetup.server
} catch {
  // MSW server not available, running unit tests only
}

if (server) {
  beforeAll(() => server?.listen({ onUnhandledRequest: 'error' }))
  afterEach(() => server?.resetHandlers())
  afterAll(() => server?.close())
}
