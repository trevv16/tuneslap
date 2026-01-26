import { test, expect } from '@playwright/test'
import { BoardPage } from '../pages/board.page'
import { login } from '../fixtures/auth.fixture'

// E2E Test Board ID - must match server/config/demo.go E2ETestBoardID
const E2E_TEST_BOARD_ID = '000000000000000000000098'

test.describe('Board Detail', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should open add key sheet', async ({ page }) => {
    await boardPage.goto(E2E_TEST_BOARD_ID)
    await boardPage.openAddKeySheet()

    await boardPage.expectAddKeySheetVisible()
  })

  test('should have edit and add key buttons visible', async ({ page }) => {
    await boardPage.goto(E2E_TEST_BOARD_ID)

    await expect(boardPage.editButton).toBeVisible()
    await expect(boardPage.addKeyButton).toBeVisible()
  })
})

test.describe('Board Key Interaction', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should respond to keyboard hotkey press', async ({ page }) => {
    await boardPage.goto(E2E_TEST_BOARD_ID)

    // Press the hotkey for the first key (seeded test data has hotkey "1")
    await boardPage.pressHotkey('1')
  })
})

test.describe('Board Edit Page', () => {
  let boardPage: BoardPage

  test.beforeEach(async ({ page }) => {
    boardPage = new BoardPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should load edit page', async () => {
    await boardPage.gotoEdit(E2E_TEST_BOARD_ID)

    await boardPage.expectEditPageLoaded()
  })
})
