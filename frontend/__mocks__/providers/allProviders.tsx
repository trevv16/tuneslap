// Combined providers wrapper for testing
import type { QueryClient } from '@tanstack/react-query'
import type { UserResponse } from '@/api/models'
import React from 'react'
import { MockAuthProvider } from './authContext'
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

export function AllProviders({ children, queryClient, authState }: AllProvidersProps) {
  const client = queryClient ?? createTestQueryClient()
  
  return (
    <QueryClientWrapper client={client}>
      <MockAuthProvider value={authState}>
        {children}
      </MockAuthProvider>
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
