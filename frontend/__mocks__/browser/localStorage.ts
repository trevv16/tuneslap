// localStorage mock utilities

interface StorageData {
  [key: string]: string
}

export function createLocalStorageMock(initialData: StorageData = {}) {
  let store: StorageData = { ...initialData }

  return {
    getItem: jest.fn((key: string) => store[key] ?? null),
    setItem: jest.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: jest.fn((key: string) => {
      delete store[key]
    }),
    clear: jest.fn(() => {
      store = {}
    }),
    get length() {
      return Object.keys(store).length
    },
    key: jest.fn((index: number) => {
      const keys = Object.keys(store)
      return keys[index] ?? null
    }),
    // Helper methods for testing
    _getStore: () => ({ ...store }),
    _setStore: (data: StorageData) => {
      store = { ...data }
    },
  }
}

// Default mock instance
export const localStorageMock = createLocalStorageMock()

// Helper to setup localStorage mock with initial data
export function setupLocalStorage(initialData: StorageData = {}): void {
  const mock = createLocalStorageMock(initialData)
  Object.defineProperty(window, 'localStorage', { value: mock, writable: true })
}

// Helper to reset localStorage mock
export function resetLocalStorage(): void {
  localStorageMock.clear()
  localStorageMock.getItem.mockClear()
  localStorageMock.setItem.mockClear()
  localStorageMock.removeItem.mockClear()
}
