// Combined providers wrapper for testing
import type { QueryClient } from '@tanstack/react-query'
import React from 'react'
import { MockAuthProvider } from './authContext'
import { createTestQueryClient, QueryClientWrapper } from './queryClient'

interface AllProvidersProps {
  children: React.ReactNode
  queryClient?: QueryClient
  authState?: {
    isLoading?: boolean
    isAuthenticated?: boolean
    user?: unknown
  }
}

export function AllProviders({ children, queryClient, authState }: AllProvidersProps) {
  const client = queryClient ?? createTestQueryClient()
  
  return React.createElement(
    QueryClientWrapper,
    { client },
    React.createElement(MockAuthProvider, { value: authState }, children)
  )
}

// Custom render function with all providers
interface RenderOptions {
  queryClient?: QueryClient
  authState?: {
    isLoading?: boolean
    isAuthenticated?: boolean
    user?: unknown
  }
}

export function createWrapper(options: RenderOptions = {}) {
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(AllProviders, options, children)
  }
}
