'use client'

import { useRouter, useSearchParams, usePathname } from 'next/navigation'
import { useCallback } from 'react'

type QueryParamValue = string | null | undefined

/**
 * Hook for managing URL query parameters
 * Provides helpers for getting, setting, and removing query params
 */
export function useQueryParams() {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()

  /**
   * Get a query param value
   */
  const getParam = useCallback((key: string): string | null => {
    return searchParams.get(key)
  }, [searchParams])

  /**
   * Get a query param value with a default fallback
   */
  const getParamWithDefault = useCallback(<T extends string>(key: string, defaultValue: T): T => {
    const value = searchParams.get(key)
    return (value as T) || defaultValue
  }, [searchParams])

  /**
   * Set a single query param (preserves other params)
   */
  const setParam = useCallback((key: string, value: QueryParamValue) => {
    const params = new URLSearchParams(searchParams.toString())
    
    if (value === null || value === undefined || value === '') {
      params.delete(key)
    } else {
      params.set(key, value)
    }
    
    const queryString = params.toString()
    router.push(queryString ? `${pathname}?${queryString}` : pathname)
  }, [router, pathname, searchParams])

  /**
   * Set multiple query params at once (preserves other params)
   */
  const setParams = useCallback((newParams: Record<string, QueryParamValue>) => {
    const params = new URLSearchParams(searchParams.toString())
    
    Object.entries(newParams).forEach(([key, value]) => {
      if (value === null || value === undefined || value === '') {
        params.delete(key)
      } else {
        params.set(key, value)
      }
    })
    
    const queryString = params.toString()
    router.push(queryString ? `${pathname}?${queryString}` : pathname)
  }, [router, pathname, searchParams])

  /**
   * Remove a query param
   */
  const removeParam = useCallback((key: string) => {
    const params = new URLSearchParams(searchParams.toString())
    params.delete(key)
    
    const queryString = params.toString()
    router.push(queryString ? `${pathname}?${queryString}` : pathname)
  }, [router, pathname, searchParams])

  /**
   * Clear all query params
   */
  const clearParams = useCallback(() => {
    router.push(pathname)
  }, [router, pathname])

  /**
   * Check if a param exists
   */
  const hasParam = useCallback((key: string): boolean => {
    return searchParams.has(key)
  }, [searchParams])

  /**
   * Get all params as an object
   */
  const getAllParams = useCallback((): Record<string, string> => {
    const params: Record<string, string> = {}
    searchParams.forEach((value, key) => {
      params[key] = value
    })
    return params
  }, [searchParams])

  return {
    getParam,
    getParamWithDefault,
    setParam,
    setParams,
    removeParam,
    clearParams,
    hasParam,
    getAllParams,
    searchParams,
  }
}

// Library-specific query param types and defaults
export type LibraryTab = 'all' | 'audio' | 'images'
export type LibraryView = 'grid' | 'list'

const VALID_TABS: LibraryTab[] = ['all', 'audio', 'images']
const VALID_VIEWS: LibraryView[] = ['grid', 'list']

/**
 * Library-specific hook for managing tab and view query params
 * Reads directly from URL and validates values
 */
export function useLibraryParams() {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()

  // Read and validate tab from URL (default to 'all' if invalid or missing)
  const rawTab = searchParams.get('tab')
  const tab: LibraryTab = VALID_TABS.includes(rawTab as LibraryTab) ? (rawTab as LibraryTab) : 'all'

  // Read and validate view from URL (default to 'grid' if invalid or missing)
  const rawView = searchParams.get('view')
  const view: LibraryView = VALID_VIEWS.includes(rawView as LibraryView) ? (rawView as LibraryView) : 'grid'

  // Derive mediaType directly from validated tab
  const mediaType: 'audio' | 'image' | undefined =
    tab === 'audio' ? 'audio' :
    tab === 'images' ? 'image' :
    undefined

  const setTab = useCallback((newTab: LibraryTab) => {
    const params = new URLSearchParams(searchParams.toString())
    if (newTab === 'all') {
      params.delete('tab')
    } else {
      params.set('tab', newTab)
    }
    const query = params.toString()
    router.push(query ? `${pathname}?${query}` : pathname)
  }, [router, pathname, searchParams])

  const setView = useCallback((newView: LibraryView) => {
    const params = new URLSearchParams(searchParams.toString())
    params.set('view', newView)
    const query = params.toString()
    router.push(query ? `${pathname}?${query}` : pathname)
  }, [router, pathname, searchParams])

  return {
    tab,
    view,
    mediaType,
    setTab,
    setView,
  }
}
