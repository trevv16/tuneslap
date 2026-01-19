import { Skeleton } from "@/components/ui/skeleton"
import { Metadata } from "next"
import { Suspense } from "react"
import PageTemplate from "../PageTemplate"
import LibraryClient from "./LibraryClient"

export const metadata: Metadata = {
  title: "Library",
  description: "Library",
}

function LibraryLoadingSkeleton() {
  return (
    <div className="space-y-6">
      {/* Header skeleton */}
      <div className="flex items-center justify-between">
        <Skeleton className="h-8 w-32" />
        <Skeleton className="h-10 w-28" />
      </div>

      {/* Tabs skeleton */}
      <div className="flex items-center justify-between border-b pb-4">
        <div className="flex gap-2">
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-10 w-20" />
          <Skeleton className="h-10 w-20" />
        </div>
        <div className="flex gap-1">
          <Skeleton className="h-8 w-8" />
          <Skeleton className="h-8 w-8" />
        </div>
      </div>

      {/* Grid skeleton */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="space-y-2">
            <Skeleton className="aspect-square w-full rounded-lg" />
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-3 w-1/2" />
          </div>
        ))}
      </div>
    </div>
  )
}

export default function LibraryPage() {
  return (
    <PageTemplate>
      <Suspense fallback={<LibraryLoadingSkeleton />}>
        <LibraryClient />
      </Suspense>
    </PageTemplate>
  )
}
