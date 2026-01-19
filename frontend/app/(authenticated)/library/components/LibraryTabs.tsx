'use client'

import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import type { LibraryTab, LibraryView } from '@/hooks/useQueryParams'
import { cn } from '@/lib/utils'
import { LayoutGrid, List } from 'lucide-react'

const tabs: { name: string; value: LibraryTab }[] = [
  { name: 'All Media', value: 'all' },
  { name: 'Audio', value: 'audio' },
  { name: 'Images', value: 'images' },
]

type LibraryTabsProps = {
  currentTab: LibraryTab
  onTabChange: (tab: LibraryTab) => void
  viewMode: LibraryView
  onViewModeChange: (mode: LibraryView) => void
}

export default function LibraryTabs({ currentTab, onTabChange, viewMode, onViewModeChange }: LibraryTabsProps) {
  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b pb-4">
      {/* Mobile: Select dropdown */}
      <div className="sm:hidden">
        <Select value={currentTab} onValueChange={(value) => onTabChange(value as LibraryTab)}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select tab" />
          </SelectTrigger>
          <SelectContent>
            {tabs.map((tab) => (
              <SelectItem key={tab.value} value={tab.value}>
                {tab.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Desktop: Tabs */}
      <div className="hidden sm:block">
        <Tabs value={currentTab} onValueChange={(value) => onTabChange(value as LibraryTab)}>
          <TabsList>
            {tabs.map((tab) => (
              <TabsTrigger key={tab.value} value={tab.value}>
                {tab.name}
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      </div>

      {/* View mode toggle */}
      <div className="flex items-center gap-1 rounded-lg bg-muted p-1">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onViewModeChange('list')}
          className={cn(
            "h-8 w-8 p-0",
            viewMode === 'list' && "bg-background shadow-sm"
          )}
        >
          <List className="h-4 w-4" />
          <span className="sr-only">List view</span>
        </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onViewModeChange('grid')}
          className={cn(
            "h-8 w-8 p-0",
            viewMode === 'grid' && "bg-background shadow-sm"
          )}
        >
          <LayoutGrid className="h-4 w-4" />
          <span className="sr-only">Grid view</span>
        </Button>
      </div>
    </div>
  )
}
