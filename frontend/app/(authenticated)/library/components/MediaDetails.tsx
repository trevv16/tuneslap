'use client'

import type { MediaListItem as Media } from '@/api/models'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { Textarea } from '@/components/ui/textarea'
import { MediaEditFormData, useMediaEdit } from '@/hooks/media'
import { AlertTriangle, Pencil } from 'lucide-react'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'

type MediaDetailsProps = {
  item: Media
  onClose: () => void
  onDownload?: () => void
  onDelete?: () => void
}

export default function MediaDetails({
  item,
  onClose,
  onDownload,
  onDelete
}: MediaDetailsProps) {
  return (
    <Sheet open={true} onOpenChange={(open) => !open && onClose()}>
      <SheetContent side="right" className="w-full max-w-sm sm:max-w-md lg:max-w-lg overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Details</SheetTitle>
        </SheetHeader>
        <div className="px-6 py-4">
          <MediaDetailsContent
            item={item}
            onDownload={onDownload}
            onDelete={onDelete}
          />
        </div>
      </SheetContent>
    </Sheet>
  )
}

function MediaDetailsContent({
  item,
  onDownload,
  onDelete
}: {
  item: Media
  onDownload?: () => void
  onDelete?: () => void
}) {
  const [isEditing, setIsEditing] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)

  const {
    form,
    handleSubmit,
    isSubmitting
  } = useMediaEdit(item.id ?? '', {
    fileName: item.fileName ?? '',
    description: item.description,
  })

  useEffect(() => {
    form.reset({
      description: item.description,
    })
  }, [item.description, form])

  const handleEditToggle = () => {
    if (isEditing) {
      form.reset({
        description: item.description,
      })
      setIsEditing(false)
    } else {
      setIsEditing(true)
    }
  }

  const onFormSubmit = async (data: MediaEditFormData) => {
    const result = await handleSubmit(data)
    if (result.success) {
      toast.success('Media updated successfully')
      setIsEditing(false)
      form.reset({
        description: data.description,
      })
    } else {
      toast.error('Failed to update media')
    }
  }

  const handleDeleteConfirm = async () => {
    try {
      if (onDelete) {
        await onDelete()
        toast.success('Media deleted successfully!')
        setDeleteModalOpen(false)
      }
    } catch (error) {
      toast.error('Failed to delete media. Please try again.')
      console.error('Delete media error:', error)
    }
  }

  return (
    <>
      {/* Media Preview */}
      <div>
        {item.mediaType === 'image' ? (
          <img
            alt=""
            src={(item.fileUrl && item.fileUrl !== "") ? item.fileUrl : "/defaultKey.png"}
            className="block aspect-[10/7] w-full rounded-lg object-cover"
          />
        ) : (
          <div className="block aspect-[10/7] w-full rounded-lg bg-muted flex items-center justify-center">
            <div className="text-muted-foreground text-6xl">🎵</div>
          </div>
        )}

        {/* Editable File Name and Description */}
        <div className="mt-4">
          <form onSubmit={form.handleSubmit(onFormSubmit)} className="space-y-4">
            {/* File Name - Display Only */}
            <div>
              <Label>File Name</Label>
              <p className="text-sm text-muted-foreground italic mt-1">
                {item.fileName}
              </p>
            </div>

            {/* Description */}
            <div>
              <Label htmlFor="description">Description</Label>
              {isEditing ? (
                <Textarea
                  {...form.register("description")}
                  className="mt-1"
                  rows={3}
                  placeholder="Enter description (optional)"
                />
              ) : (
                <div className="flex items-start justify-between mt-1">
                  <p className="text-sm text-muted-foreground italic">
                    {form.watch("description") || item.description || "No description available"}
                  </p>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={handleEditToggle}
                    className="h-8 w-8"
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                </div>
              )}
              {form.formState.errors.description && (
                <p className="mt-1 text-sm text-destructive">{form.formState.errors.description.message}</p>
              )}
            </div>

            {/* Edit/Cancel Buttons */}
            {isEditing && (
              <div className="flex gap-2">
                <Button type="submit" disabled={isSubmitting} className="flex-1">
                  {isSubmitting ? 'Saving...' : 'Save'}
                </Button>
                <Button type="button" variant="outline" onClick={handleEditToggle} className="flex-1">
                  Cancel
                </Button>
              </div>
            )}
          </form>
        </div>

        <div className="mt-4">
          <Label>File Size</Label>
          <p className="text-sm text-muted-foreground mt-1">
            {item.fileSize ? (item.fileSize / 1024 / 1024).toFixed(1) : 'N/A'} MB
          </p>
        </div>
      </div>

      {/* Information */}
      <div className="mt-6">
        <h3 className="font-medium text-foreground">Information</h3>
        <Separator className="my-2" />
        <dl className="divide-y">
          <div className="flex justify-between py-3 text-sm">
            <dt className="text-muted-foreground">Type</dt>
            <dd className="font-medium capitalize">{item.mediaType}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm">
            <dt className="text-muted-foreground">Content Type</dt>
            <dd className="font-medium">{item.contentType}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm">
            <dt className="text-muted-foreground">Status</dt>
            <dd className="font-medium capitalize">{item.status}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm">
            <dt className="text-muted-foreground">Created</dt>
            <dd className="font-medium">{item.createdAt ? new Date(item.createdAt).toLocaleDateString() : 'N/A'}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm">
            <dt className="text-muted-foreground">Last modified</dt>
            <dd className="font-medium">{item.updatedAt ? new Date(item.updatedAt).toLocaleDateString() : 'N/A'}</dd>
          </div>
        </dl>
      </div>

      {/* Processing Activity */}
      <div className="mt-6">
        <h3 className="font-medium text-foreground">Processing Activity</h3>
        <Separator className="my-2" />
        <div className="text-sm text-muted-foreground space-y-4">
          {item.processingActivity?.map((activity, index) => (
            <div key={index + (activity.status || '') + (activity.createdAt?.toString() || '')} className="space-y-1">
              <p className="text-foreground italic">{activity.message}</p>
              <p className="text-xs">Status: {activity.status || 'N/A'}</p>
              <p className="text-xs">Created: {activity.createdAt ? new Date(activity.createdAt).toLocaleString() : 'N/A'}</p>
              <p className="text-xs">Updated: {activity.updatedAt ? new Date(activity.updatedAt).toLocaleString() : 'N/A'}</p>
              {index < (item.processingActivity?.length || 0) - 1 && <Separator className="mt-2" />}
            </div>
          ))}
        </div>
      </div>

      {/* Download and Delete buttons */}
      <div className="mt-6 flex gap-3">
        <Button onClick={onDownload} className="flex-1">
          Download
        </Button>
        <Button
          variant="destructive"
          onClick={() => setDeleteModalOpen(true)}
          className="flex-1"
        >
          Delete
        </Button>
      </div>

      {/* Delete Confirmation Modal */}
      <Dialog open={deleteModalOpen} onOpenChange={setDeleteModalOpen}>
        <DialogContent>
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>
            <DialogTitle className="text-center">Delete media</DialogTitle>
            <DialogDescription className="text-center">
              Are you sure you want to delete <strong>{item.fileName}</strong>? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="sm:grid sm:grid-cols-2 sm:gap-3">
            <Button variant="outline" onClick={() => setDeleteModalOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteConfirm}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
