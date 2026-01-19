'use client'

import type { MediaListItem as Media } from '@/api/models'
import { cn } from '@/lib/utils'

type MediaGridViewProps = {
  items: Media[]
  selectedItem: Media | null
  onItemClick: (item: Media) => void
}

export default function MediaGridView({ items, selectedItem, onItemClick }: MediaGridViewProps) {
  return (
    <ul className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-3 xl:grid-cols-4 xl:gap-x-8">
      {items.map((item) => (
        <li key={item.id}>
          <button
            className="relative w-full text-left"
            onClick={() => onItemClick(item)}
          >
            <div
              className={cn(
                'group block w-full overflow-hidden rounded-lg',
                item.id === selectedItem?.id
                  ? 'ring-2 ring-primary ring-offset-2 ring-offset-background'
                  : 'focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 focus-within:ring-offset-background'
              )}
            >
              {item.mediaType === 'image' ? (
                <img
                  alt=""
                  src={item.fileUrl ? item.fileUrl : "/defaultKey.png"}
                  className={cn(
                    'pointer-events-none aspect-[10/7] object-cover w-full',
                    item.id !== selectedItem?.id && 'group-hover:opacity-75'
                  )}
                />
              ) : (
                <div className={cn(
                  'pointer-events-none aspect-[10/7] bg-muted flex items-center justify-center w-full',
                  item.id !== selectedItem?.id && 'group-hover:opacity-75'
                )}>
                  <div className="text-muted-foreground text-2xl">🎵</div>
                </div>
              )}
            </div>
            <p className="pointer-events-none mt-2 block truncate text-sm font-medium text-foreground">
              {item.fileName}
            </p>
            <p className="pointer-events-none block text-sm font-medium text-muted-foreground">
              {item.fileSize ? (item.fileSize / 1024 / 1024).toFixed(1) : 'N/A'} MB
            </p>
          </button>
        </li>
      ))}
    </ul>
  )
}
