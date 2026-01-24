// Combined providers wrapper for testing
import type { QueryClient } from '@tanstack/react-query'
import type { UserResponse } from '@/api/models'
import React from 'react'
import { AuthContext } from '@/contexts/AuthContext'
import { createTestQueryClient, QueryClientWrapper } from './queryClient'

interface AuthState {
  isLoading?: boolean
  isAuthenticated?: boolean
  user?: UserResponse | null
}

interface AllProvidersProps {
  children: React.ReactNode
  queryClient?: QueryClient
  authState?: AuthState
}

const defaultUser: UserResponse = {
  id: 'user-1',
  email: 'test@example.com',
  name: 'Test User',
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export function AllProviders({ children, queryClient, authState }: AllProvidersProps) {
  const client = queryClient ?? createTestQueryClient()
  
  // Use the real AuthContext with mock values
  const authValue = {
    isLoading: authState?.isLoading ?? false,
    isAuthenticated: authState?.isAuthenticated ?? true,
    user: authState?.user ?? defaultUser,
    setUser: jest.fn(),
    signOut: jest.fn(),
  }
  
  return (
    <QueryClientWrapper client={client}>
      <AuthContext.Provider value={authValue}>
        {children}
      </AuthContext.Provider>
    </QueryClientWrapper>
  )
}

// Custom render function with all providers
interface RenderOptions {
  queryClient?: QueryClient
  authState?: AuthState
}

export function createWrapper(options: RenderOptions = {}) {
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return (
      <AllProviders {...options}>
        {children}
      </AllProviders>
    )
  }
}
