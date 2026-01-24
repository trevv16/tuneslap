// Reusable test data fixtures
import type {
  BoardResponse,
  KeyResponse,
  UserResponse,
  MediaResponse,
  CollaboratorResponse,
} from '@/api/models'
import {
  BoardResponseLayoutEnum,
  CollaboratorResponseRoleEnum,
  MediaResponseMediaTypeEnum,
  MediaResponseStatusEnum,
  MediaResponseContentTypeEnum,
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
  audioMediaId: 'media-1',
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
    audioMediaId: 'media-2',
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
    audioMediaId: 'media-3',
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
  email: 'user2@example.com',
  name: 'Second User',
  role: CollaboratorResponseRoleEnum.Editor,
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
  layout: BoardResponseLayoutEnum.Grid,
  authorId: 'user-1',
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
    layout: BoardResponseLayoutEnum.List,
    authorId: 'user-1',
    keys: [],
    collaborators: [],
    createdAt: new Date('2024-01-02'),
    updatedAt: new Date('2024-01-02'),
  },
]

// Media fixtures
export const mockAudioMedia: MediaResponse = {
  id: 'media-1',
  fileName: 'Test Audio',
  mediaType: MediaResponseMediaTypeEnum.Audio,
  fileUrl: 'https://example.com/audio/test.mp3',
  fileSize: 1024 * 1024,
  contentType: MediaResponseContentTypeEnum.AudioMp3,
  status: MediaResponseStatusEnum.Done,
  authorId: 'user-1',
  createdAt: new Date('2024-01-01'),
  updatedAt: new Date('2024-01-01'),
}

export const mockImageMedia: MediaResponse = {
  id: 'media-2',
  fileName: 'Test Image',
  mediaType: MediaResponseMediaTypeEnum.Image,
  fileUrl: 'https://example.com/images/test.png',
  fileSize: 512 * 1024,
  contentType: MediaResponseContentTypeEnum.ImagePng,
  status: MediaResponseStatusEnum.Done,
  authorId: 'user-1',
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
