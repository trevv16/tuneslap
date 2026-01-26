import type { Page, Locator } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Page object for the media library page.
 */
export class LibraryPage {
  readonly page: Page

  // Header elements
  readonly pageTitle: Locator
  readonly addFileButton: Locator

  // Tabs
  readonly tabsList: Locator
  readonly allTab: Locator
  readonly audioTab: Locator
  readonly imagesTab: Locator

  // View toggle
  readonly viewToggle: Locator
  readonly gridViewButton: Locator
  readonly listViewButton: Locator

  // Media gallery
  readonly mediaGallery: Locator
  readonly mediaItems: Locator
  readonly mediaGridItems: Locator
  readonly mediaListItems: Locator

  // Empty state
  readonly emptyState: Locator
  readonly emptyUploadButton: Locator

  // Media details sidebar
  readonly detailsSidebar: Locator
  readonly detailsCloseButton: Locator
  readonly detailsFileName: Locator
  readonly detailsDownloadButton: Locator
  readonly detailsDeleteButton: Locator

  // Upload modal
  readonly uploadModal: Locator
  readonly fileInput: Locator
  readonly uploadButton: Locator

  constructor(page: Page) {
    this.page = page

    // Header
    this.pageTitle = page.getByRole('heading', { level: 1 })
    this.addFileButton = page.getByRole('button', { name: /add|upload/i })

    // Tabs
    this.tabsList = page.getByRole('tablist')
    this.allTab = page.getByRole('tab', { name: /all/i })
    this.audioTab = page.getByRole('tab', { name: /audio/i })
    this.imagesTab = page.getByRole('tab', { name: /images/i })

    // View toggle
    this.viewToggle = page.locator('[data-testid="view-toggle"]').or(
      page.locator('button').filter({ has: page.locator('svg') }).first()
    )
    this.gridViewButton = page.getByRole('button', { name: /grid/i })
    this.listViewButton = page.getByRole('button', { name: /list/i })

    // Media gallery
    this.mediaGallery = page.locator('main')
    this.mediaItems = page.locator('[data-testid="media-item"]').or(
      page.locator('button').filter({ has: page.locator('img') })
    )
    this.mediaGridItems = page.locator('[data-testid="media-grid-item"]')
    this.mediaListItems = page.locator('[data-testid="media-list-item"]')

    // Empty state
    this.emptyState = page.getByText(/no media|upload your first/i)
    this.emptyUploadButton = page.getByRole('button', { name: /upload/i })

    // Media details sidebar
    this.detailsSidebar = page.locator('aside').or(page.locator('[data-testid="media-details"]'))
    this.detailsCloseButton = page.getByRole('button', { name: /close/i })
    this.detailsFileName = page.locator('[data-testid="media-filename"]').or(
      page.locator('aside h2, aside h3')
    )
    this.detailsDownloadButton = page.getByRole('button', { name: /download/i })
    this.detailsDeleteButton = page.getByRole('button', { name: /delete/i })

    // Upload modal
    this.uploadModal = page.getByRole('dialog')
    this.fileInput = page.locator('input[type="file"]')
    this.uploadButton = page.getByRole('button', { name: /upload|submit/i })
  }

  // Navigation
  async goto(): Promise<void> {
    await this.page.goto('/library')
  }

  // Tab actions
  async switchToAllTab(): Promise<void> {
    await this.allTab.click()
  }

  async switchToAudioTab(): Promise<void> {
    await this.audioTab.click()
  }

  async switchToImagesTab(): Promise<void> {
    await this.imagesTab.click()
  }

  async switchTab(tab: 'all' | 'audio' | 'images'): Promise<void> {
    switch (tab) {
      case 'all':
        await this.switchToAllTab()
        break
      case 'audio':
        await this.switchToAudioTab()
        break
      case 'images':
        await this.switchToImagesTab()
        break
    }
  }

  // View actions
  async switchToGridView(): Promise<void> {
    await this.gridViewButton.click()
  }

  async switchToListView(): Promise<void> {
    await this.listViewButton.click()
  }

  async switchView(view: 'grid' | 'list'): Promise<void> {
    if (view === 'grid') {
      await this.switchToGridView()
    } else {
      await this.switchToListView()
    }
  }

  // Media actions
  async selectMedia(fileName: string): Promise<void> {
    await this.page.getByText(fileName).click()
  }

  async selectMediaByIndex(index: number): Promise<void> {
    await this.mediaItems.nth(index).click()
  }

  async closeDetailsSidebar(): Promise<void> {
    await this.detailsCloseButton.click()
  }

  async downloadSelectedMedia(): Promise<void> {
    await this.detailsDownloadButton.click()
  }

  async deleteSelectedMedia(): Promise<void> {
    await this.detailsDeleteButton.click()
  }

  // Upload actions
  async openUploadModal(): Promise<void> {
    await this.addFileButton.click()
  }

  async uploadFile(filePath: string): Promise<void> {
    await this.openUploadModal()
    await this.fileInput.setInputFiles(filePath)
    await this.uploadButton.click()
  }

  // Assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.page).toHaveURL(/\/library/)
    // Check for library-specific elements
    await expect(this.pageTitle).toBeVisible()
    await expect(this.mediaGallery).toBeVisible()
  }

  async expectEmptyState(): Promise<void> {
    await expect(this.emptyState).toBeVisible()
  }

  async expectMediaVisible(): Promise<void> {
    await expect(this.mediaItems.first()).toBeVisible()
  }

  async expectMediaCount(count: number): Promise<void> {
    await expect(this.mediaItems).toHaveCount(count)
  }

  async expectMediaItemVisible(fileName: string): Promise<void> {
    await expect(this.page.getByText(fileName)).toBeVisible()
  }

  async expectDetailsSidebarVisible(): Promise<void> {
    await expect(this.detailsSidebar).toBeVisible()
  }

  async expectDetailsSidebarHidden(): Promise<void> {
    await expect(this.detailsSidebar).not.toBeVisible()
  }

  async expectTabActive(tab: 'all' | 'audio' | 'images'): Promise<void> {
    let tabElement: Locator
    switch (tab) {
      case 'all':
        tabElement = this.allTab
        break
      case 'audio':
        tabElement = this.audioTab
        break
      case 'images':
        tabElement = this.imagesTab
        break
    }
    await expect(tabElement).toHaveAttribute('data-state', 'active')
  }

  async expectUploadModalVisible(): Promise<void> {
    await expect(this.uploadModal).toBeVisible()
  }

  async expectUploadModalHidden(): Promise<void> {
    await expect(this.uploadModal).not.toBeVisible()
  }
}
