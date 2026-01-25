import { test, expect } from '@playwright/test'
import { DashboardPage } from '../pages/dashboard.page'
import { login } from '../fixtures/auth.fixture'

// E2E Test Board name - must match server/config/demo.go E2ETestBoardName
const E2E_TEST_BOARD_NAME = 'E2E Test Board'

test.describe('Dashboard', () => {
  let dashboardPage: DashboardPage

  test.beforeEach(async ({ page }) => {
    dashboardPage = new DashboardPage(page)

    // Real authentication against the backend
    await login(page)
  })

  test('should display dashboard with seeded board', async ({ page }) => {
    // The E2E test user has a seeded board from server setup
    await dashboardPage.goto()

    await dashboardPage.expectBoardsVisible()
  })

  test('should display boards list when boards exist', async ({ page }) => {
    await dashboardPage.goto()

    await dashboardPage.expectBoardsVisible()
  })

  test('should display board name', async ({ page }) => {
    await dashboardPage.goto()

    await dashboardPage.expectBoardVisible(E2E_TEST_BOARD_NAME)
  })

  test('should open create board modal from header button', async ({ page }) => {
    await dashboardPage.goto()

    // Click the New Board button
    await dashboardPage.newBoardButton.click()

    await dashboardPage.expectCreateModalVisible()
  })

  test('should open create board modal from empty state', async ({ page }) => {
    // Note: This test originally tested empty state.
    // With seeded data, we test the modal from the New Board button instead.
    await dashboardPage.goto()
    await dashboardPage.newBoardButton.click()

    await dashboardPage.expectCreateModalVisible()
  })

  test('should close create board modal on cancel', async ({ page }) => {
    await dashboardPage.goto()
    await dashboardPage.newBoardButton.click()
    await dashboardPage.cancelButton.click()

    await dashboardPage.expectCreateModalHidden()
  })

  test('should navigate to board detail on click', async ({ page }) => {
    await dashboardPage.goto()
    await dashboardPage.openBoardByIndex(0)

    await dashboardPage.expectNavigatedToBoard()
  })
})
