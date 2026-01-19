'use client'

import type { MediaListItem } from '@/api/models'
import {
  useAudioProcessingForm,
  useImageProcessingForm,
  AudioProcessingFormData,
  ImageProcessingFormData,
} from '@/hooks/media'
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle } from '@headlessui/react'
import { CogIcon, XMarkIcon } from '@heroicons/react/20/solid'
import { toast } from 'react-hot-toast'
import AudioProcessingForm from './AudioProcessingForm'
import ImageProcessingForm from './ImageProcessingForm'
import MediaPreview from './MediaPreview'

type ProcessMediaModalProps = {
  media: MediaListItem
  open: boolean
  onClose: () => void
  onSuccess?: () => void
}

export default function ProcessMediaModal({
  media,
  open,
  onClose,
  onSuccess,
}: ProcessMediaModalProps) {
  const isAudio = media.mediaType === 'audio'
  const isImage = media.mediaType === 'image'

  const audioForm = useAudioProcessingForm(media.id || '')
  const imageForm = useImageProcessingForm(media.id || '')

  const handleAudioSubmit = async (data: AudioProcessingFormData) => {
    const result = await audioForm.handleSubmit(data)
    if (result.success) {
      toast.success('Processing started! Check back soon for your processed file.')
      onSuccess?.()
      onClose()
    } else {
      const errorMessage = result.error instanceof Error
        ? result.error.message
        : typeof result.error === 'string'
          ? result.error
          : 'Failed to start processing'
      toast.error(errorMessage)
    }
  }

  const handleImageSubmit = async (data: ImageProcessingFormData) => {
    const result = await imageForm.handleSubmit(data)
    if (result.success) {
      toast.success('Processing started! Check back soon for your processed file.')
      onSuccess?.()
      onClose()
    } else {
      const errorMessage = result.error instanceof Error
        ? result.error.message
        : typeof result.error === 'string'
          ? result.error
          : 'Failed to start processing'
      toast.error(errorMessage)
    }
  }

  const handleSkip = () => {
    toast.success('Media uploaded. You can process it later from the library.')
    onClose()
  }

  const isSubmitting = audioForm.isSubmitting || imageForm.isSubmitting

  return (
    <Dialog open={open} onClose={onClose} className="relative z-50">
      <DialogBackdrop
        transition
        className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
      />

      <div className="fixed inset-0 z-50 w-screen overflow-y-auto">
        <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <DialogPanel
            transition
            className="relative transform overflow-hidden rounded-lg bg-elevated text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-2xl data-closed:sm:translate-y-0 data-closed:sm:scale-95"
          >
            {/* Header */}
            <div className="flex items-center justify-between border-b border-highlight px-6 py-4">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-accent">
                  <CogIcon className="h-5 w-5 text-inverted" />
                </div>
                <div>
                  <DialogTitle as="h3" className="text-lg font-semibold text-highlight">
                    Process Media
                  </DialogTitle>
                  <p className="text-sm text-muted">
                    Configure processing options for your {isAudio ? 'audio' : 'image'} file
                  </p>
                </div>
              </div>
              <button
                onClick={onClose}
                className="rounded-full p-1 text-base hover:bg-muted"
              >
                <XMarkIcon className="h-6 w-6" />
              </button>
            </div>

            {/* Content */}
            <div className="px-6 py-4">
              {/* Preview */}
              <div className="mb-6">
                <MediaPreview media={media} className="mb-4" />
              </div>

              {/* Processing Form */}
              {isAudio && (
                <form onSubmit={audioForm.form.handleSubmit(handleAudioSubmit)}>
                  <AudioProcessingForm form={audioForm.form} />
                  
                  {/* Actions */}
                  <div className="mt-6 flex gap-3">
                    <button
                      type="button"
                      onClick={handleSkip}
                      disabled={isSubmitting}
                      className="flex-1 rounded-md bg-muted px-4 py-2.5 text-sm font-semibold text-base hover:bg-highlight disabled:opacity-50"
                    >
                      Skip for now
                    </button>
                    <button
                      type="submit"
                      disabled={isSubmitting}
                      className="flex-1 rounded-md bg-accent px-4 py-2.5 text-sm font-semibold text-inverted hover:bg-highlight disabled:opacity-50"
                    >
                      {isSubmitting ? 'Starting...' : 'Start Processing'}
                    </button>
                  </div>
                </form>
              )}

              {isImage && (
                <form onSubmit={imageForm.form.handleSubmit(handleImageSubmit)}>
                  <ImageProcessingForm form={imageForm.form} />
                  
                  {/* Actions */}
                  <div className="mt-6 flex gap-3">
                    <button
                      type="button"
                      onClick={handleSkip}
                      disabled={isSubmitting}
                      className="flex-1 rounded-md bg-muted px-4 py-2.5 text-sm font-semibold text-base hover:bg-highlight disabled:opacity-50"
                    >
                      Skip for now
                    </button>
                    <button
                      type="submit"
                      disabled={isSubmitting}
                      className="flex-1 rounded-md bg-accent px-4 py-2.5 text-sm font-semibold text-inverted hover:bg-highlight disabled:opacity-50"
                    >
                      {isSubmitting ? 'Starting...' : 'Start Processing'}
                    </button>
                  </div>
                </form>
              )}

              {!isAudio && !isImage && (
                <div className="text-center py-8">
                  <p className="text-muted">Unsupported media type</p>
                  <button
                    type="button"
                    onClick={onClose}
                    className="mt-4 rounded-md bg-accent px-4 py-2 text-sm font-semibold text-inverted"
                  >
                    Close
                  </button>
                </div>
              )}
            </div>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  )
}
