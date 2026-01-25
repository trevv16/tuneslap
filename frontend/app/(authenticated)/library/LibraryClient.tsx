'use client'

import type { MediaListItem as Media } from '@/api/models'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useAllMedia, useDeleteMedia } from '@/hooks/media'
import { useLibraryParams } from '@/hooks/useQueryParams'
import { Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import CreateMediaForm from './components/CreateMediaForm'
import LibraryHeader from './components/LibraryHeader'
import LibraryTabs from './components/LibraryTabs'
import MediaDetails from './components/MediaDetails'
import MediaGallery from './components/MediaGallery'
import MediaGallerySkeleton from './components/MediaGallerySkeleton'

export default function LibraryClient() {
  const [selectedItem, setSelectedItem] = useState<Media | null>(null)
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [hasMounted, setHasMounted] = useState(false)

  const { tab, view, mediaType, setTab, setView } = useLibraryParams()

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setHasMounted(true)
  }, [])

  const { data: mediaResponse, isLoading, error } = useAllMedia({ mediaType })
  const deleteMediaMutation = useDeleteMedia()
  const mediaItems = mediaResponse || []

  const handleItemClick = (item: Media) => {
    if (selectedItem?.id === item.id) {
      setSelectedItem(null)
      return
    }
    setSelectedItem(item)
  }

  const closeSidebar = () => {
    setSelectedItem(null)
  }

  const handleAddFile = () => {
    setCreateModalOpen(true)
  }

  const handleUploadClick = () => {
    setCreateModalOpen(true)
  }

  const handleDownload = () => {
    if (!selectedItem) return

    const downloadUrl = selectedItem.processedUrl || selectedItem.fileUrl

    if (!downloadUrl) {
      toast.error(`No download URL available for: ${selectedItem.fileName || 'item'}`)
      return
    }

    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = selectedItem.fileName || 'download'
    link.target = '_blank'

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const handleDelete = () => {
    if (selectedItem) {
      if (selectedItem.id) {
        deleteMediaMutation.mutate(selectedItem.id)
      }
      setSelectedItem(null)
    }
  }

  return (
    <>
      <LibraryHeader onAddFile={handleAddFile} />

      <LibraryTabs
        currentTab={tab}
        onTabChange={setTab}
        viewMode={view}
        onViewModeChange={setView}
      />

      {!hasMounted || isLoading ? (
        <MediaGallerySkeleton viewMode={view} />
      ) : error ? (
        <div className="mt-24 text-center">
          <p className="text-destructive">Error loading media: {error.message}</p>
        </div>
      ) : (
        <MediaGallery
          items={mediaItems}
          selectedItem={selectedItem}
          onItemClick={handleItemClick}
          onUploadClick={handleUploadClick}
          viewMode={view}
        />
      )}

      {/* Details drawer - renders as overlay via Sheet portal */}
      {selectedItem && (
        <MediaDetails
          item={selectedItem}
          onClose={closeSidebar}
          onDownload={handleDownload}
          onDelete={handleDelete}
        />
      )}

      <Dialog open={createModalOpen} onOpenChange={setCreateModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary">
              <Plus className="h-6 w-6 text-primary-foreground" />
            </div>
            <DialogTitle className="text-center">Add new media</DialogTitle>
            <DialogDescription className="text-center">
              Upload a new image or audio file to your library. The file will be processed automatically.
            </DialogDescription>
          </DialogHeader>
          <CreateMediaForm setOpen={setCreateModalOpen} />
        </DialogContent>
      </Dialog>
    </>
  )
}
