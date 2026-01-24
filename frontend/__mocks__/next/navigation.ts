// Mock Next.js navigation hooks

export const mockRouter = {
  push: jest.fn(),
  replace: jest.fn(),
  back: jest.fn(),
  forward: jest.fn(),
  refresh: jest.fn(),
  prefetch: jest.fn(),
}

export const mockPathname = '/'
export const mockSearchParams = new URLSearchParams()

export const useRouter = jest.fn(() => mockRouter)
export const usePathname = jest.fn(() => mockPathname)
export const useSearchParams = jest.fn(() => mockSearchParams)
export const useParams = jest.fn(() => ({}))

// Helper to set custom values for tests
export function setMockPathname(pathname: string): void {
  usePathname.mockReturnValue(pathname)
}

export function setMockSearchParams(params: URLSearchParams): void {
  useSearchParams.mockReturnValue(params)
}

export function setMockParams(params: Record<string, string>): void {
  useParams.mockReturnValue(params)
}

export function resetNavigationMocks(): void {
  mockRouter.push.mockReset()
  mockRouter.replace.mockReset()
  mockRouter.back.mockReset()
  mockRouter.forward.mockReset()
  mockRouter.refresh.mockReset()
  mockRouter.prefetch.mockReset()
  useRouter.mockReturnValue(mockRouter)
  usePathname.mockReturnValue('/')
  useSearchParams.mockReturnValue(new URLSearchParams())
  useParams.mockReturnValue({})
}
