import type { Config } from 'jest'
import path from 'path'

const rootDir = path.resolve(__dirname)

const config: Config = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: [path.join(rootDir, 'jest.setup.ts')],
  moduleNameMapper: {
    '^@/(.*)$': path.join(rootDir, '$1'),
  },
  testPathIgnorePatterns: [
    path.join(rootDir, 'node_modules'),
    path.join(rootDir, '.next'),
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
  },
  transformIgnorePatterns: [
    'node_modules/(?!(.*\\.mjs$))',
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
