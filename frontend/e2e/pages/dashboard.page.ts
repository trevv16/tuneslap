import type { Page, Locator } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Page object for the dashboard page.
 */
export class DashboardPage {
  readonly page: Page

  // Header elements
  readonly pageTitle: Locator
  readonly newBoardButton: Locator

  // Boards list
  readonly boardsList: Locator
  readonly boardCards: Locator

  // Empty state
  readonly emptyState: Locator
  readonly emptyStateTitle: Locator
  readonly emptyStateButton: Locator

  // Create board modal
  readonly createModal: Locator
  readonly boardNameInput: Locator
  readonly boardDescriptionInput: Locator
  readonly createBoardButton: Locator
  readonly cancelButton: Locator

  constructor(page: Page) {
    this.page = page

    // Header
    this.pageTitle = page.getByRole('heading', { level: 2 })
    this.newBoardButton = page.getByRole('button', { name: /create/i })

    // Boards list
    this.boardsList = page.getByTestId('boards-list')
    this.boardCards = page.locator('a[href^="/boards/"]')

    // Empty state
    this.emptyState = page.getByTestId('empty-state')
    this.emptyStateTitle = page.getByRole('heading', { name: /create your first/i })
    this.emptyStateButton = this.emptyState.getByRole('button', { name: /new board/i })

    // Create board modal
    this.createModal = page.getByRole('dialog')
    this.boardNameInput = page.getByLabel('Name')
    this.boardDescriptionInput = page.getByLabel('Description')
    this.createBoardButton = page.getByRole('button', { name: /create board/i })
    this.cancelButton = page.getByRole('button', { name: /cancel/i })
  }

  // Navigation
  async goto(): Promise<void> {
    await this.page.goto('/dashboard')
    // Wait for page to stabilize after data loading
    await this.page.waitForLoadState('networkidle')
    // Wait for either boards list or empty state to be visible
    await this.boardsList.or(this.emptyState).waitFor({ state: 'visible', timeout: 10000 })
  }

  // Actions
  async openCreateModal(): Promise<void> {
    // Try the header button first, then the empty state button
    const headerButton = this.newBoardButton
    const emptyButton = this.emptyStateButton

    if (await headerButton.isVisible()) {
      await headerButton.click()
    } else if (await emptyButton.isVisible()) {
      await emptyButton.click()
    }
  }

  async fillCreateBoardForm(name: string, description?: string): Promise<void> {
    await this.boardNameInput.fill(name)
    if (description) {
      await this.boardDescriptionInput.fill(description)
    }
  }

  async submitCreateBoard(): Promise<void> {
    await this.createBoardButton.click()
  }

  async createBoard(name: string, description?: string): Promise<void> {
    await this.openCreateModal()
    await this.fillCreateBoardForm(name, description)
    await this.submitCreateBoard()
  }

  async closeCreateModal(): Promise<void> {
    await this.cancelButton.click()
  }

  async openBoard(boardName: string): Promise<void> {
    await this.page.getByRole('link', { name: boardName }).click()
  }

  async openBoardByIndex(index: number): Promise<void> {
    await this.boardCards.nth(index).click()
  }

  // Assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.page).toHaveURL(/\/dashboard/)
    // Check for dashboard-specific elements
    await expect(this.pageTitle).toBeVisible()
    // Either boards list or empty state should be visible
    await expect(this.boardsList.or(this.emptyState)).toBeVisible()
  }

  async expectEmptyState(): Promise<void> {
    await expect(this.emptyState).toBeVisible()
  }

  async expectBoardsVisible(): Promise<void> {
    await expect(this.boardCards.first()).toBeVisible()
  }

  async expectBoardCount(count: number): Promise<void> {
    await expect(this.boardCards).toHaveCount(count)
  }

  async expectBoardVisible(name: string): Promise<void> {
    await expect(this.page.getByRole('link', { name })).toBeVisible()
  }

  async expectBoardWithDescription(name: string, description: string): Promise<void> {
    const boardCard = this.page.locator('a[href^="/boards/"]', { hasText: name })
    await expect(boardCard).toBeVisible()
    await expect(boardCard.getByText(description)).toBeVisible()
  }

  async expectCreateModalVisible(): Promise<void> {
    await expect(this.createModal).toBeVisible()
    await expect(this.boardNameInput).toBeVisible()
  }

  async expectCreateModalHidden(): Promise<void> {
    await expect(this.createModal).not.toBeVisible()
  }

  async expectNavigatedToBoard(): Promise<void> {
    await expect(this.page).toHaveURL(/\/boards\/[a-z0-9]+/)
  }
}
