import { test, expect } from '@playwright/test'
import { BoardPage } from '../pages/board.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockBoard, mockEmptyBoard, mockKeys, mockUserResponse } from '../fixtures/test-data'

test.describe('Board Detail', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should display board name in header', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    await boardPage.expectPageLoaded(mockBoard.name)
  })

  test('should display empty state when no keys exist', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.getEmpty(mockEmptyBoard.id)

    await boardPage.goto(mockEmptyBoard.id)

    await boardPage.expectEmptyState()
  })

  test('should display keys grid when keys exist', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    await boardPage.expectKeysVisible()
  })

  test('should open add key sheet', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)
    await apiMocks.media.list()

    await boardPage.goto(mockBoard.id)
    await boardPage.openAddKeySheet()

    await boardPage.expectAddKeySheetVisible()
  })

  test('should navigate to edit page', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)
    await boardPage.clickEditButton()

    await boardPage.expectNavigatedToEdit()
  })

  test('should have edit and add key buttons visible', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    await expect(boardPage.editButton).toBeVisible()
    await expect(boardPage.addKeyButton).toBeVisible()
  })
})

test.describe('Board Key Interaction', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should be able to click on a key', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    // Click on the first key - should not throw
    await boardPage.clickKeyByIndex(0)
  })

  test('should respond to keyboard hotkey press', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    // Press the hotkey for the first key
    const firstKey = mockKeys[0]
    await boardPage.pressHotkey(firstKey.hotKey)
  })
})

test.describe('Board Edit Page', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate and set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should load edit page', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.gotoEdit(mockBoard.id)

    await expect(page).toHaveURL(/\/boards\/[a-z0-9-]+\/edit/)
  })
})
