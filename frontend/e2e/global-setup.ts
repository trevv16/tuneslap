import type { FullConfig } from '@playwright/test'

const BACKEND_URL = 'http://localhost:8082'
const MAX_RETRIES = 30
const RETRY_DELAY = 2000

/**
 * Wait for the backend API to be healthy before running E2E tests.
 * This ensures the full stack is ready for true end-to-end testing.
 */
async function waitForBackend(): Promise<void> {
  console.log(`Checking backend health at ${BACKEND_URL}/health...`)

  for (let i = 0; i < MAX_RETRIES; i++) {
    try {
      const response = await fetch(`${BACKEND_URL}/health`)
      if (response.ok) {
        console.log('Backend is ready!')
        return
      }
      console.log(`Backend returned status ${response.status}, retrying...`)
    } catch (error) {
      // Backend not ready yet - this is expected during startup
      if (i === 0) {
        console.log('Waiting for backend to start...')
      }
    }

    if (i < MAX_RETRIES - 1) {
      console.log(`Waiting for backend... (${i + 1}/${MAX_RETRIES})`)
      await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY))
    }
  }

  throw new Error(
    `Backend did not become ready in time (${MAX_RETRIES * RETRY_DELAY / 1000}s). ` +
    'Make sure the backend is running with: ./scripts/e2e-setup.sh'
  )
}

/**
 * Global setup for Playwright E2E tests.
 * Runs once before all tests to ensure the test environment is ready.
 */
export default async function globalSetup(config: FullConfig): Promise<void> {
  console.log('')
  console.log('======================================')
  console.log('E2E Global Setup')
  console.log('======================================')

  await waitForBackend()

  console.log('')
  console.log('E2E Global Setup Complete!')
  console.log('======================================')
  console.log('')
}
