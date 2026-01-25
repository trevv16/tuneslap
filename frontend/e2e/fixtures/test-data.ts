/**
 * Reusable test data for Playwright e2e tests.
 *
 * Note: Once PR #36 (unit testing infrastructure) is merged, this file
 * can be refactored to import from __mocks__/data/fixtures.ts to avoid
 * duplication. For now, we define the data here for standalone usage.
 */

// Test user credentials
export const testUser = {
  id: 'user-1',
  email: 'test@example.com',
  password: 'testpassword123',
  name: 'Test User',
}

export const testUser2 = {
  id: 'user-2',
  email: 'user2@example.com',
  password: 'password456',
  name: 'Second User',
}

// Type definitions
interface MockKey {
  id: string
  name: string
  hotKey: string
  audioUrl: string
  imageUrl?: string
  boardId: string
  createdAt: string
  updatedAt: string
}

// Key fixtures
export const mockKey: MockKey = {
  id: 'key-1',
  name: 'Test Key',
  hotKey: 'A',
  audioUrl: 'https://example.com/audio/test.mp3',
  imageUrl: 'https://example.com/images/test.png',
  boardId: 'board-1',
  createdAt: '2024-01-01T00:00:00.000Z',
  updatedAt: '2024-01-01T00:00:00.000Z',
}

export const mockKeys: MockKey[] = [
  mockKey,
  {
    id: 'key-2',
    name: 'Second Key',
    hotKey: 'B',
    audioUrl: 'https://example.com/audio/test2.mp3',
    imageUrl: 'https://example.com/images/test2.png',
    boardId: 'board-1',
    createdAt: '2024-01-02T00:00:00.000Z',
    updatedAt: '2024-01-02T00:00:00.000Z',
  },
  {
    id: 'key-3',
    name: 'Third Key',
    hotKey: 'C',
    audioUrl: 'https://example.com/audio/test3.mp3',
    imageUrl: undefined,
    boardId: 'board-1',
    createdAt: '2024-01-03T00:00:00.000Z',
    updatedAt: '2024-01-03T00:00:00.000Z',
  },
]

interface MockCollaborator {
  id: string
  userId: string
  boardId: string
  role: string
  user: {
    id: string
    email: string
    name: string
  }
  createdAt: string
  updatedAt: string
}

interface MockBoard {
  id: string
  name: string
  description: string
  imageUrl?: string
  layout: string
  authorId: string
  keys: MockKey[]
  collaborators: MockCollaborator[]
  createdAt: string
  updatedAt: string
}

// Collaborator fixtures
export const mockCollaborator: MockCollaborator = {
  id: 'collab-1',
  userId: 'user-2',
  boardId: 'board-1',
  role: 'editor',
  user: {
    id: 'user-2',
    email: 'user2@example.com',
    name: 'Second User',
  },
  createdAt: '2024-01-01T00:00:00.000Z',
  updatedAt: '2024-01-01T00:00:00.000Z',
}

// Board fixtures
export const mockBoard: MockBoard = {
  id: 'board-1',
  name: 'Test Board',
  description: 'A test soundboard',
  imageUrl: 'https://example.com/images/board.png',
  layout: 'grid',
  authorId: 'user-1',
  keys: mockKeys,
  collaborators: [mockCollaborator],
  createdAt: '2024-01-01T00:00:00.000Z',
  updatedAt: '2024-01-01T00:00:00.000Z',
}

export const mockBoards: MockBoard[] = [
  mockBoard,
  {
    id: 'board-2',
    name: 'Second Board',
    description: 'Another test soundboard',
    imageUrl: undefined,
    layout: 'list',
    authorId: 'user-1',
    keys: [],
    collaborators: [],
    createdAt: '2024-01-02T00:00:00.000Z',
    updatedAt: '2024-01-02T00:00:00.000Z',
  },
]

export const mockEmptyBoard: MockBoard = {
  id: 'board-3',
  name: 'Empty Board',
  description: 'A board with no keys',
  imageUrl: undefined,
  layout: 'grid',
  authorId: 'user-1',
  keys: [],
  collaborators: [],
  createdAt: '2024-01-03T00:00:00.000Z',
  updatedAt: '2024-01-03T00:00:00.000Z',
}

// Media fixtures
export const mockAudioMedia = {
  id: 'media-1',
  fileName: 'test-audio.mp3',
  mediaType: 'audio',
  fileUrl: 'https://example.com/audio/test.mp3',
  processedUrl: 'https://example.com/audio/test-processed.webm',
  fileSize: 1024 * 1024,
  contentType: 'audio/mpeg',
  status: 'done',
  duration: 5.5,
  createdAt: '2024-01-01T00:00:00.000Z',
  updatedAt: '2024-01-01T00:00:00.000Z',
}

export const mockImageMedia = {
  id: 'media-2',
  fileName: 'test-image.png',
  mediaType: 'image',
  fileUrl: 'https://example.com/images/test.png',
  processedUrl: 'https://example.com/images/test-processed.webp',
  fileSize: 512 * 1024,
  contentType: 'image/png',
  status: 'done',
  width: 800,
  height: 600,
  createdAt: '2024-01-01T00:00:00.000Z',
  updatedAt: '2024-01-01T00:00:00.000Z',
}

export const mockMedia = [mockAudioMedia, mockImageMedia]

// API response fixtures
export const mockUserResponse = {
  id: testUser.id,
  email: testUser.email,
  name: testUser.name,
  createdAt: '2024-01-01T00:00:00.000Z',
  updatedAt: '2024-01-01T00:00:00.000Z',
}

export const mockSigninResponse = {
  success: true,
  message: 'Login successful',
  data: {
    token: 'mock-jwt-token-for-testing',
    user: mockUserResponse,
  },
}

export const mockSignupResponse = {
  success: true,
  message: 'Registration successful',
}

export const mockBoardsListResponse = mockBoards

export const mockMediaListResponse = mockMedia

// Error response fixtures
export const mockErrorResponse = {
  success: false,
  message: 'An error occurred',
  error: 'INTERNAL_ERROR',
}

export const mockValidationErrorResponse = {
  success: false,
  message: 'Validation failed',
  error: 'VALIDATION_ERROR',
  details: [{ field: 'email', message: 'Invalid email format' }],
}

export const mockUnauthorizedResponse = {
  success: false,
  message: 'Unauthorized',
  error: 'UNAUTHORIZED',
}
