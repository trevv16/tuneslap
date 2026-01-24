// Mock AuthContext for testing
import React from 'react'
import type { UserResponse } from '@/api/models'

interface MockAuthContextValue {
  isLoading: boolean
  isAuthenticated: boolean
  user: UserResponse | null
  setUser: React.Dispatch<React.SetStateAction<UserResponse | null>>
  signOut: () => void
}

const defaultMockAuthContext: MockAuthContextValue = {
  isLoading: false,
  isAuthenticated: true,
  user: {
    id: 'user-1',
    email: 'test@example.com',
    name: 'Test User',
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date('2024-01-01'),
  },
  setUser: jest.fn(),
  signOut: jest.fn(),
}

export const MockAuthContext = React.createContext<MockAuthContextValue | null>(null)

interface MockAuthProviderProps {
  children: React.ReactNode
  value?: Partial<MockAuthContextValue>
}

export function MockAuthProvider({ children, value = {} }: MockAuthProviderProps) {
  const contextValue: MockAuthContextValue = {
    ...defaultMockAuthContext,
    ...value,
  }
  return React.createElement(MockAuthContext.Provider, { value: contextValue }, children)
}

export function useMockAuthContext(): MockAuthContextValue {
  const context = React.useContext(MockAuthContext)
  if (!context) {
    throw new Error('useMockAuthContext must be used within MockAuthProvider')
  }
  return context
}

// Preset configurations for common test scenarios
export const mockAuthStates = {
  authenticated: {
    isLoading: false,
    isAuthenticated: true,
    user: defaultMockAuthContext.user,
  },
  unauthenticated: {
    isLoading: false,
    isAuthenticated: false,
    user: null,
  },
  loading: {
    isLoading: true,
    isAuthenticated: false,
    user: null,
  },
}
