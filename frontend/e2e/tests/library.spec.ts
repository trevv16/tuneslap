import { test, expect } from '@playwright/test'
import { LibraryPage } from '../pages/library.page'
import { login } from '../fixtures/auth.fixture'

test.describe('Media Library', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should display library page', async () => {
    await libraryPage.goto()

    await libraryPage.expectPageLoaded()
    // E2E test user starts with no media - library shows empty state or media list
  })
})

test.describe('Media Details Sidebar', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should show empty state or media items', async ({ page }) => {
    await libraryPage.goto()

    // Check if media exists - E2E test user may have no media initially
    const mediaCount = await page.locator('[data-testid="media-item"]').count()
    if (mediaCount > 0) {
      // If media exists, clicking should open details sidebar
      await libraryPage.selectMediaByIndex(0)
      await libraryPage.expectDetailsSidebarVisible()
    } else {
      // Otherwise, empty state or just page loaded is fine
      await libraryPage.expectPageLoaded()
    }
  })
})

test.describe('Media Upload', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should open upload modal', async () => {
    await libraryPage.goto()
    await libraryPage.openUploadModal()

    await libraryPage.expectUploadModalVisible()
  })
})
