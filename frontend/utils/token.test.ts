import { getStoredToken, setStoredToken, removeStoredToken } from './token'

describe('token utilities', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear()
  })

  describe('getStoredToken', () => {
    it('should return token from localStorage', () => {
      const mockToken = 'test-jwt-token'
      localStorage.setItem('tuneslap_api_token', mockToken)

      const result = getStoredToken()

      expect(result).toBe(mockToken)
    })

    it('should return null when no token exists', () => {
      const result = getStoredToken()

      expect(result).toBeNull()
    })
  })

  describe('setStoredToken', () => {
    it('should store token in localStorage', () => {
      const token = 'new-jwt-token'

      setStoredToken(token)

      expect(localStorage.getItem('tuneslap_api_token')).toBe(token)
    })

    it('should handle empty token', () => {
      setStoredToken('')

      expect(localStorage.getItem('tuneslap_api_token')).toBe('')
    })
  })

  describe('removeStoredToken', () => {
    it('should remove token from localStorage', () => {
      localStorage.setItem('tuneslap_api_token', 'some-token')

      removeStoredToken()

      expect(localStorage.getItem('tuneslap_api_token')).toBeNull()
    })
  })
})

describe('token utilities - SSR handling', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('getStoredToken should handle being called (SSR check is in the function)', () => {
    const result = getStoredToken()
    expect(result).toBeNull()
  })

  it('setStoredToken should handle being called safely', () => {
    expect(() => {
      setStoredToken('token')
    }).not.toThrow()
  })

  it('removeStoredToken should handle being called safely', () => {
    expect(() => {
      removeStoredToken()
    }).not.toThrow()
  })
})
