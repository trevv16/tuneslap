import { test, expect } from '@playwright/test'
import { LibraryPage } from '../pages/library.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockMedia, mockAudioMedia, mockImageMedia, mockUserResponse } from '../fixtures/test-data'

test.describe('Media Library', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
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

  test('should filter by audio type', async ({ page }) => {
    const audioMedia = mockMedia.filter((m) => m.mediaType === 'audio')

    // Set up route to return filtered results
    await page.route('**/api/v1/media**', async (route) => {
      const url = route.request().url()
      if (url.includes('type=audio') || url.includes('mediaType=audio')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(audioMedia),
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(mockMedia),
        })
      }
    })

    await libraryPage.goto()
    await libraryPage.switchToAudioTab()

    // Should show only audio items
    await libraryPage.expectMediaItemVisible(mockAudioMedia.fileName)
  })

  test('should filter by images type', async ({ page }) => {
    const imageMedia = mockMedia.filter((m) => m.mediaType === 'image')

    // Set up route to return filtered results
    await page.route('**/api/v1/media**', async (route) => {
      const url = route.request().url()
      if (url.includes('type=image') || url.includes('mediaType=image')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(imageMedia),
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(mockMedia),
        })
      }
    })

    await libraryPage.goto()
    await libraryPage.switchToImagesTab()

    // Should show only image items
    await libraryPage.expectMediaItemVisible(mockImageMedia.fileName)
  })

  test('should show all media in all tab', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.media.list(mockMedia)

    await libraryPage.goto()
    await libraryPage.switchToAllTab()

    await libraryPage.expectMediaCount(mockMedia.length)
  })
})

test.describe('Media Details Sidebar', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.media.list(mockMedia)
  })

  test('should open media details sidebar on item click', async ({ page }) => {
    await libraryPage.goto()
    await libraryPage.selectMediaByIndex(0)

    await libraryPage.expectDetailsSidebarVisible()
  })

  test('should close media details sidebar', async ({ page }) => {
    await libraryPage.goto()
    await libraryPage.selectMediaByIndex(0)
    await libraryPage.closeDetailsSidebar()

    await libraryPage.expectDetailsSidebarHidden()
  })

  test('should toggle sidebar when clicking same item', async ({ page }) => {
    await libraryPage.goto()

    // First click opens sidebar
    await libraryPage.selectMediaByIndex(0)
    await libraryPage.expectDetailsSidebarVisible()

    // Second click on same item closes sidebar
    await libraryPage.selectMediaByIndex(0)
    await libraryPage.expectDetailsSidebarHidden()
  })

  test('should have download button in details sidebar', async ({ page }) => {
    await libraryPage.goto()
    await libraryPage.selectMediaByIndex(0)

    await expect(libraryPage.detailsDownloadButton).toBeVisible()
  })

  test('should have delete button in details sidebar', async ({ page }) => {
    await libraryPage.goto()
    await libraryPage.selectMediaByIndex(0)

    await expect(libraryPage.detailsDeleteButton).toBeVisible()
  })
})

test.describe('Media Upload', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.media.list(mockMedia)
  })

  test('should open upload modal', async ({ page }) => {
    await libraryPage.goto()
    await libraryPage.openUploadModal()

    await libraryPage.expectUploadModalVisible()
  })

  test('should have file input in upload modal', async ({ page }) => {
    await libraryPage.goto()
    await libraryPage.openUploadModal()

    await expect(libraryPage.fileInput).toBeAttached()
  })
})

test.describe('View Modes', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
    await apiMocks.media.list(mockMedia)
  })

  test('should have view toggle buttons', async ({ page }) => {
    await libraryPage.goto()

    // Should have grid and list view options
    await expect(
      libraryPage.gridViewButton.or(libraryPage.listViewButton).or(libraryPage.viewToggle)
    ).toBeVisible()
  })
})
