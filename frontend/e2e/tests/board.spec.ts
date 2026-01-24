import { test, expect } from '@playwright/test'
import { BoardPage } from '../pages/board.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockBoard, mockEmptyBoard, mockKeys, mockUserResponse } from '../fixtures/test-data'

test.describe('Board Detail', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
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
    await boardPage.expectKeyCount(mockKeys.length)
  })

  test('should display key names', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    for (const key of mockKeys) {
      await boardPage.expectKeyVisible(key.name)
    }
  })

  test('should open add key sheet', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)
    await apiMocks.media.list()

    await boardPage.goto(mockBoard.id)
    await boardPage.openAddKeySheet()

    await boardPage.expectAddKeySheetVisible()
  })

  test('should open add key sheet from empty state', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.getEmpty(mockEmptyBoard.id)
    await apiMocks.media.list()

    await boardPage.goto(mockEmptyBoard.id)
    await boardPage.openAddKeySheet()

    await boardPage.expectAddKeySheetVisible()
  })

  test('should close add key sheet on cancel', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)
    await apiMocks.media.list()

    await boardPage.goto(mockBoard.id)
    await boardPage.openAddKeySheet()
    await boardPage.closeAddKeySheet()

    await boardPage.expectAddKeySheetHidden()
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

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should be able to click on a key', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    // Click on the first key
    await boardPage.clickKeyByIndex(0)

    // The key should respond to click (visual feedback or audio play)
    // This test verifies the key is clickable without errors
  })

  test('should respond to keyboard hotkey press', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.goto(mockBoard.id)

    // Press the hotkey for the first key
    const firstKey = mockKeys[0]
    await boardPage.pressHotkey(firstKey.hotKey)

    // The key should respond to the hotkey
    // This test verifies the hotkey handling works
  })
})

test.describe('Board Edit Page', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
  })

  test('should load edit page', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await boardPage.gotoEdit(mockBoard.id)

    await expect(page).toHaveURL(/\/boards\/[a-z0-9-]+\/edit/)
  })
})
