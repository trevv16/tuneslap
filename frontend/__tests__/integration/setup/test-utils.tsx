// Integration test utilities
// Extends existing AllProviders from __mocks__/providers - DO NOT duplicate provider logic
import { render, RenderOptions, RenderResult } from '@testing-library/react'
import { QueryClient } from '@tanstack/react-query'
import type { UserResponse } from '@/api/models'
import React from 'react'
import { AllProviders } from '@/__mocks__/providers/allProviders'
import { createTestQueryClient } from '@/__mocks__/providers/queryClient'
import { Toaster } from 'sonner'

interface AuthState {
  isLoading?: boolean
  isAuthenticated?: boolean
  user?: UserResponse | null
}

interface IntegrationRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  queryClient?: QueryClient
  authState?: AuthState
}

// Extend AllProviders with Toaster for integration tests that need toast assertions
function IntegrationProviders({
  children,
  queryClient,
  authState,
}: {
  children: React.ReactNode
  queryClient?: QueryClient
  authState?: AuthState
}) {
  return (
    <AllProviders queryClient={queryClient} authState={authState}>
      {children}
      <Toaster data-testid="sonner-toaster" />
    </AllProviders>
  )
}

/**
 * Render function for integration tests with all providers and Toaster.
 * Use this instead of RTL's render() in integration tests.
 */
export function renderWithProviders(
  ui: React.ReactElement,
  options: IntegrationRenderOptions = {}
): RenderResult {
  const { queryClient, authState, ...renderOptions } = options
  const client = queryClient ?? createTestQueryClient()

  return render(ui, {
    wrapper: ({ children }) => (
      <IntegrationProviders queryClient={client} authState={authState}>
        {children}
      </IntegrationProviders>
    ),
    ...renderOptions,
  })
}

/**
 * Create a wrapper function for testing hooks that need providers.
 * Use with renderHook() from @testing-library/react.
 */
export function createIntegrationWrapper(options: IntegrationRenderOptions = {}) {
  const client = options.queryClient ?? createTestQueryClient()
  
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return (
      <IntegrationProviders queryClient={client} authState={options.authState}>
        {children}
      </IntegrationProviders>
    )
  }
}

// Re-export commonly used testing utilities
export { screen, waitFor, within, fireEvent } from '@testing-library/react'
export { default as userEvent } from '@testing-library/user-event'

// Re-export fixtures for convenience
export {
  mockUser,
  mockUsers,
  mockBoard,
  mockBoards,
  mockKey,
  mockKeys,
  mockMedia,
  mockAudioMedia,
  mockImageMedia,
  mockSigninResponse,
  mockSignupResponse,
  mockErrorResponse,
} from '@/__mocks__/data/fixtures'

// Re-export MSW server and error handlers
export { server } from './msw-server'
export { errorHandlers } from './msw-handlers'
