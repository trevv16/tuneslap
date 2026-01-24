import { renderHook, act } from '@testing-library/react'
import { useQueryParams, useLibraryParams } from './useQueryParams'

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
  usePathname: jest.fn(),
  useSearchParams: jest.fn(),
}))

import { useRouter, usePathname, useSearchParams } from 'next/navigation'

describe('useQueryParams', () => {
  const mockRouter = {
    push: jest.fn(),
    replace: jest.fn(),
  }

  beforeEach(() => {
    jest.clearAllMocks()
    ;(useRouter as jest.Mock).mockReturnValue(mockRouter)
    ;(usePathname as jest.Mock).mockReturnValue('/test')
    ;(useSearchParams as jest.Mock).mockReturnValue(new URLSearchParams())
  })

  describe('getParam', () => {
    it('should return param value when exists', () => {
      const searchParams = new URLSearchParams('foo=bar')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      expect(result.current.getParam('foo')).toBe('bar')
    })

    it('should return null when param does not exist', () => {
      const { result } = renderHook(() => useQueryParams())

      expect(result.current.getParam('nonexistent')).toBeNull()
    })
  })

  describe('getParamWithDefault', () => {
    it('should return param value when exists', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      expect(result.current.getParamWithDefault('tab', 'all')).toBe('audio')
    })

    it('should return default value when param does not exist', () => {
      const { result } = renderHook(() => useQueryParams())

      expect(result.current.getParamWithDefault('tab', 'all')).toBe('all')
    })
  })

  describe('setParam', () => {
    it('should add new param', () => {
      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParam('tab', 'audio')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test?tab=audio')
    })

    it('should update existing param', () => {
      const searchParams = new URLSearchParams('tab=all')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParam('tab', 'audio')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test?tab=audio')
    })

    it('should remove param when value is null', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParam('tab', null)
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test')
    })

    it('should remove param when value is empty string', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParam('tab', '')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test')
    })

    it('should preserve other params', () => {
      const searchParams = new URLSearchParams('view=grid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParam('tab', 'audio')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test?view=grid&tab=audio')
    })
  })

  describe('setParams', () => {
    it('should set multiple params at once', () => {
      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParams({ tab: 'audio', view: 'list' })
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test?tab=audio&view=list')
    })

    it('should remove params with null values', () => {
      const searchParams = new URLSearchParams('tab=audio&view=grid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.setParams({ tab: null, view: 'list' })
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test?view=list')
    })
  })

  describe('removeParam', () => {
    it('should remove specified param', () => {
      const searchParams = new URLSearchParams('tab=audio&view=grid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.removeParam('tab')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test?view=grid')
    })
  })

  describe('clearParams', () => {
    it('should remove all params', () => {
      const searchParams = new URLSearchParams('tab=audio&view=grid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      act(() => {
        result.current.clearParams()
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/test')
    })
  })

  describe('hasParam', () => {
    it('should return true when param exists', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      expect(result.current.hasParam('tab')).toBe(true)
    })

    it('should return false when param does not exist', () => {
      const { result } = renderHook(() => useQueryParams())

      expect(result.current.hasParam('tab')).toBe(false)
    })
  })

  describe('getAllParams', () => {
    it('should return all params as object', () => {
      const searchParams = new URLSearchParams('tab=audio&view=grid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useQueryParams())

      expect(result.current.getAllParams()).toEqual({
        tab: 'audio',
        view: 'grid',
      })
    })

    it('should return empty object when no params', () => {
      const { result } = renderHook(() => useQueryParams())

      expect(result.current.getAllParams()).toEqual({})
    })
  })
})

describe('useLibraryParams', () => {
  const mockRouter = {
    push: jest.fn(),
  }

  beforeEach(() => {
    jest.clearAllMocks()
    ;(useRouter as jest.Mock).mockReturnValue(mockRouter)
    ;(usePathname as jest.Mock).mockReturnValue('/library')
    ;(useSearchParams as jest.Mock).mockReturnValue(new URLSearchParams())
  })

  describe('tab validation', () => {
    it('should default to "all" when no tab param', () => {
      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.tab).toBe('all')
    })

    it('should return valid tab value', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.tab).toBe('audio')
    })

    it('should default to "all" for invalid tab', () => {
      const searchParams = new URLSearchParams('tab=invalid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.tab).toBe('all')
    })

    it('should accept "images" as valid tab', () => {
      const searchParams = new URLSearchParams('tab=images')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.tab).toBe('images')
    })
  })

  describe('view validation', () => {
    it('should default to "grid" when no view param', () => {
      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.view).toBe('grid')
    })

    it('should return valid view value', () => {
      const searchParams = new URLSearchParams('view=list')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.view).toBe('list')
    })

    it('should default to "grid" for invalid view', () => {
      const searchParams = new URLSearchParams('view=invalid')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.view).toBe('grid')
    })
  })

  describe('mediaType derivation', () => {
    it('should return undefined for "all" tab', () => {
      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.mediaType).toBeUndefined()
    })

    it('should return "audio" for audio tab', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.mediaType).toBe('audio')
    })

    it('should return "image" for images tab', () => {
      const searchParams = new URLSearchParams('tab=images')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      expect(result.current.mediaType).toBe('image')
    })
  })

  describe('setTab', () => {
    it('should set tab param', () => {
      const { result } = renderHook(() => useLibraryParams())

      act(() => {
        result.current.setTab('audio')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/library?tab=audio')
    })

    it('should remove tab param when setting to "all"', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      act(() => {
        result.current.setTab('all')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/library')
    })

    it('should preserve other params', () => {
      const searchParams = new URLSearchParams('view=list')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      act(() => {
        result.current.setTab('audio')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/library?view=list&tab=audio')
    })
  })

  describe('setView', () => {
    it('should set view param', () => {
      const { result } = renderHook(() => useLibraryParams())

      act(() => {
        result.current.setView('list')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/library?view=list')
    })

    it('should preserve other params', () => {
      const searchParams = new URLSearchParams('tab=audio')
      ;(useSearchParams as jest.Mock).mockReturnValue(searchParams)

      const { result } = renderHook(() => useLibraryParams())

      act(() => {
        result.current.setView('list')
      })

      expect(mockRouter.push).toHaveBeenCalledWith('/library?tab=audio&view=list')
    })
  })
})
