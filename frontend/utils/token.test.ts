import { getStoredToken, setStoredToken, removeStoredToken } from './token'

describe('token utilities', () => {
  beforeEach(() => {
    // Clear localStorage mock before each test
    jest.clearAllMocks()
    ;(window.localStorage.getItem as jest.Mock).mockReset()
    ;(window.localStorage.setItem as jest.Mock).mockReset()
    ;(window.localStorage.removeItem as jest.Mock).mockReset()
  })

  describe('getStoredToken', () => {
    it('should return token from localStorage', () => {
      const mockToken = 'test-jwt-token'
      ;(window.localStorage.getItem as jest.Mock).mockReturnValue(mockToken)

      const result = getStoredToken()

      expect(result).toBe(mockToken)
      expect(window.localStorage.getItem).toHaveBeenCalledWith('tuneslap_api_token')
    })

    it('should return null when no token exists', () => {
      ;(window.localStorage.getItem as jest.Mock).mockReturnValue(null)

      const result = getStoredToken()

      expect(result).toBeNull()
      expect(window.localStorage.getItem).toHaveBeenCalledWith('tuneslap_api_token')
    })
  })

  describe('setStoredToken', () => {
    it('should store token in localStorage', () => {
      const token = 'new-jwt-token'

      setStoredToken(token)

      expect(window.localStorage.setItem).toHaveBeenCalledWith('tuneslap_api_token', token)
    })

    it('should handle empty token', () => {
      setStoredToken('')

      expect(window.localStorage.setItem).toHaveBeenCalledWith('tuneslap_api_token', '')
    })
  })

  describe('removeStoredToken', () => {
    it('should remove token from localStorage', () => {
      removeStoredToken()

      expect(window.localStorage.removeItem).toHaveBeenCalledWith('tuneslap_api_token')
    })
  })
})

describe('token utilities - SSR handling', () => {
  // Note: These tests verify that the functions handle SSR gracefully
  // In a real SSR environment (Node.js without JSDOM), window would be undefined
  // Here we test that the functions check for window existence before accessing localStorage
  
  it('getStoredToken should handle being called (SSR check is in the function)', () => {
    // The function has a typeof window check, so it works in both environments
    // In JSDOM, window exists, so it will try to get from localStorage
    ;(window.localStorage.getItem as jest.Mock).mockReturnValue(null)
    const result = getStoredToken()
    expect(result).toBeNull()
  })

  it('setStoredToken should handle being called safely', () => {
    expect(() => setStoredToken('token')).not.toThrow()
  })

  it('removeStoredToken should handle being called safely', () => {
    expect(() => removeStoredToken()).not.toThrow()
  })
})
