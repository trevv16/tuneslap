import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { AuthContextProvider, useAuthContext } from './AuthContext'
import { mockUser } from '@/__mocks__/data/fixtures'
import React from 'react'

// Mock dependencies
jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}))

jest.mock('@/utils/token', () => ({
  getStoredToken: jest.fn(),
  removeStoredToken: jest.fn(),
}))

jest.mock('../hooks/users', () => ({
  useGetMe: jest.fn(),
}))

import { useRouter } from 'next/navigation'
import { getStoredToken, removeStoredToken } from '@/utils/token'
import { useGetMe } from '../hooks/users'

// Test component that uses the auth context
function TestConsumer() {
  const { isLoading, isAuthenticated, user, signOut } = useAuthContext()
  
  return (
    <div>
      <div data-testid="loading">{isLoading.toString()}</div>
      <div data-testid="authenticated">{isAuthenticated.toString()}</div>
      <div data-testid="user">{user?.name ?? 'null'}</div>
      <button onClick={signOut} data-testid="signout">Sign Out</button>
    </div>
  )
}

describe('AuthContextProvider', () => {
  const mockRouter = {
    push: jest.fn(),
  }

  beforeEach(() => {
    jest.clearAllMocks()
    const mockUseRouter = useRouter as jest.Mock
    mockUseRouter.mockReturnValue(mockRouter)
    const mockGetStoredToken = getStoredToken as jest.Mock
    mockGetStoredToken.mockReturnValue(null)
    const mockUseGetMe = useGetMe as jest.Mock
    mockUseGetMe.mockReturnValue({
      data: undefined,
      error: undefined,
    })
  })

  it('should show loading state when token exists but user not resolved', () => {
    const mockGetStoredToken = getStoredToken as jest.Mock
    mockGetStoredToken.mockReturnValue('mock-token')
    const mockUseGetMe = useGetMe as jest.Mock
    mockUseGetMe.mockReturnValue({
      data: undefined,
      error: undefined,
    })

    render(
      <AuthContextProvider>
        <TestConsumer />
      </AuthContextProvider>
    )

    expect(screen.getByTestId('loading')).toHaveTextContent('true')
    expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
  })

  it('should be authenticated when token and user exist', async () => {
    const mockGetStoredToken = getStoredToken as jest.Mock
    mockGetStoredToken.mockReturnValue('mock-token')
    const mockUseGetMe = useGetMe as jest.Mock
    mockUseGetMe.mockReturnValue({
      data: { data: mockUser },
      error: undefined,
    })

    render(
      <AuthContextProvider>
        <TestConsumer />
      </AuthContextProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('false')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
      expect(screen.getByTestId('user')).toHaveTextContent(mockUser.name ?? '')
    })
  })

  it('should not be authenticated when no token', () => {
    const mockGetStoredToken = getStoredToken as jest.Mock
    mockGetStoredToken.mockReturnValue(null)

    render(
      <AuthContextProvider>
        <TestConsumer />
      </AuthContextProvider>
    )

    expect(screen.getByTestId('loading')).toHaveTextContent('false')
    expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
    expect(screen.getByTestId('user')).toHaveTextContent('null')
  })

  it('should clear user on error', async () => {
    const mockGetStoredToken = getStoredToken as jest.Mock
    mockGetStoredToken.mockReturnValue('mock-token')
    const mockUseGetMe = useGetMe as jest.Mock
    mockUseGetMe.mockReturnValue({
      data: undefined,
      error: new Error('Unauthorized'),
    })

    render(
      <AuthContextProvider>
        <TestConsumer />
      </AuthContextProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('null')
    })
  })

  it('should handle signOut correctly', async () => {
    const user = userEvent.setup()
    const mockGetStoredToken = getStoredToken as jest.Mock
    mockGetStoredToken.mockReturnValue('mock-token')
    const mockUseGetMe = useGetMe as jest.Mock
    mockUseGetMe.mockReturnValue({
      data: { data: mockUser },
      error: undefined,
    })

    render(
      <AuthContextProvider>
        <TestConsumer />
      </AuthContextProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
    })

    await user.click(screen.getByTestId('signout'))

    expect(removeStoredToken).toHaveBeenCalled()
    expect(mockRouter.push).toHaveBeenCalledWith('/auth/signin')
  })
})

describe('useAuthContext', () => {
  it('should throw error when used outside provider', () => {
    // Suppress console.error for this test
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation()

    function TestComponent() {
      useAuthContext()
      return null
    }

    expect(() => render(<TestComponent />)).toThrow(
      'Must be a child of AuthContextProvider'
    )

    consoleSpy.mockRestore()
  })
})
