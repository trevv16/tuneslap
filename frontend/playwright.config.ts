import { defineConfig, devices } from '@playwright/test'

// E2E environment configuration - real API URLs for true end-to-end testing
const E2E_ENV = {
  NEXT_PUBLIC_API_URL: 'http://localhost:8082/api/v1',
  INTERNAL_API_URL: 'http://localhost:8082/api/v1',
  NEXT_PUBLIC_SITE_URL: 'http://localhost:3000',
  NEXT_PUBLIC_DEMO_MODE: 'true',
}

/**
 * Playwright configuration for TuneSlap frontend e2e tests.
 * See https://playwright.dev/docs/test-configuration
 */
export default defineConfig({
  // Test directory
  testDir: './e2e/tests',

  // Global setup to wait for backend health
  globalSetup: require.resolve('./e2e/global-setup.ts'),

  // Run tests in parallel
  fullyParallel: true,

  // Fail the build on CI if you accidentally left test.only in the source code
  forbidOnly: !!process.env.CI,

  // Retry on CI only
  retries: process.env.CI ? 2 : 0,

  // Limit parallel workers on CI
  workers: process.env.CI ? 1 : undefined,

  // Reporter configuration
  reporter: [
    ['html', { outputFolder: 'playwright-report' }],
    ['list'],
  ],

  // Shared settings for all projects
  use: {
    // Base URL for navigation
    baseURL: 'http://localhost:3000',

    // Collect trace when retrying a failed test
    trace: 'on-first-retry',

    // Screenshot on failure
    screenshot: 'only-on-failure',

    // Video on failure
    video: 'on-first-retry',
  },

  // Configure projects for major browsers
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
      },
    },

    {
      name: 'firefox',
      use: {
        ...devices['Desktop Firefox'],
      },
    },

    {
      name: 'webkit',
      use: {
        ...devices['Desktop Safari'],
      },
    },

    // Mobile viewports
    {
      name: 'mobile-chrome',
      use: {
        ...devices['Pixel 5'],
      },
    },

    {
      name: 'mobile-safari',
      use: {
        ...devices['iPhone 12'],
      },
    },
  ],

  // Run local dev server before starting the tests
  webServer: {
    command: 'npx next dev --webpack',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 60 * 1000, // 1 min timeout for dev server startup
    stdout: 'ignore', // Suppress verbose output
    stderr: 'pipe', // Show errors
    env: {
      ...process.env,
      ...E2E_ENV,
    },
  },

  // Global timeout for each test
  timeout: 30 * 1000,

  // Timeout for each action (click, fill, etc.)
  expect: {
    timeout: 5 * 1000,
  },
})
