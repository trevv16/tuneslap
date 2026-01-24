// Note: The runtime module is a stub that allows tests to run without generated code

// Mock the token utility
jest.mock('@/utils/token', () => ({
  getStoredToken: jest.fn(),
}))

import { getApiConfig, getServerApiConfig } from './config'
import { getStoredToken } from '@/utils/token'

describe('getApiConfig', () => {
  const originalEnv = process.env

  beforeEach(() => {
    jest.resetModules()
    process.env = { ...originalEnv }
    ;(getStoredToken as jest.Mock).mockReset()
  })

  afterAll(() => {
    process.env = originalEnv
  })

  it('should return configuration with basePath from env', () => {
    process.env.NEXT_PUBLIC_API_URL = 'https://api.test.com/api/v1'
    ;(getStoredToken as jest.Mock).mockReturnValue('test-token')

    const config = getApiConfig()

    expect(config.basePath).toBe('https://api.test.com/api/v1')
  })

  it('should use production URL when env is not set', () => {
    delete process.env.NEXT_PUBLIC_API_URL
    ;(getStoredToken as jest.Mock).mockReturnValue('test-token')

    const config = getApiConfig()

    expect(config.basePath).toBe('https://api.tuneslap.com/api/v1')
  })

  it('should return accessToken function that gets token', () => {
    ;(getStoredToken as jest.Mock).mockReturnValue('my-jwt-token')

    const config = getApiConfig()
    const token = config.accessToken?.()

    expect(token).toBe('my-jwt-token')
    expect(getStoredToken).toHaveBeenCalled()
  })

  it('should return empty string when no token exists', () => {
    ;(getStoredToken as jest.Mock).mockReturnValue(null)

    const config = getApiConfig()
    const token = config.accessToken?.()

    expect(token).toBe('')
  })
})

describe('getServerApiConfig', () => {
  const originalEnv = process.env

  beforeEach(() => {
    jest.resetModules()
    process.env = { ...originalEnv }
  })

  afterAll(() => {
    process.env = originalEnv
  })

  it('should prefer INTERNAL_API_URL over NEXT_PUBLIC_API_URL', () => {
    process.env.INTERNAL_API_URL = 'http://server:8082/api/v1'
    process.env.NEXT_PUBLIC_API_URL = 'https://api.test.com/api/v1'

    const config = getServerApiConfig()

    expect(config.basePath).toBe('http://server:8082/api/v1')
  })

  it('should fallback to NEXT_PUBLIC_API_URL when INTERNAL_API_URL not set', () => {
    delete process.env.INTERNAL_API_URL
    process.env.NEXT_PUBLIC_API_URL = 'https://api.test.com/api/v1'

    const config = getServerApiConfig()

    expect(config.basePath).toBe('https://api.test.com/api/v1')
  })

  it('should use production URL when no env vars set', () => {
    delete process.env.INTERNAL_API_URL
    delete process.env.NEXT_PUBLIC_API_URL

    const config = getServerApiConfig()

    expect(config.basePath).toBe('https://api.tuneslap.com/api/v1')
  })

  it('should not include accessToken for server config', () => {
    const config = getServerApiConfig()

    expect(config.accessToken).toBeUndefined()
  })
})
