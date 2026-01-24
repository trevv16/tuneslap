import { render, screen, fireEvent } from '@testing-library/react'
import BoardsList from './BoardsList'
import { mockBoards, mockBoard } from '@/__mocks__/data/fixtures'

// Mock next/link
jest.mock('next/link', () => {
  const React = require('react')
  return function MockLink({ href, children, onMouseEnter, className }: {
    href: string
    children: React.ReactNode
    onMouseEnter?: () => void
    className?: string
  }) {
    return React.createElement('a', {
      href,
      onMouseEnter,
      className,
      'data-testid': 'board-link',
    }, children)
  }
})

// Mock the usePreloadBoard hook
jest.mock('@/hooks/usePreloadBoard', () => ({
  usePreloadBoard: jest.fn(() => ({
    preloadBoard: jest.fn(),
    isPreloading: jest.fn(() => false),
  })),
}))

import { usePreloadBoard } from '@/hooks/usePreloadBoard'

describe('BoardsList', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    ;(usePreloadBoard as jest.Mock).mockReturnValue({
      preloadBoard: jest.fn(),
      isPreloading: jest.fn(() => false),
    })
  })

  it('should render list of boards', () => {
    render(<BoardsList boards={mockBoards} />)

    mockBoards.forEach(board => {
      expect(screen.getByText(board.name)).toBeInTheDocument()
    })
  })

  it('should render board descriptions', () => {
    render(<BoardsList boards={mockBoards} />)

    mockBoards.forEach(board => {
      if (board.description) {
        expect(screen.getByText(board.description)).toBeInTheDocument()
      }
    })
  })

  it('should show "No description" for boards without description', () => {
    const boardsWithNoDescription = [{ ...mockBoard, description: undefined }]
    render(<BoardsList boards={boardsWithNoDescription} />)

    expect(screen.getByText('No description')).toBeInTheDocument()
  })

  it('should render correct number of boards', () => {
    render(<BoardsList boards={mockBoards} />)

    const boardLinks = screen.getAllByTestId('board-link')
    expect(boardLinks).toHaveLength(mockBoards.length)
  })

  it('should link to board detail page', () => {
    render(<BoardsList boards={[mockBoard]} />)

    const link = screen.getByTestId('board-link')
    expect(link).toHaveAttribute('href', `/boards/${mockBoard.id}`)
  })

  it('should display key count', () => {
    render(<BoardsList boards={[mockBoard]} />)

    // Board has keys
    const keyCount = mockBoard.keys?.length ?? 0
    expect(screen.getByText(keyCount.toString())).toBeInTheDocument()
  })

  it('should display collaborator count', () => {
    render(<BoardsList boards={[mockBoard]} />)

    const collabCount = mockBoard.collaborators?.length ?? 0
    expect(screen.getByText(collabCount.toString())).toBeInTheDocument()
  })

  it('should render board image when available', () => {
    render(<BoardsList boards={[mockBoard]} />)

    const image = screen.getByAltText(mockBoard.name)
    expect(image).toBeInTheDocument()
    expect(image).toHaveAttribute('src', mockBoard.imageUrl)
  })

  it('should render fallback icon when no image', () => {
    const boardWithNoImage = { ...mockBoard, imageUrl: undefined }
    render(<BoardsList boards={[boardWithNoImage]} />)

    // Should show the LayoutGrid icon (SVG)
    expect(screen.queryByAltText(boardWithNoImage.name)).not.toBeInTheDocument()
  })

  it('should show layout badge', () => {
    render(<BoardsList boards={[mockBoard]} />)

    expect(screen.getByText(mockBoard.layout?.toUpperCase() ?? '')).toBeInTheDocument()
  })

  it('should call preloadBoard on mouse enter', () => {
    const mockPreloadBoard = jest.fn()
    ;(usePreloadBoard as jest.Mock).mockReturnValue({
      preloadBoard: mockPreloadBoard,
      isPreloading: jest.fn(() => false),
    })

    render(<BoardsList boards={[mockBoard]} />)

    const link = screen.getByTestId('board-link')
    fireEvent.mouseEnter(link)

    // The preloadBoard should be called when board.id exists
    expect(mockPreloadBoard).toHaveBeenCalled()
  })

  it('should show preloading indicator when preloading', () => {
    ;(usePreloadBoard as jest.Mock).mockReturnValue({
      preloadBoard: jest.fn(),
      isPreloading: jest.fn((id: string) => id === mockBoard.id),
    })

    render(<BoardsList boards={[mockBoard]} />)

    // Should have a pulsing indicator
    const indicator = document.querySelector('.animate-pulse')
    expect(indicator).toBeInTheDocument()
  })

  it('should format date correctly', () => {
    const boardWithDate = {
      ...mockBoard,
      createdAt: new Date('2024-01-15T00:00:00.000Z'),
    }
    render(<BoardsList boards={[boardWithDate]} />)

    // Check for formatted date - the format is "Mon DD, YYYY"
    // Using a flexible matcher since timezone can affect the exact date
    const dateElement = screen.getByText(/Jan \d{1,2}, 2024/)
    expect(dateElement).toBeInTheDocument()
  })

  it('should show "Unknown" for missing date', () => {
    const boardWithNoDate = { ...mockBoard, createdAt: undefined }
    render(<BoardsList boards={[boardWithNoDate]} />)

    expect(screen.getByText('Unknown')).toBeInTheDocument()
  })
})
