// Test QueryClient wrapper
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'

// Create a new QueryClient for each test to avoid shared state
export function createTestQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
        staleTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  })
}

interface QueryClientWrapperProps {
  children: React.ReactNode
  client?: QueryClient
}

export function QueryClientWrapper({ children, client }: QueryClientWrapperProps) {
  const queryClient = client ?? createTestQueryClient()
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}

// Helper to wrap a component with QueryClient for testing
export function withQueryClient(ui: React.ReactElement, client?: QueryClient): React.ReactElement {
  return (
    <QueryClientWrapper client={client}>
      {ui}
    </QueryClientWrapper>
  )
}
