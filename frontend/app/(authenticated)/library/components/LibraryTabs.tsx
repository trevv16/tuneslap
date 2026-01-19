'use client'

import { classNames } from '@/utils/helpers'
import { ChevronDownIcon } from '@heroicons/react/16/solid'
import { Bars4Icon, Squares2X2Icon as Squares2X2IconMini } from '@heroicons/react/20/solid'
import type { LibraryTab, LibraryView } from '@/hooks/useQueryParams'

const tabs: { name: string; value: LibraryTab }[] = [
  { name: 'All Media', value: 'all' },
  { name: 'Audio', value: 'audio' },
  { name: 'Images', value: 'images' },
]

type LibraryTabsProps = {
  currentTab: LibraryTab;
  onTabChange: (tab: LibraryTab) => void;
  viewMode: LibraryView;
  onViewModeChange: (mode: LibraryView) => void;
}

export default function LibraryTabs({ currentTab, onTabChange, viewMode, onViewModeChange }: LibraryTabsProps) {
  return (
    <div className="mt-3 sm:mt-2">
      <div className="grid grid-cols-1 sm:hidden">
        <select
          value={currentTab}
          onChange={(e) => onTabChange(e.target.value as LibraryTab)}
          aria-label="Select a tab"
          className="col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-2 pr-8 pl-3 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600"
        >
          {tabs.map((tab) => (
            <option key={tab.value} value={tab.value}>
              {tab.name}
            </option>
          ))}
        </select>
        <ChevronDownIcon
          aria-hidden="true"
          className="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end fill-gray-500"
        />
      </div>
      <div className="hidden sm:block">
        <div className="flex items-center border-b border-gray-200">
          <nav aria-label="Tabs" className="-mb-px flex flex-1 space-x-6 xl:space-x-8">
            {tabs.map((tab) => (
              <button
                key={tab.value}
                type="button"
                onClick={() => onTabChange(tab.value)}
                aria-current={tab.value === currentTab ? 'page' : undefined}
                className={classNames(
                  tab.value === currentTab
                    ? 'border-primary-500 text-accent'
                    : 'border-transparent text-base hover:border-highlight hover:text-highlight',
                  'border-b-2 px-1 py-4 text-sm font-medium whitespace-nowrap',
                )}
              >
                {tab.name}
              </button>
            ))}
          </nav>
          <div className="ml-6 hidden items-center rounded-lg bg-highlight p-0.5 sm:flex">
            <button
              type="button"
              onClick={() => onViewModeChange('list')}
              className={classNames(
                viewMode === 'list'
                  ? 'bg-elevated text-highlight shadow-xs'
                  : 'bg-highlight text-base hover:bg-elevated hover:text-highlight',
                'rounded-md p-1.5 focus:ring-2 focus:ring-primary-500 focus:outline-hidden focus:ring-inset'
              )}
            >
              <Bars4Icon aria-hidden="true" className="size-5" />
              <span className="sr-only">Use list view</span>
            </button>
            <button
              type="button"
              onClick={() => onViewModeChange('grid')}
              className={classNames(
                viewMode === 'grid'
                  ? 'bg-elevated text-highlight shadow-xs'
                  : 'bg-highlight text-base hover:bg-elevated hover:text-highlight',
                'ml-0.5 rounded-md p-1.5 focus:ring-2 focus:ring-primary-500 focus:outline-hidden focus:ring-inset'
              )}
            >
              <Squares2X2IconMini aria-hidden="true" className="size-5" />
              <span className="sr-only">Use grid view</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
