'use client'

import EmptyState from '@/components/EmptyState';
import { classNames } from '@/utils/helpers';

import type { MediaListItem as Media } from '@/api/models';

type MediaGalleryProps = {
  items: Media[];
  selectedItem: Media | null;
  onItemClick: (item: Media) => void;
  onUploadClick?: () => void;
}

export default function MediaGallery({ items, selectedItem, onItemClick, onUploadClick }: MediaGalleryProps) {
  if (items.length === 0) {
    return (
      <div className="mt-24">
        <EmptyState
          title="Create Your First Media"
          description="Get started by uploading a file."
          buttonText="Upload Media"
          buttonOnClick={onUploadClick || (() => { })}
        />
      </div>
    )
  }

  return (
    <section aria-labelledby="gallery-heading" className="mt-8 pb-16">
      <h2 id="gallery-heading" className="sr-only">
        Recently viewed
      </h2>
      <ul className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 sm:gap-x-6 md:grid-cols-4 lg:grid-cols-3 xl:grid-cols-4 xl:gap-x-8">
        {items.map((item) => (
          <button
            key={item.id}
            className="relative"
            onClick={() => onItemClick(item)}
          >
            <div
              className={classNames(
                item.id === selectedItem?.id
                  ? 'ring-2 ring-primary-500 ring-offset-2'
                  : 'focus-within:ring-2 focus-within:ring-primary-500 focus-within:ring-offset-2 focus-within:ring-offset-gray-100',
                'group block w-full overflow-hidden rounded-lg',
              )}
            >
              {item.mediaType === 'image' ? (
                <img
                  alt=""
                  src={item.fileUrl ? item.fileUrl : "/defaultKey.png"}
                  className={classNames(
                    item.id === selectedItem?.id ? '' : 'group-hover:opacity-75',
                    'pointer-events-none aspect-10/7 object-cover',
                  )}
                />
              ) : (
                <div className={classNames(
                  item.id === selectedItem?.id ? '' : 'group-hover:opacity-75',
                  'pointer-events-none aspect-10/7 bg-gray-200 flex items-center justify-center',
                )}>
                  <div className="text-gray-500 text-2xl">🎵</div>
                </div>
              )}
            </div>
            <p className="pointer-events-none mt-2 block truncate text-sm font-medium text-base">
              {item.fileName}
            </p>
            <p className="pointer-events-none block text-sm font-medium text-neutral-300">
              {item.fileSize ? (item.fileSize / 1024 / 1024).toFixed(1) : 'N/A'} MB
            </p>
          </button>
        ))}
      </ul>
    </section>
  )
} 