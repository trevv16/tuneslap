// AudioContext mock utilities

interface AudioBufferSourceNodeMock {
  buffer: AudioBuffer | null
  connect: jest.Mock
  start: jest.Mock
  stop: jest.Mock
  disconnect: jest.Mock
  onended: (() => void) | null
}

interface AudioContextMockInstance {
  createBufferSource: jest.Mock<AudioBufferSourceNodeMock>
  decodeAudioData: jest.Mock
  destination: object
  close: jest.Mock
  state: 'running' | 'suspended' | 'closed'
  resume: jest.Mock
  suspend: jest.Mock
}

export function createAudioContextMock(): AudioContextMockInstance {
  const createBufferSource = jest.fn((): AudioBufferSourceNodeMock => ({
    buffer: null,
    connect: jest.fn(),
    start: jest.fn(),
    stop: jest.fn(),
    disconnect: jest.fn(),
    onended: null,
  }))

  return {
    createBufferSource,
    decodeAudioData: jest.fn().mockResolvedValue({
      duration: 1.0,
      length: 44100,
      numberOfChannels: 2,
      sampleRate: 44100,
    }),
    destination: {},
    close: jest.fn().mockResolvedValue(undefined),
    state: 'running',
    resume: jest.fn().mockResolvedValue(undefined),
    suspend: jest.fn().mockResolvedValue(undefined),
  }
}

// Global mock class
export class MockAudioContext implements AudioContextMockInstance {
  createBufferSource: jest.Mock<AudioBufferSourceNodeMock>
  decodeAudioData: jest.Mock
  destination: object
  close: jest.Mock
  state: 'running' | 'suspended' | 'closed'
  resume: jest.Mock
  suspend: jest.Mock

  constructor() {
    const mock = createAudioContextMock()
    this.createBufferSource = mock.createBufferSource
    this.decodeAudioData = mock.decodeAudioData
    this.destination = mock.destination
    this.close = mock.close
    this.state = mock.state
    this.resume = mock.resume
    this.suspend = mock.suspend
  }
}

// Helper to setup AudioContext mock
export function setupAudioContext(): void {
  Object.defineProperty(window, 'AudioContext', {
    value: MockAudioContext,
    writable: true,
  })
}

// Helper to create a mock AudioBuffer for testing
export function createMockAudioBuffer(): AudioBuffer {
  return {
    duration: 1.0,
    length: 44100,
    numberOfChannels: 2,
    sampleRate: 44100,
    getChannelData: jest.fn(),
    copyFromChannel: jest.fn(),
    copyToChannel: jest.fn(),
  } as unknown as AudioBuffer
}
