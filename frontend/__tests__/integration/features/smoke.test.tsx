// Smoke test to verify integration test infrastructure works
// This test validates MSW is properly intercepting requests
import { renderWithProviders, screen, waitFor, mockUser } from '../setup/test-utils'

// Simple component that displays user info
function TestComponent({ showUser }: { showUser: boolean }) {
  return (
    <div>
      <h1>Integration Test Smoke Test</h1>
      {showUser && <p data-testid="user-name">{mockUser.name}</p>}
    </div>
  )
}

describe('Integration Test Infrastructure', () => {
  it('renders components with providers', () => {
    renderWithProviders(<TestComponent showUser={false} />)
    
    expect(screen.getByRole('heading', { name: /integration test smoke test/i })).toBeInTheDocument()
  })

  it('provides mock data from fixtures', () => {
    renderWithProviders(<TestComponent showUser={true} />)
    
    expect(screen.getByTestId('user-name')).toHaveTextContent(mockUser.name)
  })

  it('supports authenticated state', () => {
    renderWithProviders(<TestComponent showUser={true} />, {
      authState: { isAuthenticated: true, user: mockUser },
    })
    
    expect(screen.getByTestId('user-name')).toBeInTheDocument()
  })

  it('supports async waitFor assertions', async () => {
    renderWithProviders(<TestComponent showUser={true} />)
    
    await waitFor(() => {
      expect(screen.getByTestId('user-name')).toBeInTheDocument()
    })
  })
})
