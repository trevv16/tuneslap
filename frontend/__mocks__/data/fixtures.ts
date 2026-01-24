// Reusable test data fixtures
import type {
  BoardResponse,
  KeyResponse,
  UserResponse,
  MediaResponse,
  CollaboratorResponse,
} from '@/api/models'

// User fixtures
export const mockUser: UserResponse = {
  id: 'user-1',
  email: 'test@example.com',
  name: 'Test User',
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockUsers: UserResponse[] = [
  mockUser,
  {
    id: 'user-2',
    email: 'user2@example.com',
    name: 'Second User',
    createdAt: new Date('2024-01-02'),
    updatedAt: new Date('2024-01-02'),
  },
]

// Key fixtures
export const mockKey: KeyResponse = {
  id: 'key-1',
  name: 'Test Key',
  hotKey: 'A',
  audioUrl: 'https://example.com/audio/test.mp3',
  imageUrl: 'https://example.com/images/test.png',
  boardId: 'board-1',
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockKeys: KeyResponse[] = [
  mockKey,
  {
    id: 'key-2',
    name: 'Second Key',
    hotKey: 'B',
    audioUrl: 'https://example.com/audio/test2.mp3',
    imageUrl: 'https://example.com/images/test2.png',
    boardId: 'board-1',
    createdAt: new Date('2024-01-02'),
    updatedAt: new Date('2024-01-02'),
  },
  {
    id: 'key-3',
    name: 'Third Key',
    hotKey: 'C',
    audioUrl: 'https://example.com/audio/test3.mp3',
    imageUrl: undefined,
    boardId: 'board-1',
    createdAt: new Date('2024-01-03'),
    updatedAt: new Date('2024-01-03'),
  },
]

// Collaborator fixtures
export const mockCollaborator: CollaboratorResponse = {
  id: 'collab-1',
  userId: 'user-2',
  boardId: 'board-1',
  role: 'editor',
  user: mockUsers[1],
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockCollaborators: CollaboratorResponse[] = [
  mockCollaborator,
]

// Board fixtures
export const mockBoard: BoardResponse = {
  id: 'board-1',
  name: 'Test Board',
  description: 'A test soundboard',
  imageUrl: 'https://example.com/images/board.png',
  layout: 'grid',
  ownerId: 'user-1',
  keys: mockKeys,
  collaborators: mockCollaborators,
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockBoards: BoardResponse[] = [
  mockBoard,
  {
    id: 'board-2',
    name: 'Second Board',
    description: 'Another test soundboard',
    imageUrl: undefined,
    layout: 'list',
    ownerId: 'user-1',
    keys: [],
    collaborators: [],
    createdAt: new Date('2024-01-02'),
    updatedAt: new Date('2024-01-02'),
  },
]

// Media fixtures
export const mockAudioMedia: MediaResponse = {
  id: 'media-1',
  name: 'Test Audio',
  type: 'audio',
  url: 'https://example.com/audio/test.mp3',
  size: 1024 * 1024,
  mimeType: 'audio/mpeg',
  ownerId: 'user-1',
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockImageMedia: MediaResponse = {
  id: 'media-2',
  name: 'Test Image',
  type: 'image',
  url: 'https://example.com/images/test.png',
  size: 512 * 1024,
  mimeType: 'image/png',
  ownerId: 'user-1',
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockMedia: MediaResponse[] = [
  mockAudioMedia,
  mockImageMedia,
]

// API response fixtures
export const mockSigninResponse = {
  success: true,
  message: 'Login successful',
  data: {
    token: 'mock-jwt-token',
    user: mockUser,
  },
}

export const mockSignupResponse = {
  success: true,
  message: 'Registration successful',
}

export const mockBoardsResponse = {
  success: true,
  data: {
    boards: mockBoards,
  },
}

export const mockMediaResponse = {
  success: true,
  data: {
    media: mockMedia,
  },
}

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
  details: [
    { field: 'email', message: 'Invalid email format' },
  ],
}
