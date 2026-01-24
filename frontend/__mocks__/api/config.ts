// Mock API configuration

export const mockConfiguration = {
  basePath: 'http://localhost:8082/api/v1',
  accessToken: jest.fn(() => 'mock-token'),
}

export const getApiConfig = jest.fn(() => mockConfiguration)

export const getServerApiConfig = jest.fn(() => ({
  basePath: 'http://localhost:8082/api/v1',
}))
