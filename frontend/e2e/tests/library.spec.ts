import { test, expect } from '@playwright/test'
import { LibraryPage } from '../pages/library.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockMedia, mockAudioMedia, mockImageMedia, mockUserResponse } from '../fixtures/test-data'

test.describe('Media Library', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should display media gallery', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.media.list(mockMedia)

    await libraryPage.goto()

    await libraryPage.expectPageLoaded()
    await libraryPage.expectMediaVisible()
  })

  test('should display empty state when no media exists', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.media.listEmpty()

    await libraryPage.goto()

    await libraryPage.expectEmptyState()
  })

  test('should display correct media count', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.media.list(mockMedia)

    await libraryPage.goto()

    await libraryPage.expectMediaCount(mockMedia.length)
  })
})

test.describe('Media Details Sidebar', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.media.list(mockMedia)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should open media details sidebar on item click', async () => {
    await libraryPage.goto()
    await libraryPage.selectMediaByIndex(0)

    await libraryPage.expectDetailsSidebarVisible()
  })

  test('should close media details sidebar', async () => {
    await libraryPage.goto()
    await libraryPage.selectMediaByIndex(0)
    await libraryPage.closeDetailsSidebar()

    await libraryPage.expectDetailsSidebarHidden()
  })
})

test.describe('Media Upload', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.media.list(mockMedia)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should open upload modal', async () => {
    await libraryPage.goto()
    await libraryPage.openUploadModal()

    await libraryPage.expectUploadModalVisible()
  })
})
