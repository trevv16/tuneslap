'use client'

import type { MediaListItem } from '@/api/models'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  useAudioProcessingForm,
  useImageProcessingForm,
  AudioProcessingFormData,
  ImageProcessingFormData,
} from '@/hooks/media'
import { Cog } from 'lucide-react'
import { toast } from 'sonner'
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
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && onClose()}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary">
              <Cog className="h-5 w-5 text-primary-foreground" />
            </div>
            <div>
              <DialogTitle>Process Media</DialogTitle>
              <DialogDescription>
                Configure processing options for your {isAudio ? 'audio' : 'image'} file
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        {/* Preview */}
        <div className="mb-4">
          <MediaPreview media={media} className="mb-4" />
        </div>

        {/* Processing Form */}
        {isAudio && (
          <form onSubmit={audioForm.form.handleSubmit(handleAudioSubmit)}>
            <AudioProcessingForm form={audioForm.form} />
            
            {/* Actions */}
            <div className="mt-6 flex gap-3">
              <Button
                type="button"
                variant="outline"
                onClick={handleSkip}
                disabled={isSubmitting}
                className="flex-1"
              >
                Skip for now
              </Button>
              <Button
                type="submit"
                disabled={isSubmitting}
                className="flex-1"
              >
                {isSubmitting ? 'Starting...' : 'Start Processing'}
              </Button>
            </div>
          </form>
        )}

        {isImage && (
          <form onSubmit={imageForm.form.handleSubmit(handleImageSubmit)}>
            <ImageProcessingForm form={imageForm.form} />
            
            {/* Actions */}
            <div className="mt-6 flex gap-3">
              <Button
                type="button"
                variant="outline"
                onClick={handleSkip}
                disabled={isSubmitting}
                className="flex-1"
              >
                Skip for now
              </Button>
              <Button
                type="submit"
                disabled={isSubmitting}
                className="flex-1"
              >
                {isSubmitting ? 'Starting...' : 'Start Processing'}
              </Button>
            </div>
          </form>
        )}

        {!isAudio && !isImage && (
          <div className="text-center py-8">
            <p className="text-muted-foreground">Unsupported media type</p>
            <Button
              type="button"
              onClick={onClose}
              className="mt-4"
            >
              Close
            </Button>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
