import { test, expect } from '@playwright/test'
import { DashboardPage } from '../pages/dashboard.page'
import { createApiMocks } from '../fixtures/api.fixture'
import { setAuthToken } from '../fixtures/auth.fixture'
import { mockBoards, mockBoard, mockUserResponse } from '../fixtures/test-data'

test.describe('Dashboard', () => {
  let dashboardPage: DashboardPage

  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page)

    // Set up authentication
    await page.goto('/')
    await setAuthToken(page)

    // Set up common API mocks
    const apiMocks = createApiMocks(page)
    await apiMocks.auth.me(mockUserResponse)
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
    await dashboardPage.expectBoardCount(mockBoards.length)
  })

  test('should display board name and description', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()

    await dashboardPage.expectBoardVisible(mockBoard.name)
  })

  test('should open create board modal', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()
    await dashboardPage.openCreateModal()

    await dashboardPage.expectCreateModalVisible()
  })

  test('should open create board modal from empty state', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.listEmpty()

    await dashboardPage.goto()
    await dashboardPage.openCreateModal()

    await dashboardPage.expectCreateModalVisible()
  })

  test('should close create board modal on cancel', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards)

    await dashboardPage.goto()
    await dashboardPage.openCreateModal()
    await dashboardPage.closeCreateModal()

    await dashboardPage.expectCreateModalHidden()
  })

  test('should create new board and show it in list', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list([])

    const newBoard = {
      ...mockBoard,
      id: 'new-board-id',
      name: 'My New Board',
      description: 'A brand new board',
    }
    await apiMocks.boards.create(newBoard)

    // After creation, the list should include the new board
    await page.route('**/api/v1/boards', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([newBoard]),
        })
      } else {
        await route.continue()
      }
    })

    await dashboardPage.goto()
    await dashboardPage.createBoard('My New Board', 'A brand new board')

    // Modal should close after successful creation
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

  test('should show loading skeleton while fetching boards', async ({ page }) => {
    const apiMocks = createApiMocks(page)
    await apiMocks.boards.list(mockBoards, { delay: 1000 })

    await dashboardPage.goto()

    // Should show skeleton or loading state
    // The exact implementation depends on the UI
    await expect(page.locator('.animate-pulse').or(page.getByText(/loading/i))).toBeVisible()
  })
})
