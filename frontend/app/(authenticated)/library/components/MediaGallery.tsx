'use client'

import EmptyState from '@/components/EmptyState';
import type { MediaListItem as Media } from '@/api/models';
import MediaGridView from './MediaGridView';
import MediaListView from './MediaListView';

type MediaGalleryProps = {
  items: Media[];
  selectedItem: Media | null;
  onItemClick: (item: Media) => void;
  onUploadClick?: () => void;
  viewMode?: 'grid' | 'list';
}

export default function MediaGallery({ items, selectedItem, onItemClick, onUploadClick, viewMode = 'grid' }: MediaGalleryProps) {
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
      {viewMode === 'grid' ? (
        <MediaGridView
          items={items}
          selectedItem={selectedItem}
          onItemClick={onItemClick}
        />
      ) : (
        <MediaListView
          items={items}
          selectedItem={selectedItem}
          onItemClick={onItemClick}
        />
      )}
    </section>
  )
}
