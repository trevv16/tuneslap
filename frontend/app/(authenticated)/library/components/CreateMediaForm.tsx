'use client'

import type { MediaListItem } from '@/api/models'
import DemoBanner from '@/components/DemoBanner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Progress } from '@/components/ui/progress'
import { Textarea } from '@/components/ui/textarea'
import { MediaCreateFormData, useMediaCreate, useMediaStats } from '@/hooks/media'
import { cn } from '@/lib/utils'
import { formatBytes } from '@/utils/helpers'
import { CloudUpload } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import { toast } from 'sonner'
import ProcessMediaModal from './ProcessMediaModal'

type CreateMediaFormProps = {
  setOpen: (open: boolean) => void
}

export default function CreateMediaForm({ setOpen }: CreateMediaFormProps) {
  const [dragActive, setDragActive] = useState(false)
  const [createdMedia, setCreatedMedia] = useState<MediaListItem | null>(null)
  const [showProcessModal, setShowProcessModal] = useState(false)

  const {
    form,
    handleSubmit,
    isSubmitting
  } = useMediaCreate()

  const { data: mediaStats } = useMediaStats()

  const selectedFile = form.watch('file')
  const fileName = form.watch('fileName')

  const storageError = useMemo(() => {
    if (!selectedFile || !mediaStats?.data) return null

    const availableStorage = mediaStats.data?.availableStorage
    if (availableStorage === -1 || availableStorage === undefined) return null

    if (selectedFile.size > availableStorage) {
      return `File size (${formatBytes(selectedFile.size)}) exceeds available storage (${formatBytes(availableStorage)})`
    }

    return null
  }, [selectedFile, mediaStats])

  const getFileExtension = (file: File | undefined): string => {
    if (!file) return ''

    const mimeToExtension: Record<string, string> = {
      'audio/mp3': '.mp3',
      'audio/mpeg': '.mp3',
      'audio/wav': '.wav',
      'audio/webm': '.webm',
      'audio/ogg': '.ogg',
      'audio/aac': '.aac',
      'image/jpeg': '.jpg',
      'image/jpg': '.jpg',
      'image/png': '.png',
      'image/gif': '.gif',
      'image/webp': '.webp',
      'image/svg+xml': '.svg',
    }

    const mimeExtension = mimeToExtension[file.type]
    if (mimeExtension) {
      return mimeExtension
    }

    const fileName = file.name
    const lastDotIndex = fileName.lastIndexOf('.')
    if (lastDotIndex !== -1) {
      return fileName.substring(lastDotIndex)
    }

    return ''
  }

  const fileExtension = getFileExtension(selectedFile)

  useEffect(() => {
    if (selectedFile && !fileName) {
      const baseName = selectedFile.name.replace(/\.[^/.]+$/, "")
      form.setValue('fileName', baseName)
    }
  }, [selectedFile, fileName, form])

  const handleFormSubmit = async (data: MediaCreateFormData) => {
    const result = await handleSubmit(data)
    if (result.success && result.data) {
      toast.success('Media uploaded! Configure processing options.')
      const mediaItem = {
        id: result.data.id,
        authorId: result.data.authorId,
        mediaType: result.data.mediaType as MediaListItem['mediaType'],
        fileName: result.data.fileName,
        description: result.data.description,
        fileUrl: result.data.fileUrl,
        processedUrl: result.data.processedUrl,
        waveformUrl: result.data.waveformUrl,
        contentType: result.data.contentType as MediaListItem['contentType'],
        fileSize: result.data.fileSize,
        status: result.data.status as MediaListItem['status'],
        processingParams: result.data.processingParams,
        processingActivity: result.data.processingActivity,
        dimensions: result.data.dimensions,
        duration: result.data.duration,
        createdAt: result.data.createdAt,
        updatedAt: result.data.updatedAt,
      } as MediaListItem
      setCreatedMedia(mediaItem)
      setShowProcessModal(true)
      form.reset()
    } else {
      const errorMessage = result.error instanceof Error
        ? result.error.message
        : typeof result.error === 'string'
          ? result.error
          : 'Failed to create media'
      toast.error(errorMessage)
    }
  }

  const handleProcessModalClose = () => {
    setShowProcessModal(false)
    setCreatedMedia(null)
    setOpen(false)
  }

  const handleProcessSuccess = () => {
    setShowProcessModal(false)
    setCreatedMedia(null)
    setOpen(false)
  }

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true)
    } else if (e.type === "dragleave") {
      setDragActive(false)
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0]
      form.setValue('file', file)

      if (!form.getValues('fileName')) {
        const fileName = file.name.replace(/\.[^/.]+$/, "")
        form.setValue('fileName', fileName)
      }
    }
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0]
      form.setValue('file', file)

      if (!form.getValues('fileName')) {
        const fileName = file.name.replace(/\.[^/.]+$/, "")
        form.setValue('fileName', fileName)
      }
    }
  }

  return (
    <form onSubmit={form.handleSubmit(handleFormSubmit)} className="space-y-4">
      <DemoBanner message="Demo mode: Files will be deleted within one hour. Max file size: 10MB. Max uploads: 5." />
      
      {/* File Upload */}
      <div>
        <Label>File</Label>
        <div
          className={cn(
            "relative border-2 border-dashed rounded-lg p-6 text-center mt-2 transition-colors",
            dragActive
              ? 'border-primary bg-primary/5'
              : 'border-muted hover:border-muted-foreground/50'
          )}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          {selectedFile ? (
            <div className="space-y-2">
              <CloudUpload className="mx-auto h-12 w-12 text-primary" />
              <p className="text-sm font-medium">{selectedFile.name}</p>
              <p className="text-xs text-muted-foreground">
                {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
              </p>
              <Button
                type="button"
                variant="link"
                className="text-destructive hover:text-destructive/80"
                onClick={() => form.setValue('file', undefined as unknown as File)}
              >
                Remove file
              </Button>
            </div>
          ) : (
            <div className="space-y-2">
              <CloudUpload className="mx-auto h-12 w-12 text-muted-foreground" />
              <p className="text-sm text-muted-foreground">
                Drag and drop your file here, or{' '}
                <label className="text-primary hover:text-primary/80 cursor-pointer">
                  browse
                  <input
                    type="file"
                    className="hidden"
                    accept="image/*,audio/*"
                    onChange={handleFileSelect}
                  />
                </label>
              </p>
              <p className="text-xs text-muted-foreground">
                Supports images and audio files
              </p>
            </div>
          )}
        </div>
        {form.formState.errors.file && (
          <p className="mt-1 text-sm text-destructive">{form.formState.errors.file.message}</p>
        )}
        {storageError && (
          <p className="mt-1 text-sm text-destructive">{storageError}</p>
        )}
        {isSubmitting && (
          <div className="mt-2 space-y-1">
            <div className="flex items-center justify-between text-sm text-muted-foreground">
              <span>Uploading file...</span>
              <span>50%</span>
            </div>
            <Progress value={50} />
          </div>
        )}
      </div>

      {/* File Name */}
      <div>
        <Label htmlFor="fileName">File Name</Label>
        <div className="relative mt-2">
          <Input
            {...form.register("fileName")}
            placeholder="Enter file name"
            className={cn(fileExtension && "pr-16")}
          />
          {fileExtension && (
            <div className="absolute inset-y-0 right-0 flex items-center pr-3">
              <span className="text-muted-foreground text-sm font-mono">
                {fileExtension}
              </span>
            </div>
          )}
        </div>
        {form.formState.errors.fileName && (
          <p className="mt-1 text-sm text-destructive">{form.formState.errors.fileName.message}</p>
        )}
        {fileExtension && (
          <p className="mt-1 text-xs text-muted-foreground">
            Final filename will be: <span className="font-mono">{fileName || 'filename'}{fileExtension}</span>
          </p>
        )}
      </div>

      {/* Description */}
      <div>
        <Label htmlFor="description">Description (Optional)</Label>
        <Textarea
          {...form.register("description")}
          className="mt-2"
          rows={3}
          placeholder="Enter description"
        />
        {form.formState.errors.description && (
          <p className="mt-1 text-sm text-destructive">{form.formState.errors.description.message}</p>
        )}
      </div>

      {/* Submit Buttons */}
      <div className="flex gap-3 pt-2">
        <Button
          type="button"
          variant="outline"
          onClick={() => setOpen(false)}
          className="flex-1"
        >
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={isSubmitting || !!storageError}
          className="flex-1"
        >
          {isSubmitting ? 'Uploading...' : 'Upload Media'}
        </Button>
      </div>

      {/* Process Media Modal */}
      {createdMedia && (
        <ProcessMediaModal
          media={createdMedia}
          open={showProcessModal}
          onClose={handleProcessModalClose}
          onSuccess={handleProcessSuccess}
        />
      )}
    </form>
  )
}
