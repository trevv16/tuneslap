'use client'

import type { MediaListItem } from '@/api/models'

type MediaPreviewProps = {
  media: MediaListItem
  className?: string
}

export default function MediaPreview({ media, className = '' }: MediaPreviewProps) {
  const fileUrl = media.fileUrl || ''

  if (media.mediaType === 'audio') {
    return (
      <div className={`flex flex-col items-center justify-center ${className}`}>
        <div className="w-full max-w-md p-6 bg-elevated rounded-lg">
          <div className="flex items-center justify-center mb-4">
            <div className="w-24 h-24 bg-muted rounded-full flex items-center justify-center">
              <span className="text-4xl">🎵</span>
            </div>
          </div>
          <p className="text-center text-sm text-highlight font-medium mb-4 truncate">
            {media.fileName || 'Audio file'}
          </p>
          <audio
            controls
            className="w-full"
            src={fileUrl}
            preload="metadata"
          >
            Your browser does not support the audio element.
          </audio>
        </div>
      </div>
    )
  }

  if (media.mediaType === 'image') {
    return (
      <div className={`flex items-center justify-center ${className}`}>
        <div className="max-w-md max-h-80 overflow-hidden rounded-lg">
          <img
            src={fileUrl || '/defaultKey.png'}
            alt={media.fileName || 'Image preview'}
            className="w-full h-full object-contain"
          />
        </div>
      </div>
    )
  }

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <div className="p-6 bg-elevated rounded-lg text-center">
        <p className="text-muted">Preview not available</p>
      </div>
    </div>
  )
}
