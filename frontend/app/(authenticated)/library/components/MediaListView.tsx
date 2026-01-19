'use client'

import type { MediaListItem as Media } from '@/api/models'
import { cn } from '@/lib/utils'

type MediaListViewProps = {
  items: Media[]
  selectedItem: Media | null
  onItemClick: (item: Media) => void
}

export default function MediaListView({ items, selectedItem, onItemClick }: MediaListViewProps) {
  return (
    <ul className="flex flex-col space-y-1">
      {items.map((item) => (
        <li key={item.id}>
          <button
            onClick={() => onItemClick(item)}
            className={cn(
              'group flex w-full items-center gap-x-4 rounded-lg p-3 text-sm font-medium leading-6 transition-all duration-150',
              item.id === selectedItem?.id
                ? 'bg-primary/10 ring-1 ring-primary'
                : 'hover:bg-accent hover:shadow-sm'
            )}
          >
            <div className={cn(
              'relative h-12 w-12 flex-none overflow-hidden rounded-lg transition-transform duration-150',
              item.id !== selectedItem?.id && 'group-hover:scale-105'
            )}>
              {item.mediaType === 'image' ? (
                <img
                  src={item.fileUrl || "/defaultKey.png"}
                  alt=""
                  className="h-full w-full object-cover"
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-muted">
                  <span className="text-xl">🎵</span>
                </div>
              )}
            </div>
            <div className="flex-auto text-left min-w-0">
              <p className={cn(
                'font-semibold truncate transition-colors duration-150',
                item.id === selectedItem?.id ? 'text-primary' : 'text-foreground group-hover:text-foreground'
              )}>
                {item.fileName}
              </p>
              <p className="text-xs text-muted-foreground mt-0.5">
                {item.mediaType === 'audio' ? 'Audio' : 'Image'}
              </p>
            </div>
            <div className={cn(
              'flex-none text-sm transition-colors duration-150',
              item.id === selectedItem?.id ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'
            )}>
              {item.fileSize ? (item.fileSize / 1024 / 1024).toFixed(2) : 'N/A'} MB
            </div>
          </button>
        </li>
      ))}
    </ul>
  )
}
