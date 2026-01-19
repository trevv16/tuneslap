'use client'

import { classNames } from '@/utils/helpers';
import type { MediaListItem as Media } from '@/api/models';

type MediaListViewProps = {
  items: Media[];
  selectedItem: Media | null;
  onItemClick: (item: Media) => void;
}

export default function MediaListView({ items, selectedItem, onItemClick }: MediaListViewProps) {
  return (
    <ul className="flex flex-col space-y-1">
      {items.map((item) => (
        <li key={item.id}>
          <button
            onClick={() => onItemClick(item)}
            className={classNames(
              item.id === selectedItem?.id
                ? 'bg-primary-500/10 ring-1 ring-primary-500'
                : 'hover:bg-elevated hover:shadow-sm',
              'group flex w-full items-center gap-x-4 rounded-lg p-3 text-sm font-medium leading-6 transition-all duration-150'
            )}
          >
            <div className={classNames(
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
                <div className="flex h-full w-full items-center justify-center bg-gray-200">
                  <span className="text-xl">🎵</span>
                </div>
              )}
            </div>
            <div className="flex-auto text-left min-w-0">
              <p className={classNames(
                'font-semibold truncate transition-colors duration-150',
                item.id === selectedItem?.id ? 'text-primary-400' : 'text-base group-hover:text-highlight'
              )}>
                {item.fileName}
              </p>
              <p className="text-xs text-neutral-500 mt-0.5">
                {item.mediaType === 'audio' ? 'Audio' : 'Image'}
              </p>
            </div>
            <div className={classNames(
              'flex-none text-sm transition-colors duration-150',
              item.id === selectedItem?.id ? 'text-primary-400' : 'text-neutral-400 group-hover:text-neutral-300'
            )}>
              {item.fileSize ? (item.fileSize / 1024 / 1024).toFixed(2) : 'N/A'} MB
            </div>
          </button>
        </li>
      ))}
    </ul>
  );
}
