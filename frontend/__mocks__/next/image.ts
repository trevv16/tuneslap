// Mock Next.js Image component
import React from 'react'

interface ImageProps {
  src: string
  alt: string
  width?: number
  height?: number
  fill?: boolean
  priority?: boolean
  className?: string
  onLoad?: () => void
  onError?: () => void
}

function MockImage({
  src,
  alt,
  width,
  height,
  className,
  onLoad,
  onError,
  fill: _fill,
  priority: _priority,
  ...props
}: ImageProps) {
  return React.createElement('img', {
    src,
    alt,
    width,
    height,
    className,
    onLoad,
    onError,
    'data-testid': 'next-image',
    ...props,
  })
}
}

export default MockImage
