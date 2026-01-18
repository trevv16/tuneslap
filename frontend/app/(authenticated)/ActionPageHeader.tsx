import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import {
  CheckIcon,
  ChevronDownIcon, DocumentChartBarIcon, LinkIcon,
  PencilIcon, Squares2X2Icon
} from '@heroicons/react/20/solid';

interface ActionPageHeaderProps {
  title: string;
}

export default function ActionPageHeader({ title }: ActionPageHeaderProps) {
  return (
    <div className="lg:flex lg:items-center lg:justify-between">
      <div className="min-w-0 flex-1">
        <h2 className="text-2xl/7 font-bold sm:truncate sm:text-3xl sm:tracking-tight text-base">
          {title}
        </h2>
        <div className="mt-1 flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap sm:space-x-6">

          <div className="mt-2 flex items-center text-sm text-highlight">
            <Squares2X2Icon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
            4 Boards
          </div>
          <div className="mt-2 flex items-center text-sm text-highlight">
            <DocumentChartBarIcon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
            3 MB Used
          </div>
          <div className="mt-2 flex items-center text-sm text-highlight">
            <DocumentChartBarIcon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
            2 MB Available
          </div>
        </div>
      </div>
      <div className="mt-5 flex lg:mt-0 lg:ml-4">
        <span className="hidden sm:block group">
          <button
            type="button"
            className="inline-flex items-center rounded-md bg-elevated px-3 py-2 text-sm font-semibold text-highlight shadow-xs ring-1 border border-muted ring-inset group-hover:bg-highlight group-hover:text-inverted"
          >
            <PencilIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-highlight group-hover:text-inverted" />
            Edit
          </button>
        </span>

        <span className="ml-3 hidden sm:block group">
          <button
            type="button"
            className="inline-flex items-center rounded-md bg-elevated px-3 py-2 text-sm font-semibold text-highlight shadow-xs ring-1 border border-muted ring-inset group-hover:bg-highlight group-hover:text-inverted"
          >
            <LinkIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-highlight group-hover:text-inverted" />
            View
          </button>
        </span>

        <span className="sm:ml-3">
          <button
            type="button"
            className="inline-flex items-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <CheckIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-inverted group-hover:text-inverted" />
            Publish
          </button>
        </span>

        {/* Dropdown */}
        <Menu as="div" className="relative ml-3 sm:hidden">
          <MenuButton className="group inline-flex items-center rounded-md bg-elevated px-3 py-2 text-sm font-semibold text-highlight shadow-xs ring-1 border border-muted ring-inset group-hover:bg-highlight group-hover:text-inverted">
            More
            <ChevronDownIcon aria-hidden="true" className="-mr-1 ml-1.5 size-5 text-highlight group-hover:text-inverted" />
          </MenuButton>

          <MenuItems
            transition
            className="absolute right-0 z-10 mt-2 -mr-1 w-48 origin-top-right rounded-md bg-elevated py-1 shadow-lg ring-1 border border-muted ring-inset transition focus:outline-hidden data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-200 data-enter:ease-out data-leave:duration-75 data-leave:ease-in"
          >
            <MenuItem>
              <a
                href="#"
                className="block px-4 py-2 text-sm text-base data-focus:text-inverted data-focus:bg-highlight"
              >
                Edit
              </a>
            </MenuItem>
            <MenuItem>
              <a
                href="#"
                className="block px-4 py-2 text-sm text-base data-focus:text-inverted data-focus:bg-highlight"
              >
                View
              </a>
            </MenuItem>
          </MenuItems>
        </Menu>
      </div>
    </div>
  )
}
