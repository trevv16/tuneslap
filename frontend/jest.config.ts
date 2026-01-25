import type { Config } from 'jest'

const config: Config = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/$1',
  },
  testPathIgnorePatterns: [
    '<rootDir>/node_modules/',
    '<rootDir>/.next/',
    '<rootDir>/__tests__/integration/setup/',
    '<rootDir>/e2e/',
  ],
  transform: {
    '^.+\\.(ts|tsx)$': ['ts-jest', {
      tsconfig: {
        jsx: 'react-jsx',
        esModuleInterop: true,
        module: 'commonjs',
        moduleResolution: 'node',
      },
    }],
    // Transform ESM dependencies from node_modules
    '^.+\\.js$': ['ts-jest', {
      tsconfig: {
        allowJs: true,
        esModuleInterop: true,
        module: 'commonjs',
      },
    }],
  },
  // Transform MSW v2 and its ESM dependencies
  transformIgnorePatterns: [
    'node_modules/\\.pnpm',
    'node_modules/(?!(@?msw|@mswjs|until-async|@bundled-es-modules|@open-draft|outvariant|strict-event-emitter|path-to-regexp|headers-polyfill|cookie|statuses|graphql|is-node-process|type-fest)/)',
  ],
  collectCoverageFrom: [
    'utils/**/*.{ts,tsx}',
    'hooks/**/*.{ts,tsx}',
    'contexts/**/*.{ts,tsx}',
    'components/**/*.{ts,tsx}',
    'api/config.ts',
    'api/uploadUrl.ts',
    '!**/*.d.ts',
    '!**/node_modules/**',
    '!components/ui/**',
  ],
  coverageDirectory: 'coverage',
  coverageReporters: ['text', 'lcov', 'html'],
}

export default config
