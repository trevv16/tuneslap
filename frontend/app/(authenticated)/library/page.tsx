import { Metadata } from "next";
import { Suspense } from "react";
import PageTemplate from "../PageTemplate";
import LibraryClient from "./LibraryClient";
import MediaGallerySkeleton from "./components/MediaGallerySkeleton";

export const metadata: Metadata = {
  title: "Library",
  description: "Library",
};

function LibraryLoadingSkeleton() {
  return (
    <div className="flex h-full bg-base">
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Header skeleton */}
        <div className="flex items-center justify-between border-b border-neutral-700 px-4 py-4 sm:px-6 lg:px-8">
          <div className="h-8 w-32 animate-pulse rounded bg-gray-200" />
          <div className="h-10 w-28 animate-pulse rounded bg-gray-200" />
        </div>

        {/* Main content skeleton */}
        <main className="flex-1">
          <div className="mx-auto max-w-7xl px-4 pt-8 sm:px-6 lg:px-8">
            {/* Title skeleton */}
            <div className="h-8 w-24 animate-pulse rounded bg-gray-200" />
            
            {/* Tabs skeleton */}
            <div className="mt-4 flex gap-4">
              <div className="h-10 w-20 animate-pulse rounded bg-gray-200" />
              <div className="h-10 w-20 animate-pulse rounded bg-gray-200" />
              <div className="h-10 w-20 animate-pulse rounded bg-gray-200" />
            </div>

            <MediaGallerySkeleton />
          </div>
        </main>
      </div>
    </div>
  );
}

export default function LibraryPage() {
  return (
    <PageTemplate>
      <Suspense fallback={<LibraryLoadingSkeleton />}>
        <LibraryClient />
      </Suspense>
    </PageTemplate>
  );
}
