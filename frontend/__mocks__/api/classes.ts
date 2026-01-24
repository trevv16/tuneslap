// Mock API classes for testing

export const mockAuthApi = {
  signup: jest.fn(),
  signin: jest.fn(),
  forgot: jest.fn(),
  reset: jest.fn(),
}

export const mockBoardsApi = {
  getAllBoards: jest.fn(),
  getBoardById: jest.fn(),
  createBoard: jest.fn(),
  updateBoard: jest.fn(),
  deleteBoard: jest.fn(),
}

export const mockUsersApi = {
  getMe: jest.fn(),
  updateUser: jest.fn(),
  deleteUser: jest.fn(),
  changePassword: jest.fn(),
}

export const mockMediaApi = {
  getAllMedia: jest.fn(),
  getMediaById: jest.fn(),
  createMedia: jest.fn(),
  updateMedia: jest.fn(),
  deleteMedia: jest.fn(),
  processMedia: jest.fn(),
}

export const mockKeysApi = {
  getAllKeys: jest.fn(),
  getKeyById: jest.fn(),
  createKey: jest.fn(),
  updateKey: jest.fn(),
  deleteKey: jest.fn(),
}

export const mockCollaboratorsApi = {
  getAllCollaborators: jest.fn(),
  addCollaborator: jest.fn(),
  removeCollaborator: jest.fn(),
  updateCollaborator: jest.fn(),
}

// Factory functions to create mock API instances
export const AuthApi = jest.fn(() => mockAuthApi)
export const BoardsApi = jest.fn(() => mockBoardsApi)
export const UsersApi = jest.fn(() => mockUsersApi)
export const MediaApi = jest.fn(() => mockMediaApi)
export const KeysApi = jest.fn(() => mockKeysApi)
export const CollaboratorsApi = jest.fn(() => mockCollaboratorsApi)

// Reset all mocks
export function resetApiMocks(): void {
  Object.values(mockAuthApi).forEach(fn => fn.mockReset())
  Object.values(mockBoardsApi).forEach(fn => fn.mockReset())
  Object.values(mockUsersApi).forEach(fn => fn.mockReset())
  Object.values(mockMediaApi).forEach(fn => fn.mockReset())
  Object.values(mockKeysApi).forEach(fn => fn.mockReset())
  Object.values(mockCollaboratorsApi).forEach(fn => fn.mockReset())
}
