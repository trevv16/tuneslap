'use client'

import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'

type LibraryHeaderProps = {
  onAddFile?: () => void
}

export default function LibraryHeader({ onAddFile }: LibraryHeaderProps) {
  return (
    <div className="flex items-center justify-between pb-4">
      <h1 className="text-2xl font-bold">Library</h1>
      <Button onClick={onAddFile} size="sm">
        <Plus className="mr-2 h-4 w-4" />
        Add File
      </Button>
    </div>
  )
}
