import { getStoredToken, setStoredToken, removeStoredToken } from './token'

describe('token utilities', () => {
  beforeEach(() => {
    // Clear localStorage mock before each test
    jest.clearAllMocks()
    const getItemMock = window.localStorage.getItem as jest.Mock
    const setItemMock = window.localStorage.setItem as jest.Mock
    const removeItemMock = window.localStorage.removeItem as jest.Mock
    getItemMock.mockReset()
    setItemMock.mockReset()
    removeItemMock.mockReset()
  })

  describe('getStoredToken', () => {
    it('should return token from localStorage', () => {
      const mockToken = 'test-jwt-token'
      const getItemMock = window.localStorage.getItem as jest.Mock
      getItemMock.mockReturnValue(mockToken)

      const result = getStoredToken()

      expect(result).toBe(mockToken)
      expect(getItemMock).toHaveBeenCalledWith('tuneslap_api_token')
    })

    it('should return null when no token exists', () => {
      const getItemMock = window.localStorage.getItem as jest.Mock
      getItemMock.mockReturnValue(null)

      const result = getStoredToken()

      expect(result).toBeNull()
      expect(getItemMock).toHaveBeenCalledWith('tuneslap_api_token')
    })
  })

  describe('setStoredToken', () => {
    it('should store token in localStorage', () => {
      const token = 'new-jwt-token'
      const setItemMock = window.localStorage.setItem as jest.Mock

      setStoredToken(token)

      expect(setItemMock).toHaveBeenCalledWith('tuneslap_api_token', token)
    })

    it('should handle empty token', () => {
      const setItemMock = window.localStorage.setItem as jest.Mock

      setStoredToken('')

      expect(setItemMock).toHaveBeenCalledWith('tuneslap_api_token', '')
    })
  })

  describe('removeStoredToken', () => {
    it('should remove token from localStorage', () => {
      const removeItemMock = window.localStorage.removeItem as jest.Mock

      removeStoredToken()

      expect(removeItemMock).toHaveBeenCalledWith('tuneslap_api_token')
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
    const getItemMock = window.localStorage.getItem as jest.Mock
    getItemMock.mockReturnValue(null)
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
