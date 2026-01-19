'use client'

import type { MediaListItem as Media } from '@/api/models';
import { MediaEditFormData, useMediaEdit } from '@/hooks/media';
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle } from '@headlessui/react';
import { ExclamationTriangleIcon, PencilIcon, XMarkIcon } from '@heroicons/react/20/solid';
import { useEffect, useState } from 'react';
import { toast } from 'react-hot-toast';

type MediaDetailsProps = {
  item: Media;
  onClose: () => void;
  onDownload?: () => void;
  onDelete?: () => void;
}

export default function MediaDetails({
  item,
  onClose,
  onDownload,
  onDelete
}: MediaDetailsProps) {

  return (
    <>
      {/* Mobile overlay */}
      <div className="fixed inset-0 z-50 lg:hidden">
        <div className="absolute inset-0 bg-black/55" onClick={onClose} />
        <div className="absolute right-0 top-0 h-full w-full max-w-sm bg-elevated shadow-xl">
          <div className="flex h-full flex-col overflow-y-auto">
            <div className="flex items-center justify-between border-b border-highlight p-4">
              <h2 className="text-lg font-medium text-highlight">Details</h2>
              <button
                onClick={onClose}
                className="rounded-full p-1 text-base hover:bg-highlight"
              >
                <XMarkIcon className="h-6 w-6" />
              </button>
            </div>
            <div className="flex-1 p-4">
              <MediaDetailsContent
                item={item}
                onDownload={onDownload}
                onDelete={onDelete}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Desktop sidebar */}
      <aside className="hidden w-96 overflow-y-auto border-l border-highlight bg-elevated p-8 lg:block">
        <div className="space-y-6 pb-16">
          <MediaDetailsContent
            item={item}
            onDownload={onDownload}
            onDelete={onDelete}
          />
        </div>
      </aside>
    </>
  )
}

function MediaDetailsContent({
  item,
  onDownload,
  onDelete
}: {
  item: Media;
  onDownload?: () => void;
  onDelete?: () => void;
}) {
  const [isEditing, setIsEditing] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);

  const {
    form,
    handleSubmit,
    isSubmitting
  } = useMediaEdit(item.id ?? '', {
    fileName: item.fileName ?? '',
    description: item.description,
  });

  // Update form values when item changes
  useEffect(() => {
    form.reset({
      description: item.description,
    });
  }, [item.description, form]);

  const handleEditToggle = () => {
    if (isEditing) {
      // Cancel editing
      form.reset({
        description: item.description,
      });
      setIsEditing(false);
    } else {
      // Start editing
      setIsEditing(true);
    }
  };

  const onFormSubmit = async (data: MediaEditFormData) => {
    const result = await handleSubmit(data);
    if (result.success) {
      toast.success('Media updated successfully');
      setIsEditing(false);
      // Reset form with updated values
      form.reset({
        description: data.description,
      });
    } else {
      toast.error('Failed to update media');
    }
  };

  // Handle delete confirmation
  const handleDeleteConfirm = async () => {
    try {
      if (onDelete) {
        await onDelete();
        toast.success('Media deleted successfully!');
        setDeleteModalOpen(false);
      }
    } catch (error) {
      toast.error('Failed to delete media. Please try again.');
      console.error('Delete media error:', error);
    }
  };

  return (
    <>
      {/* Media Preview */}
      <div>
        {item.mediaType === 'image' ? (
          <img
            alt=""
            src={(item.fileUrl && item.fileUrl !== "") ? item.fileUrl : "/defaultKey.png"}
            className="block aspect-10/7 w-full rounded-lg object-cover"
          />
        ) : (
          <div className="block aspect-10/7 w-full rounded-lg bg-gray-200 flex items-center justify-center">
            <div className="text-gray-500 text-6xl">🎵</div>
          </div>
        )}

        {/* Editable File Name and Description */}
        <div className="mt-4">
          <form onSubmit={form.handleSubmit(onFormSubmit)} className="space-y-4">
            {/* File Name - Display Only */}
            <div>
              <label className="block text-sm font-medium text-highlight mb-1">
                File Name
              </label>
              <div className="flex items-center justify-between">
                <h2 className="text-sm text-base italic">
                  {item.fileName}
                </h2>
              </div>
            </div>

            {/* Description */}
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-highlight mb-1">
                Description
              </label>
              {isEditing ? (
                <textarea
                  {...form.register("description")}
                  rows={3}
                  className="w-full rounded-md bg-white px-3 py-2 text-sm text-gray-900 border border-gray-300 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  placeholder="Enter description (optional)"
                />
              ) : (
                <div className="flex items-start justify-between">
                  <p className="text-sm text-base italic">
                    {form.watch("description") || item.description || "No description available"}
                  </p>
                  <button
                    type="button"
                    onClick={handleEditToggle}
                    className="flex items-center justify-center rounded-full p-1 text-gray-400 hover:text-gray-500 hover:bg-gray-100"
                  >
                    <PencilIcon className="h-4 w-4" />
                  </button>
                </div>
              )}
              {form.formState.errors.description && (
                <p className="mt-1 text-sm text-red-600">{form.formState.errors.description.message}</p>
              )}
            </div>

            {/* Edit/Cancel Button */}
            {isEditing && (
              <div className="flex gap-2">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50"
                >
                  {isSubmitting ? 'Saving...' : 'Save'}
                </button>
                <button
                  type="button"
                  onClick={handleEditToggle}
                  className="flex-1 rounded-md bg-gray-300 px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm hover:bg-gray-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-600"
                >
                  Cancel
                </button>
              </div>
            )}
          </form>
        </div>

        <p className="font-sm text-highlight mt-4">File Size</p>
        <p className="text-sm font-medium text-base mt-2">{item.fileSize ? (item.fileSize / 1024 / 1024).toFixed(1) : 'N/A'} MB</p>
      </div>

      {/* Information */}
      <div>
        <h3 className="font-medium text-highlight">Information</h3>
        <dl className="mt-2 divide-y divide-gray-200/20 border-t border-b border-highlight">
          <div className="flex justify-between py-3 text-sm font-medium">
            <dt className="text-base">Type</dt>
            <dd className="whitespace-nowrap text-highlight capitalize">{item.mediaType}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm font-medium">
            <dt className="text-base">Content Type</dt>
            <dd className="whitespace-nowrap text-highlight">{item.contentType}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm font-medium">
            <dt className="text-base">Status</dt>
            <dd className="whitespace-nowrap text-highlight capitalize">{item.status}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm font-medium">
            <dt className="text-base">Created</dt>
            <dd className="whitespace-nowrap text-highlight">{item.createdAt ? new Date(item.createdAt).toLocaleDateString() : 'N/A'}</dd>
          </div>
          <div className="flex justify-between py-3 text-sm font-medium">
            <dt className="text-base">Last modified</dt>
            <dd className="whitespace-nowrap text-highlight">{item.updatedAt ? new Date(item.updatedAt).toLocaleDateString() : 'N/A'}</dd>
          </div>
        </dl>
      </div>

      {/* Description */}
      <div>
        <h3 className="font-medium text-highlight">Processing Activity</h3>
        <hr className="mt-2 mb-6 border-highlight" />
        <div className="mt-2 flex items-center justify-between">
          <div className="text-sm text-base italic w-full">
            {item.processingActivity?.map((activity, index) => (
              <div key={index + (activity.status || '') + (activity.createdAt?.toString() || '')}>
                <p className="mb-2 text-md text-highlight italic">{activity.message}</p>
                <p className="text-sm text-base">Status: {activity.status || 'N/A'}</p>
                <p className="text-sm text-base italic">Created: {activity.createdAt ? new Date(activity.createdAt).toLocaleString() : 'N/A'}</p>
                <p className="text-sm text-base italic">Updated: {activity.updatedAt ? new Date(activity.updatedAt).toLocaleString() : 'N/A'}</p>
                <hr className="my-2" />
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Download and Delete buttons */}
      <div className="flex gap-x-3">
        <button
          type="button"
          onClick={onDownload}
          className="flex-1 rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          Download
        </button>
        <button
          type="button"
          onClick={() => setDeleteModalOpen(true)}
          className="flex-1 rounded-md bg-error text-error px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-gray-50"
        >
          Delete
        </button>
      </div>

      {/* Delete Confirmation Modal */}
      <Dialog open={deleteModalOpen} onClose={() => setDeleteModalOpen(false)} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-sm sm:p-6 data-closed:sm:translate-y-0 data-closed:sm:scale-95"
            >
              <div>
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-error">
                  <ExclamationTriangleIcon aria-hidden="true" className="size-6 text-error" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-gray-900">
                    Delete media
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Are you sure you want to delete <strong>{item.fileName}</strong>? This action cannot be undone.
                    </p>
                  </div>
                </div>
                <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                  <button
                    type="button"
                    onClick={() => setDeleteModalOpen(false)}
                    className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    onClick={handleDeleteConfirm}
                    className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>
    </>
  )
} 