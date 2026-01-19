'use client'

import type { LibraryView } from '@/hooks/useQueryParams'

type MediaGallerySkeletonProps = {
  count?: number;
  viewMode?: LibraryView;
}

function GridSkeleton({ count }: { count: number }) {
  return (
    <ul className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 sm:gap-x-6 md:grid-cols-4 lg:grid-cols-3 xl:grid-cols-4 xl:gap-x-8">
      {Array.from({ length: count }).map((_, index) => (
        <li key={index}>
          <div className="relative w-full animate-pulse">
            {/* Thumbnail skeleton */}
            <div className="block w-full overflow-hidden rounded-lg">
              <div className="aspect-10/7 bg-gray-200 w-full" />
            </div>
            {/* File name skeleton */}
            <div className="mt-2 h-4 w-3/4 rounded bg-gray-200" />
            {/* File size skeleton */}
            <div className="mt-1 h-4 w-1/3 rounded bg-gray-200" />
          </div>
        </li>
      ))}
    </ul>
  )
}

function ListSkeleton({ count }: { count: number }) {
  return (
    <ul className="flex flex-col space-y-1">
      {Array.from({ length: count }).map((_, index) => (
        <li key={index}>
          <div className="flex w-full items-center gap-x-4 rounded-lg p-3 animate-pulse">
            {/* Thumbnail skeleton */}
            <div className="h-12 w-12 flex-none rounded-lg bg-gray-200" />
            {/* Text content skeleton */}
            <div className="flex-auto min-w-0">
              {/* File name skeleton */}
              <div className="h-4 w-2/3 rounded bg-gray-200" />
              {/* Media type skeleton */}
              <div className="mt-1.5 h-3 w-12 rounded bg-gray-200" />
            </div>
            {/* File size skeleton */}
            <div className="h-4 w-16 flex-none rounded bg-gray-200" />
          </div>
        </li>
      ))}
    </ul>
  )
}

export default function MediaGallerySkeleton({ count = 8, viewMode = 'grid' }: MediaGallerySkeletonProps) {
  return (
    <section className="mt-8 pb-16">
      {viewMode === 'grid' ? (
        <GridSkeleton count={count} />
      ) : (
        <ListSkeleton count={count} />
      )}
    </section>
  );
}
