'use client'

import DemoBanner from '@/components/DemoBanner';
import { MediaCreateFormData, useMediaCreate, useMediaStats } from '@/hooks/media';
import { formatBytes } from '@/utils/helpers';
import { CloudArrowUpIcon } from '@heroicons/react/20/solid';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'react-hot-toast';

type CreateMediaFormProps = {
  setOpen: (open: boolean) => void;
}

export default function CreateMediaForm({ setOpen }: CreateMediaFormProps) {
  const [dragActive, setDragActive] = useState(false);

  const {
    form,
    handleSubmit,
    isSubmitting
  } = useMediaCreate();

  const { data: mediaStats } = useMediaStats();

  const selectedFile = form.watch('file');
  const fileName = form.watch('fileName');

  // Check if file size exceeds available storage
  const storageError = useMemo(() => {
    if (!selectedFile || !mediaStats?.data) return null;

    const availableStorage = mediaStats.data?.availableStorage;
    // -1 means unlimited storage
    if (availableStorage === -1 || availableStorage === undefined) return null;

    if (selectedFile.size > availableStorage) {
      return `File size (${formatBytes(selectedFile.size)}) exceeds available storage (${formatBytes(availableStorage)})`;
    }

    return null;
  }, [selectedFile, mediaStats]);

  // Get file extension from selected file
  const getFileExtension = (file: File | undefined): string => {
    if (!file) return '';

    // Map MIME types to extensions
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
    };

    // Try MIME type first
    const mimeExtension = mimeToExtension[file.type];
    if (mimeExtension) {
      return mimeExtension;
    }

    // Fallback: extract extension from filename
    const fileName = file.name;
    const lastDotIndex = fileName.lastIndexOf('.');
    if (lastDotIndex !== -1) {
      return fileName.substring(lastDotIndex);
    }

    return '';
  };

  const fileExtension = getFileExtension(selectedFile);

  // Update fileName when file changes to include extension
  useEffect(() => {
    if (selectedFile && !fileName) {
      const baseName = selectedFile.name.replace(/\.[^/.]+$/, ""); // Remove existing extension
      form.setValue('fileName', baseName);
    }
  }, [selectedFile, fileName, form]);

  const handleFormSubmit = async (data: MediaCreateFormData) => {
    const result = await handleSubmit(data);
    if (result.success) {
      toast.success('Media created successfully!');
      setOpen(false);
      form.reset();
    } else {
      const errorMessage = result.error instanceof Error
        ? result.error.message
        : typeof result.error === 'string'
          ? result.error
          : 'Failed to create media';
      toast.error(errorMessage);
    }
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      form.setValue('file', file);

      // Auto-fill fileName if not already set
      if (!form.getValues('fileName')) {
        const fileName = file.name.replace(/\.[^/.]+$/, ""); // Remove extension
        form.setValue('fileName', fileName);
      }
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      form.setValue('file', file);

      // Auto-fill fileName if not already set
      if (!form.getValues('fileName')) {
        const fileName = file.name.replace(/\.[^/.]+$/, ""); // Remove extension
        form.setValue('fileName', fileName);
      }
    }
  };

  return (
    <form onSubmit={form.handleSubmit(handleFormSubmit)} className="space-y-4">
      <DemoBanner message="Demo mode: Files will be deleted within one hour. Max file size: 10MB. Max uploads: 5." />
      {/* File Upload */}
      <div>
        <label className="block text-sm/6 font-medium text-base mb-2">
          File
        </label>
        <div
          className={`relative border-2 border-dashed rounded-lg p-6 text-center ${dragActive
            ? 'border-accent bg-accent/10'
            : 'border-muted hover:border-base'
            }`}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          {selectedFile ? (
            <div className="space-y-2">
              <CloudArrowUpIcon className="mx-auto h-12 w-12 text-accent" />
              <p className="text-sm text-base font-medium">{selectedFile.name}</p>
              <p className="text-xs text-muted">
                {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
              </p>
              <button
                type="button"
                onClick={() => form.setValue('file', undefined as unknown as File)}
                className="text-sm text-error hover:text-error/80"
              >
                Remove file
              </button>
            </div>
          ) : (
            <div className="space-y-2">
              <CloudArrowUpIcon className="mx-auto h-12 w-12 text-muted" />
              <p className="text-sm text-muted">
                Drag and drop your file here, or{' '}
                <label className="text-accent hover:text-accent/80 cursor-pointer">
                  browse
                  <input
                    type="file"
                    className="hidden"
                    accept="image/*,audio/*"
                    onChange={handleFileSelect}
                  />
                </label>
              </p>
              <p className="text-xs text-muted">
                Supports images and audio files
              </p>
            </div>
          )}
        </div>
        {form.formState.errors.file && (
          <p className="mt-1 text-sm text-error">{form.formState.errors.file.message}</p>
        )}
        {storageError && (
          <p className="mt-1 text-sm text-error">{storageError}</p>
        )}
        {isSubmitting && (
          <div className="mt-2">
            <div className="flex items-center justify-between text-sm text-muted mb-1">
              <span>Uploading file...</span>
              <span>{Math.round((isSubmitting ? 50 : 0))}%</span>
            </div>
            <div className="w-full bg-muted rounded-full h-2">
              <div
                className="bg-accent h-2 rounded-full transition-all duration-300"
                style={{ width: `${isSubmitting ? 50 : 0}%` }}
              ></div>
            </div>
          </div>
        )}
      </div>

      {/* File Name with Extension Display */}
      <div>
        <label htmlFor="fileName" className="block text-sm/6 font-medium text-base">
          File Name
        </label>
        <div className="mt-2">
          <div className="relative">
            <input
              {...form.register("fileName")}
              type="text"
              className={`block w-full rounded-md bg-elevated h-8 px-3 py-1.5 text-base border-1 border-muted placeholder:text-muted focus:border-accent sm:text-sm/6 pr-16 ${form.formState.errors.fileName ? 'outline-error focus:outline-error' : ''}`}
              placeholder="Enter file name"
            />
            {fileExtension && (
              <div className="absolute inset-y-0 right-0 flex items-center pr-3">
                <span className="text-muted text-sm font-mono">
                  {fileExtension}
                </span>
              </div>
            )}
          </div>
          {form.formState.errors.fileName && (
            <p className="mt-1 text-sm text-error">{form.formState.errors.fileName.message}</p>
          )}
          {fileExtension && (
            <p className="mt-1 text-xs text-muted">
              Final filename will be: <span className="font-mono">{fileName || 'filename'}{fileExtension}</span>
            </p>
          )}
        </div>
      </div>

      {/* Description */}
      <div>
        <label htmlFor="description" className="block text-sm/6 font-medium text-base">
          Description (Optional)
        </label>
        <div className="mt-2">
          <textarea
            {...form.register("description")}
            rows={3}
            className={`block w-full rounded-md bg-elevated px-3 py-1.5 text-base border-1 border-muted placeholder:text-muted focus:border-accent sm:text-sm/6 ${form.formState.errors.description ? 'outline-error focus:outline-error' : ''}`}
            placeholder="Enter description"
          />
          {form.formState.errors.description && (
            <p className="mt-1 text-sm text-error">{form.formState.errors.description.message}</p>
          )}
        </div>
      </div>

      {/* Submit Button */}
      <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
        <button
          type="button"
          onClick={() => setOpen(false)}
          className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isSubmitting || !!storageError}
          className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isSubmitting ? 'Creating...' : 'Create Media'}
        </button>
      </div>
    </form>
  );
} 