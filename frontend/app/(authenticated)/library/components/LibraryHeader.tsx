'use client'

import { MagnifyingGlassIcon, PlusIcon } from '@heroicons/react/20/solid';

type LibraryHeaderProps = {
  onAddFile?: () => void;
}

export default function LibraryHeader({ onAddFile }: LibraryHeaderProps) {
  return (
    <header className="w-full">
      <div className="relative flex h-16 shrink-0 border-b border-highlight bg-muted shadow-xs">
        <div className="flex flex-1 justify-between px-4 sm:px-6">
          <div className="flex flex-1">
            <form action="#" method="GET" className="grid flex-1 grid-cols-1">
              <input
                name="search"
                type="search"
                placeholder="Search"
                aria-label="Search"
                className="col-start-1 row-start-1 block size-full bg-muted pl-8 text-base outline-hidden placeholder:text-base placeholder:text-sm sm:text-sm/6"
              />
              <MagnifyingGlassIcon
                aria-hidden="true"
                className="pointer-events-none col-start-1 row-start-1 size-5 self-center text-gray-400"
              />
            </form>
          </div>
          <div className="ml-2 flex items-center space-x-4 sm:ml-6 sm:space-x-6">
            <button
              type="button"
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onAddFile?.();
              }}
              className="relative rounded-full bg-primary-600 p-1.5 text-white hover:bg-primary-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
            >
              <span className="absolute -inset-1.5" />
              <PlusIcon aria-hidden="true" className="size-5" />
              <span className="sr-only">Add file</span>
            </button>
          </div>
        </div>
      </div>
    </header>
  )
} 