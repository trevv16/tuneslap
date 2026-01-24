import type { Page, Locator } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Page object for the board detail page.
 */
export class BoardPage {
  readonly page: Page

  // Header elements
  readonly pageTitle: Locator
  readonly editButton: Locator
  readonly addKeyButton: Locator

  // Keys grid
  readonly keysGrid: Locator
  readonly keyCards: Locator

  // Empty state
  readonly emptyState: Locator
  readonly emptyStateAddButton: Locator

  // Add key sheet
  readonly addKeySheet: Locator
  readonly keyNameInput: Locator
  readonly keyDescriptionInput: Locator
  readonly keyHotkeyInput: Locator
  readonly audioMediaSelect: Locator
  readonly imageMediaSelect: Locator
  readonly saveKeyButton: Locator
  readonly cancelKeyButton: Locator

  constructor(page: Page) {
    this.page = page

    // Header
    this.pageTitle = page.getByRole('heading', { level: 1 })
    this.editButton = page.getByRole('link', { name: /edit/i })
    this.addKeyButton = page.getByRole('button', { name: /add key/i })

    // Keys grid
    this.keysGrid = page.locator('[data-testid="keys-grid"]').or(page.locator('main'))
    this.keyCards = page.locator('[data-testid="sound-key"]').or(
      page.locator('button').filter({ has: page.locator('img') })
    )

    // Empty state
    this.emptyState = page.getByText(/no keys|add your first key/i)
    this.emptyStateAddButton = page.getByRole('button', { name: /add key/i })

    // Add key sheet
    this.addKeySheet = page.locator('[role="dialog"]')
    this.keyNameInput = page.getByLabel('Name')
    this.keyDescriptionInput = page.getByLabel('Description')
    this.keyHotkeyInput = page.getByLabel(/hotkey|key/i)
    this.audioMediaSelect = page.getByLabel(/audio/i)
    this.imageMediaSelect = page.getByLabel(/image/i)
    this.saveKeyButton = page.getByRole('button', { name: /save|add|create/i })
    this.cancelKeyButton = page.getByRole('button', { name: /cancel/i })
  }

  // Navigation
  async goto(boardId: string): Promise<void> {
    await this.page.goto(`/boards/${boardId}`)
  }

  async gotoEdit(boardId: string): Promise<void> {
    await this.page.goto(`/boards/${boardId}/edit`)
  }

  // Actions
  async openAddKeySheet(): Promise<void> {
    const headerButton = this.addKeyButton
    const emptyButton = this.emptyStateAddButton

    if (await headerButton.isVisible()) {
      await headerButton.click()
    } else if (await emptyButton.isVisible()) {
      await emptyButton.click()
    }
  }

  async fillKeyForm(data: {
    name: string
    description?: string
    hotkey: string
  }): Promise<void> {
    await this.keyNameInput.fill(data.name)
    if (data.description) {
      await this.keyDescriptionInput.fill(data.description)
    }
    await this.keyHotkeyInput.fill(data.hotkey)
  }

  async submitAddKey(): Promise<void> {
    await this.saveKeyButton.click()
  }

  async closeAddKeySheet(): Promise<void> {
    await this.cancelKeyButton.click()
  }

  async clickKey(keyName: string): Promise<void> {
    await this.page.locator('button', { hasText: keyName }).click()
  }

  async clickKeyByIndex(index: number): Promise<void> {
    await this.keyCards.nth(index).click()
  }

  async clickEditButton(): Promise<void> {
    await this.editButton.click()
  }

  async pressHotkey(key: string): Promise<void> {
    await this.page.keyboard.press(key)
  }

  // Assertions
  async expectPageLoaded(boardName?: string): Promise<void> {
    await expect(this.page).toHaveURL(/\/boards\/[a-z0-9]+/)
    if (boardName) {
      await expect(this.pageTitle).toContainText(boardName)
    }
  }

  async expectEmptyState(): Promise<void> {
    await expect(this.emptyState).toBeVisible()
  }

  async expectKeysVisible(): Promise<void> {
    await expect(this.keyCards.first()).toBeVisible()
  }

  async expectKeyCount(count: number): Promise<void> {
    await expect(this.keyCards).toHaveCount(count)
  }

  async expectKeyVisible(keyName: string): Promise<void> {
    await expect(this.page.getByText(keyName)).toBeVisible()
  }

  async expectKeyWithHotkey(keyName: string, hotkey: string): Promise<void> {
    const keyCard = this.page.locator('button', { hasText: keyName })
    await expect(keyCard).toBeVisible()
    await expect(keyCard.getByText(hotkey)).toBeVisible()
  }

  async expectAddKeySheetVisible(): Promise<void> {
    await expect(this.addKeySheet).toBeVisible()
    await expect(this.keyNameInput).toBeVisible()
  }

  async expectAddKeySheetHidden(): Promise<void> {
    await expect(this.addKeySheet).not.toBeVisible()
  }

  async expectNavigatedToEdit(): Promise<void> {
    await expect(this.page).toHaveURL(/\/boards\/[a-z0-9]+\/edit/)
  }
}
