"use client";

import type { KeyResponse as Key, MediaListItem as Media } from "@/api/models";
import { useAllMedia } from "@/hooks/media";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const KeySchema = (existingHotKeys: string[], excludeHotKey?: string) =>
  z.object({
    name: z.string().min(3).max(100),
    description: z.string().max(500).optional(),
    hotKey: z
      .string()
      .min(1, "Hotkey is required")
      .max(1, "Hotkey must be a single character")
      .refine(
        (val) => {
          const upperVal = val.toUpperCase();
          return !existingHotKeys.some(hotKey =>
            hotKey.toUpperCase() === upperVal &&
            hotKey.toUpperCase() !== excludeHotKey?.toUpperCase()
          );
        },
        "This hotkey is already used on this board"
      ),
    audioMediaId: z.string().min(1, "Audio selection is required"),
    imageMediaId: z.string().optional(),
  });

type KeyFormData = z.infer<ReturnType<typeof KeySchema>>;

type KeyFormProps = {
  boardId: string;
  existingKeys: Key[];
  mode: 'add' | 'edit';
  initialData?: Key;
  onSubmit: (data: KeyFormData & { boardId: string }) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
};

export default function KeyForm({
  boardId,
  existingKeys,
  mode,
  initialData,
  onSubmit,
  onCancel,
  isSubmitting = false
}: KeyFormProps) {
  const { data: mediaData, isLoading } = useAllMedia();
  const [selectedAudio, setSelectedAudio] = useState<string | null>(initialData?.audioMediaId || null);
  const [selectedImage, setSelectedImage] = useState<string | null>(initialData?.imageMediaId || null);

  const existingHotKeys = useMemo(() =>
    existingKeys
      .filter((k) => k.hotKey != null)
      .map((k) => k.hotKey!.toUpperCase()),
    [existingKeys]
  );
  const excludeHotKey = mode === 'edit' && initialData?.hotKey ? initialData.hotKey : undefined;

  const form = useForm({
    resolver: zodResolver(KeySchema(existingHotKeys, excludeHotKey)),
    defaultValues: {
      name: initialData?.name || "",
      description: initialData?.description || "",
      hotKey: initialData?.hotKey || "",
      audioMediaId: initialData?.audioMediaId || "",
      imageMediaId: initialData?.imageMediaId || "",
    },
  });

  // Filter media - useAllMedia now returns array directly
  const mediaArray = mediaData || [];
  const audioMedia = useMemo(() => mediaArray.filter((m: Media) => m.mediaType === "audio"), [mediaArray]);
  const imageMedia = useMemo(() => mediaArray.filter((m: Media) => m.mediaType === "image"), [mediaArray]);

  // Handle selection
  const handleAudioSelect = (id: string) => {
    if (selectedAudio === id) {
      // Deselect if already selected
      setSelectedAudio(null);
      form.setValue("audioMediaId", "");
    } else {
      // Select new audio
      setSelectedAudio(id);
      form.setValue("audioMediaId", id);
    }
  };

  const handleImageSelect = (id: string) => {
    if (selectedImage === id) {
      // Deselect if already selected
      setSelectedImage(null);
      form.setValue("imageMediaId", "");
    } else {
      // Select new image
      setSelectedImage(id);
      form.setValue("imageMediaId", id);
    }
  };

  const handleFormSubmit = form.handleSubmit((data) => {
    onSubmit({ ...data, boardId });
  });

  return (
    <form onSubmit={handleFormSubmit}>
      <div className="my-12">
        <div>
          <label htmlFor="name" className="block text-sm/6 font-medium text-base">
            Name
          </label>
          <div className="mt-2">
            <input
              id="name"
              type="text"
              disabled={isSubmitting}
              className={`block w-full rounded-md bg-muted px-3 py-1.5 text-base border-1 border-highlight placeholder:text-highlight focus:border-accent sm:text-sm/6 ${form.formState.errors.name ? 'outline-red-500 focus:outline-red-500' : ''}`}
              placeholder="Enter key name"
              {...form.register("name")}
            />
          </div>
          {form.formState.errors.name && (
            <p className="mt-1 text-sm text-red-500">{form.formState.errors.name.message as string}</p>
          )}
        </div>

        <div className="mt-4">
          <label htmlFor="description" className="block text-sm/6 font-medium text-base">
            Description
          </label>
          <div className="mt-2">
            <textarea
              id="description"
              disabled={isSubmitting}
              className={`block w-full rounded-md bg-muted px-3 py-1.5 text-base border-1 border-highlight placeholder:text-highlight focus:border-accent sm:text-sm/6 ${form.formState.errors.description ? 'outline-red-500 focus:outline-red-500' : ''}`}
              placeholder="Enter key description"
              rows={3}
              {...form.register("description")}
            />
          </div>
          {form.formState.errors.description && (
            <p className="mt-1 text-sm text-red-500">{form.formState.errors.description.message as string}</p>
          )}
        </div>

        <div className="mt-4">
          <label htmlFor="hotKey" className="block text-sm/6 font-medium text-base">
            Hotkey
          </label>
          <div className="mt-2">
            <input
              id="hotKey"
              type="text"
              maxLength={1}
              disabled={isSubmitting}
              className={`block w-20 rounded-md bg-muted px-3 py-1.5 text-base border-1 border-highlight placeholder:text-highlight focus:border-accent sm:text-sm/6 ${form.formState.errors.hotKey ? 'outline-red-500 focus:outline-red-500' : ''}`}
              placeholder="A"
              {...form.register("hotKey")}
            />
          </div>
          {form.formState.errors.hotKey && (
            <p className="mt-1 text-sm text-red-500">{form.formState.errors.hotKey.message as string}</p>
          )}
        </div>

        <div className="mt-4">
          <label className="block text-sm/6 font-medium text-base">
            Audio (required)
          </label>
          <div className="mt-2">
            {isLoading ? (
              <p className="text-sm text-muted">Loading audio files...</p>
            ) : audioMedia.length === 0 ? (
              <p className="text-sm text-muted">No audio files found. Please upload some audio first.</p>
            ) : (
              <div className="grid grid-cols-1 gap-2 max-h-40 overflow-y-auto">
                {audioMedia.map((audio: Media) => (
                  <button
                    type="button"
                    key={audio.id || ''}
                    disabled={isSubmitting}
                    className={`text-left p-3 rounded-md border transition-colors ${selectedAudio === audio.id
                      ? 'border-primary-500 bg-primary-500/10'
                      : 'border-highlight hover:border-accent'
                      } ${isSubmitting ? 'opacity-50 cursor-not-allowed' : ''}`}
                    onClick={() => { if (audio.id) handleAudioSelect(audio.id); }}
                  >
                    <div className="font-medium text-sm text-highlight">{audio.fileName}</div>
                    <audio src={audio.fileUrl} controls className="w-full mt-2" />
                  </button>
                ))}
              </div>
            )}
          </div>
          {form.formState.errors.audioMediaId && (
            <p className="mt-1 text-sm text-red-500">{form.formState.errors.audioMediaId.message as string}</p>
          )}
        </div>

        <div className="mt-4">
          <label className="block text-sm/6 font-medium text-base">
            Image (optional)
          </label>
          <div className="mt-2">
            {isLoading ? (
              <p className="text-sm text-muted">Loading images...</p>
            ) : imageMedia.length === 0 ? (
              <p className="text-sm text-muted">No images found. Please upload some images first.</p>
            ) : (
              <div className="grid grid-cols-3 gap-2 max-h-40 overflow-y-auto">
                {imageMedia.map((img: Media) => (
                  <button
                    type="button"
                    key={img.id || ''}
                    disabled={isSubmitting}
                    className={`p-2 rounded-md border transition-colors ${selectedImage === img.id
                      ? 'border-primary-500 bg-primary-500/10'
                      : 'border-highlight hover:border-accent'
                      } ${isSubmitting ? 'opacity-50 cursor-not-allowed' : ''}`}
                    onClick={() => img.id && handleImageSelect(img.id)}
                  >
                    <img src={img.fileUrl} alt={img.fileName} className="w-full h-20 object-cover rounded" />
                    <div className="text-xs mt-1 truncate text-highlight">{img.fileName}</div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
        <button
          type="button"
          onClick={onCancel}
          disabled={isSubmitting}
          className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isSubmitting}
          className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isSubmitting ? (mode === 'add' ? 'Adding Key...' : 'Updating Key...') : (mode === 'add' ? 'Add Key' : 'Update Key')}
        </button>
      </div>
    </form>
  );
} 