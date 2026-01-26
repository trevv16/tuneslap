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

  test('should display media gallery', async () => {
    await libraryPage.goto()

    await libraryPage.expectPageLoaded()
    // E2E test user has seeded media - if none exists, something is wrong
    await libraryPage.expectMediaVisible()
  })
})

test.describe('Media Details Sidebar', () => {
  let libraryPage: LibraryPage

  test.beforeEach(async ({ page }) => {
    libraryPage = new LibraryPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should open media details sidebar on item click', async () => {
    await libraryPage.goto()

    // E2E test user has seeded media - select first item
    await libraryPage.selectMediaByIndex(0)
    await libraryPage.expectDetailsSidebarVisible()
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
