'use client'

import { classNames } from '@/utils/helpers'
import { ChevronDownIcon } from '@heroicons/react/16/solid'
import { Bars4Icon, Squares2X2Icon as Squares2X2IconMini } from '@heroicons/react/20/solid'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

const tabs = [
  { name: 'All Media', href: '/library' },
  { name: 'Audio', href: '/library/audio' },
  { name: 'Images', href: '/library/images' },
]

type LibraryTabsProps = {
  currentTab?: string;
}

export default function LibraryTabs({ currentTab }: LibraryTabsProps) {
  const pathname = usePathname();
  const activeTab = currentTab || pathname;
  
  const updatedTabs = tabs.map(tab => ({
    ...tab,
    current: tab.href === activeTab
  }))

  return (
    <div className="mt-3 sm:mt-2">
      <div className="grid grid-cols-1 sm:hidden">
        {/* Use an "onChange" listener to redirect the user to the selected tab URL. */}
        <select
          defaultValue="Recently Viewed"
          aria-label="Select a tab"
          className="col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-2 pr-8 pl-3 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600"
        >
          <option>Recently Viewed</option>
          <option>Recently Added</option>
        </select>
        <ChevronDownIcon
          aria-hidden="true"
          className="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end fill-gray-500"
        />
      </div>
      <div className="hidden sm:block">
        <div className="flex items-center border-b border-gray-200">
          <nav aria-label="Tabs" className="-mb-px flex flex-1 space-x-6 xl:space-x-8">
            {updatedTabs.map((tab) => (
              <Link
                key={tab.name}
                href={tab.href}
                aria-current={tab.current ? 'page' : undefined}
                className={classNames(
                  tab.current
                    ? 'border-primary-500 text-accent'
                    : 'border-transparent text-base hover:border-highlight hover:text-highlight',
                  'border-b-2 px-1 py-4 text-sm font-medium whitespace-nowrap',
                )}
              >
                {tab.name}
              </Link>
            ))}
          </nav>
          <div className="ml-6 hidden items-center rounded-lg bg-highlight p-0.5 sm:flex">
            <button
              type="button"
              className="rounded-md p-1.5 bg-highlight hover:bg-elevated hover:text-highlight hover:shadow-xs focus:ring-2 focus:ring-primary-500 focus:outline-hidden focus:ring-inset"
            >
              <Bars4Icon aria-hidden="true" className="size-5" />
              <span className="sr-only">Use list view</span>
            </button>
            <button
              type="button"
              className="ml-0.5 rounded-md p-1.5 bg-elevated text-highlight shadow-xs focus:ring-2 focus:ring-primary-500 focus:outline-hidden focus:ring-inset"
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