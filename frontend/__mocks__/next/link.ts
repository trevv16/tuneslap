// Mock Next.js Link component
import React from 'react'

interface LinkProps {
  href: string
  children: React.ReactNode
  className?: string
  onClick?: () => void
  onMouseEnter?: () => void
  prefetch?: boolean
}

function MockLink({ href, children, className, onClick, onMouseEnter, ...props }: LinkProps) {
  return React.createElement(
    'a',
    {
      href,
      className,
      onClick,
      onMouseEnter,
      'data-testid': 'next-link',
      ...props,
    },
    children
  )
}

export default MockLink
