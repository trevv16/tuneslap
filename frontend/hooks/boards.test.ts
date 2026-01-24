import { renderHook, waitFor } from '@testing-library/react'
import {
  useGetBoards,
  useGetBoardById,
  useCreateBoard,
  useUpdateBoard,
  useDeleteBoard,
  boardKeys,
} from './boards'
import { createWrapper } from '@/__mocks__/providers/allProviders'
import { mockBoards, mockBoard } from '@/__mocks__/data/fixtures'
import { CreateBoardRequestLayoutEnum } from '@/api/models'

// Mock the API and config modules
jest.mock('@/api', () => ({
  BoardsApi: jest.fn(),
}))

jest.mock('@/api/config', () => ({
  getApiConfig: jest.fn(() => ({})),
}))

jest.mock('@/utils/token', () => ({
  getStoredToken: jest.fn(() => 'mock-token'),
}))

import { BoardsApi } from '@/api'

describe('boardKeys', () => {
  it('should return correct key for all boards', () => {
    expect(boardKeys.all()).toEqual(['boards'])
  })

  it('should return correct key for board detail', () => {
    expect(boardKeys.detail('board-123')).toEqual(['board', 'board-123'])
  })
})

describe('useGetBoards', () => {
  const mockGetAllBoards = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockBoardsApi = BoardsApi as jest.MockedClass<typeof BoardsApi>
    MockBoardsApi.mockImplementation(() => ({
      getAllBoards: mockGetAllBoards,
    }) as unknown as InstanceType<typeof BoardsApi>)
  })

  it('should fetch all boards successfully', async () => {
    mockGetAllBoards.mockResolvedValue({
      data: { boards: mockBoards },
    })

    const { result } = renderHook(() => useGetBoards(), {
      wrapper: createWrapper(),
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(result.current.data).toEqual(mockBoards)
    expect(mockGetAllBoards).toHaveBeenCalledWith({})
  })

  it('should return empty array when no boards data', async () => {
    mockGetAllBoards.mockResolvedValue({
      data: {},
    })

    const { result } = renderHook(() => useGetBoards(), {
      wrapper: createWrapper(),
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(result.current.data).toEqual([])
  })

  it('should not fetch when token is empty', () => {
    jest.requireMock('@/utils/token').getStoredToken.mockReturnValue('')

    const { result } = renderHook(() => useGetBoards(), {
      wrapper: createWrapper(),
    })

    // Query should not be enabled
    expect(result.current.fetchStatus).toBe('idle')
    expect(mockGetAllBoards).not.toHaveBeenCalled()
  })
})

describe('useGetBoardById', () => {
  const mockGetBoardById = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    jest.requireMock('@/utils/token').getStoredToken.mockReturnValue('mock-token')
    const MockBoardsApi = BoardsApi as jest.MockedClass<typeof BoardsApi>
    MockBoardsApi.mockImplementation(() => ({
      getBoardById: mockGetBoardById,
    }) as unknown as InstanceType<typeof BoardsApi>)
  })

  it('should fetch board by ID successfully', async () => {
    mockGetBoardById.mockResolvedValue(mockBoard)

    const { result } = renderHook(() => useGetBoardById('board-1'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(result.current.data).toEqual(mockBoard)
    expect(mockGetBoardById).toHaveBeenCalledWith({ boardId: 'board-1' })
  })

  it('should handle board not found error', async () => {
    const error = new Error('Board not found')
    mockGetBoardById.mockRejectedValue(error)

    const { result } = renderHook(() => useGetBoardById('nonexistent'), {
      wrapper: createWrapper(),
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('Board not found')
  })
})

describe('useCreateBoard', () => {
  const mockCreateBoard = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockBoardsApi = BoardsApi as jest.MockedClass<typeof BoardsApi>
    MockBoardsApi.mockImplementation(() => ({
      createBoard: mockCreateBoard,
    }) as unknown as InstanceType<typeof BoardsApi>)
  })

  it('should create board successfully', async () => {
    mockCreateBoard.mockResolvedValue({
      success: true,
      data: mockBoard,
    })

    const { result } = renderHook(() => useCreateBoard(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      name: 'New Board',
      description: 'A new test board',
      layout: CreateBoardRequestLayoutEnum.Grid,
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockCreateBoard).toHaveBeenCalledWith({
      createBoardRequest: {
        name: 'New Board',
        description: 'A new test board',
        layout: CreateBoardRequestLayoutEnum.Grid,
      },
    })
  })

  it('should handle create board error', async () => {
    const error = new Error('Failed to create board')
    mockCreateBoard.mockRejectedValue(error)

    const { result } = renderHook(() => useCreateBoard(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      name: 'New Board',
      layout: CreateBoardRequestLayoutEnum.Grid,
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('Failed to create board')
  })
})

describe('useUpdateBoard', () => {
  const mockUpdateBoard = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockBoardsApi = BoardsApi as jest.MockedClass<typeof BoardsApi>
    MockBoardsApi.mockImplementation(() => ({
      updateBoard: mockUpdateBoard,
    }) as unknown as InstanceType<typeof BoardsApi>)
  })

  it('should update board successfully', async () => {
    mockUpdateBoard.mockResolvedValue({
      success: true,
      data: { ...mockBoard, name: 'Updated Board' },
    })

    const { result } = renderHook(() => useUpdateBoard(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      boardId: 'board-1',
      name: 'Updated Board',
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockUpdateBoard).toHaveBeenCalledWith({
      boardId: 'board-1',
      updateBoardRequest: {
        name: 'Updated Board',
      },
    })
  })

  it('should handle update board error', async () => {
    const error = new Error('Unauthorized')
    mockUpdateBoard.mockRejectedValue(error)

    const { result } = renderHook(() => useUpdateBoard(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      boardId: 'board-1',
      name: 'Updated Board',
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('Unauthorized')
  })
})

describe('useDeleteBoard', () => {
  const mockDeleteBoard = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockBoardsApi = BoardsApi as jest.MockedClass<typeof BoardsApi>
    MockBoardsApi.mockImplementation(() => ({
      deleteBoard: mockDeleteBoard,
    }) as unknown as InstanceType<typeof BoardsApi>)
  })

  it('should delete board successfully', async () => {
    mockDeleteBoard.mockResolvedValue(undefined)

    const { result } = renderHook(() => useDeleteBoard(), {
      wrapper: createWrapper(),
    })

    result.current.mutate('board-1')

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockDeleteBoard).toHaveBeenCalledWith({ boardId: 'board-1' })
  })

  it('should handle delete board error', async () => {
    const error = new Error('Board not found')
    mockDeleteBoard.mockRejectedValue(error)

    const { result } = renderHook(() => useDeleteBoard(), {
      wrapper: createWrapper(),
    })

    result.current.mutate('nonexistent')

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('Board not found')
  })
})
