// MSW handlers for integration tests
// Reuses fixtures from __mocks__/data/fixtures.ts - DO NOT duplicate fixture data
import { http, HttpResponse } from 'msw'
import {
  mockUser,
  mockBoards,
  mockBoard,
  mockMedia,
  mockSigninResponse,
  mockSignupResponse,
  mockErrorResponse,
  mockKeys,
  mockKey,
} from '@/__mocks__/data/fixtures'

// Base URL pattern to match API requests
const API_BASE = '*/api/v1'

export const handlers = [
  // Auth endpoints
  http.post(`${API_BASE}/auth/signin`, async () => {
    return HttpResponse.json(mockSigninResponse)
  }),

  http.post(`${API_BASE}/auth/signup`, async () => {
    return HttpResponse.json(mockSignupResponse)
  }),

  http.post(`${API_BASE}/auth/forgot`, async () => {
    return HttpResponse.json({
      success: true,
      message: 'Password reset email sent',
    })
  }),

  http.post(`${API_BASE}/auth/reset`, async () => {
    return HttpResponse.json({
      success: true,
      message: 'Password reset successful',
    })
  }),

  // Users endpoints
  http.get(`${API_BASE}/users/me`, () => {
    return HttpResponse.json({
      success: true,
      data: mockUser,
    })
  }),

  http.put(`${API_BASE}/users/me`, async () => {
    return HttpResponse.json({
      success: true,
      data: mockUser,
    })
  }),

  // Boards endpoints
  http.get(`${API_BASE}/boards`, () => {
    return HttpResponse.json({
      success: true,
      data: { boards: mockBoards },
    })
  }),

  http.post(`${API_BASE}/boards`, async ({ request }) => {
    const body = await request.json() as { name: string; description?: string }
    return HttpResponse.json({
      success: true,
      data: {
        id: 'new-board-id',
        name: body.name,
        description: body.description,
        imageUrl: undefined,
        layout: 'grid',
        ownerId: mockUser.id,
        keys: [],
        collaborators: [],
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.get(`${API_BASE}/boards/:boardId`, ({ params }) => {
    const { boardId } = params
    if (boardId === mockBoard.id) {
      return HttpResponse.json(mockBoard)
    }
    return HttpResponse.json(mockErrorResponse, { status: 404 })
  }),

  http.put(`${API_BASE}/boards/:boardId`, async ({ params, request }) => {
    const { boardId } = params
    const body = await request.json() as { name?: string; description?: string }
    return HttpResponse.json({
      success: true,
      data: {
        ...mockBoard,
        id: boardId,
        ...body,
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.delete(`${API_BASE}/boards/:boardId`, () => {
    return HttpResponse.json({ success: true })
  }),

  // Keys endpoints
  http.get(`${API_BASE}/boards/:boardId/keys`, () => {
    return HttpResponse.json({
      success: true,
      data: { keys: mockKeys },
    })
  }),

  http.post(`${API_BASE}/boards/:boardId/keys`, async ({ request }) => {
    const body = await request.json() as { name: string; hotKey?: string; audioUrl?: string; imageUrl?: string }
    return HttpResponse.json({
      success: true,
      data: {
        id: 'new-key-id',
        ...body,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.get(`${API_BASE}/boards/:boardId/keys/:keyId`, ({ params }) => {
    const { keyId } = params
    const key = mockKeys.find(k => k.id === keyId)
    if (key) {
      return HttpResponse.json(key)
    }
    return HttpResponse.json(mockErrorResponse, { status: 404 })
  }),

  http.put(`${API_BASE}/boards/:boardId/keys/:keyId`, async ({ params, request }) => {
    const { keyId } = params
    const body = await request.json() as Partial<typeof mockKey>
    return HttpResponse.json({
      success: true,
      data: {
        ...mockKey,
        id: keyId,
        ...body,
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.delete(`${API_BASE}/boards/:boardId/keys/:keyId`, () => {
    return HttpResponse.json({ success: true })
  }),

  // Media endpoints
  http.get(`${API_BASE}/media/stats`, () => {
    return HttpResponse.json({
      success: true,
      data: {
        totalCount: mockMedia.length,
        imageCount: mockMedia.filter(m => m.mediaType === 'image').length,
        audioCount: mockMedia.filter(m => m.mediaType === 'audio').length,
        totalSize: mockMedia.reduce((sum, m) => sum + (m.fileSize || 0), 0),
      },
    })
  }),

  http.get(`${API_BASE}/media`, ({ request }) => {
    const url = new URL(request.url)
    const mediaType = url.searchParams.get('mediaType')
    
    let filteredMedia = mockMedia
    if (mediaType === 'image') {
      filteredMedia = mockMedia.filter(m => m.mediaType === 'image')
    } else if (mediaType === 'audio') {
      filteredMedia = mockMedia.filter(m => m.mediaType === 'audio')
    }
    
    return HttpResponse.json({
      success: true,
      data: { media: filteredMedia },
    })
  }),

  http.post(`${API_BASE}/media`, async ({ request }) => {
    const body = await request.json() as { fileName: string; mediaType: string; fileUrl: string }
    return HttpResponse.json({
      success: true,
      data: {
        id: 'new-media-id',
        fileName: body.fileName,
        mediaType: body.mediaType,
        fileUrl: body.fileUrl,
        fileSize: 1024,
        contentType: body.mediaType === 'image' ? 'image/png' : 'audio/mpeg',
        authorId: mockUser.id,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.get(`${API_BASE}/media/:mediaId`, ({ params }) => {
    const { mediaId } = params
    const media = mockMedia.find(m => m.id === mediaId)
    if (media) {
      return HttpResponse.json(media)
    }
    return HttpResponse.json(mockErrorResponse, { status: 404 })
  }),

  http.put(`${API_BASE}/media/:mediaId`, async ({ params, request }) => {
    const { mediaId } = params
    const body = await request.json() as { description?: string }
    const media = mockMedia.find(m => m.id === mediaId)
    return HttpResponse.json({
      success: true,
      data: {
        ...media,
        ...body,
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.delete(`${API_BASE}/media/:mediaId`, () => {
    return HttpResponse.json({ success: true })
  }),

  http.post(`${API_BASE}/media/:mediaId/process`, () => {
    return HttpResponse.json({
      success: true,
      message: 'Processing started',
    })
  }),

  // Collaborators endpoints
  http.get(`${API_BASE}/boards/:boardId/collaborators`, () => {
    return HttpResponse.json({
      success: true,
      data: { collaborators: mockBoard.collaborators },
    })
  }),

  http.post(`${API_BASE}/boards/:boardId/collaborators`, async ({ request }) => {
    const body = await request.json() as { email: string; role: string }
    return HttpResponse.json({
      success: true,
      data: {
        id: 'new-collab-id',
        ...body,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    })
  }),

  http.delete(`${API_BASE}/boards/:boardId/collaborators/:collaboratorId`, () => {
    return HttpResponse.json({ success: true })
  }),

  // Upload URL endpoint
  http.post(`${API_BASE}/media/upload-url`, async ({ request }) => {
    const body = await request.json() as { fileName: string; contentType: string }
    return HttpResponse.json({
      success: true,
      data: {
        signedUrl: 'https://storage.example.com/upload?signed=token',
        fileUrl: `https://storage.example.com/files/${body.fileName}`,
      },
    })
  }),
]

// Error handler overrides for testing error scenarios
export const errorHandlers = {
  signinError: http.post(`${API_BASE}/auth/signin`, () => {
    return HttpResponse.json(
      { success: false, message: 'Invalid credentials' },
      { status: 401 }
    )
  }),

  signupError: http.post(`${API_BASE}/auth/signup`, () => {
    return HttpResponse.json(
      { success: false, message: 'Email already exists' },
      { status: 409 }
    )
  }),

  boardsError: http.get(`${API_BASE}/boards`, () => {
    return HttpResponse.json(mockErrorResponse, { status: 500 })
  }),

  emptyBoards: http.get(`${API_BASE}/boards`, () => {
    return HttpResponse.json({
      success: true,
      data: { boards: [] },
    })
  }),

  emptyMedia: http.get(`${API_BASE}/media`, () => {
    return HttpResponse.json({
      success: true,
      data: { media: [] },
    })
  }),

  networkError: http.get(`${API_BASE}/boards`, () => {
    return HttpResponse.error()
  }),
}
