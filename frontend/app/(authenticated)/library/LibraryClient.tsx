'use client'

import type { MediaListItem as Media } from '@/api/models'
import { useAllMedia, useDeleteMedia } from '@/hooks/media'
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle } from '@headlessui/react'
import { CheckIcon } from '@heroicons/react/20/solid'
import { usePathname } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'react-hot-toast'
import CreateMediaForm from './components/CreateMediaForm'
import LibraryHeader from './components/LibraryHeader'
import LibraryTabs from './components/LibraryTabs'
import MediaDetails from './components/MediaDetails'
import MediaGallery from './components/MediaGallery'

type LibraryClientProps = {
  mediaType?: 'audio' | 'image';
}

export default function LibraryClient({ mediaType }: LibraryClientProps = {}) {
  const [selectedItem, setSelectedItem] = useState<Media | null>(null);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const pathname = usePathname();

  // Fetch media data with optional filtering by type
  const { data: mediaResponse, isLoading, error } = useAllMedia(
    mediaType ? { mediaType } : undefined
  );
  const deleteMediaMutation = useDeleteMedia();
  const mediaItems = mediaResponse || [];

  const handleItemClick = (item: Media) => {
    if (selectedItem?.id === item.id) {
      setSelectedItem(null);
      return;
    }
    setSelectedItem(item);
  }

  const closeSidebar = () => {
    setSelectedItem(null);
  }

  const handleAddFile = () => {
    setCreateModalOpen(true);
  }

  const handleUploadClick = () => {
    // TODO: Implement upload functionality
    console.log('Upload clicked');
  }

  const handleDownload = () => {
    if (!selectedItem) return;

    // Use processed URL if available, otherwise use original file URL
    const downloadUrl = selectedItem.processedUrl || selectedItem.fileUrl;

    if (!downloadUrl) {
      toast.error(`No download URL available for: ${selectedItem.fileName || 'item'}`);
      return;
    }

    // Create a temporary anchor element to trigger the download
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = selectedItem.fileName || 'download';
    link.target = '_blank';

    // Append to body, click, and remove
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  const handleDelete = () => {
    if (selectedItem) {
      if (selectedItem.id) {
        deleteMediaMutation.mutate(selectedItem.id);
      }
      setSelectedItem(null);
    }
  }

  return (
    <>
      <div className="flex h-full bg-base">
        {/* Content area */}
        <div className="flex flex-1 flex-col overflow-hidden">
          <LibraryHeader onAddFile={handleAddFile} />

          {/* Main content */}
          <div className="flex flex-1 items-stretch overflow-hidden">
            <main className={selectedItem ? "lg:flex-1" : "flex-1"}>
              <div className="mx-auto max-w-7xl px-4 pt-8 sm:px-6 lg:px-8">
                <div className="flex">
                  <h1 className="flex-1 text-2xl font-bold text-highlight">Library</h1>
                </div>

                <LibraryTabs currentTab={pathname} />

                {isLoading ? (
                  <div className="mt-24 text-center">
                    <p className="text-base">Loading media...</p>
                  </div>
                ) : error ? (
                  <div className="mt-24 text-center">
                    <p className="text-base text-red-500">Error loading media: {error.message}</p>
                  </div>
                ) : (
                  <MediaGallery
                    items={mediaItems}
                    selectedItem={selectedItem}
                    onItemClick={handleItemClick}
                    onUploadClick={handleUploadClick}
                  />
                )}
              </div>
            </main>

            {/* Details sidebar */}
            {selectedItem && (
              <MediaDetails
                item={selectedItem}
                onClose={closeSidebar}
                onDownload={handleDownload}
                onDelete={handleDelete}
              />
            )}
          </div>
        </div>
      </div>

      {/* Create Media Modal */}
      <Dialog open={createModalOpen} onClose={() => setCreateModalOpen(false)} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-elevated px-4 pt-5 pb-4 text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-lg sm:p-6 data-closed:sm:translate-y-0 data-closed:sm:scale-95"
            >
              <div>
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-accent">
                  <CheckIcon aria-hidden="true" className="size-6 text-inverted" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-highlight">
                    Add new media
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-base">
                      Upload a new image or audio file to your library. The file will be processed automatically.
                    </p>
                  </div>
                </div>
                <CreateMediaForm setOpen={setCreateModalOpen} />
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>
    </>
  )
}
