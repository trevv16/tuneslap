'use client'

import { Skeleton } from '@/components/ui/skeleton'
import type { LibraryView } from '@/hooks/useQueryParams'

type MediaGallerySkeletonProps = {
  count?: number
  viewMode?: LibraryView
}

function GridSkeleton({ count }: { count: number }) {
  return (
    <ul className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-3 xl:grid-cols-4 xl:gap-x-8">
      {Array.from({ length: count }).map((_, index) => (
        <li key={index} className="space-y-2">
          <Skeleton className="aspect-[10/7] w-full rounded-lg" />
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-4 w-1/3" />
        </li>
      ))}
    </ul>
  )
}

function ListSkeleton({ count }: { count: number }) {
  return (
    <ul className="flex flex-col space-y-1">
      {Array.from({ length: count }).map((_, index) => (
        <li key={index} className="flex w-full items-center gap-x-4 rounded-lg p-3">
          <Skeleton className="h-12 w-12 flex-none rounded-lg" />
          <div className="flex-auto min-w-0 space-y-2">
            <Skeleton className="h-4 w-2/3" />
            <Skeleton className="h-3 w-12" />
          </div>
          <Skeleton className="h-4 w-16 flex-none" />
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
  )
}
