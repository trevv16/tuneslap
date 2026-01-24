import { test, expect } from '@playwright/test'
import { DashboardPage } from '../pages/dashboard.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockBoards, mockBoard, mockUserResponse } from '../fixtures/test-data'

test.describe('Dashboard', () => {
  let dashboardPage: DashboardPage

  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page)

    // Set up common API mocks BEFORE any navigation
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)

    // Navigate to set auth token
    await page.goto('/')
    await setAuthToken(page)
  })

  test('should display empty state when no boards exist', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.listEmpty()

    await dashboardPage.goto()

    await dashboardPage.expectEmptyState()
  })

  test('should display boards list when boards exist', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()

    await dashboardPage.expectBoardsVisible()
  })

  test('should display board name', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()

    await dashboardPage.expectBoardVisible(mockBoard.name)
  })

  test('should open create board modal from header button', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()

    // Click the New Board button
    await dashboardPage.newBoardButton.click()

    await dashboardPage.expectCreateModalVisible()
  })

  test('should open create board modal from empty state', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.listEmpty()

    await dashboardPage.goto()
    await dashboardPage.emptyStateButton.click()

    await dashboardPage.expectCreateModalVisible()
  })

  test('should close create board modal on cancel', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()
    await dashboardPage.newBoardButton.click()
    await dashboardPage.cancelButton.click()

    await dashboardPage.expectCreateModalHidden()
  })

  test('should navigate to board detail on click', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)
    await apiMocks.boards.get(mockBoard.id, mockBoard)

    await dashboardPage.goto()
    await dashboardPage.openBoardByIndex(0)

    await dashboardPage.expectNavigatedToBoard()
  })
})
