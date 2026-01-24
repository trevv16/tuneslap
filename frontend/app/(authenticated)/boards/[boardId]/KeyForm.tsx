"use client";

import type { KeyResponse as Key, MediaListItem as Media } from "@/api/models";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Skeleton } from "@/components/ui/skeleton";
import { Textarea } from "@/components/ui/textarea";
import { useAllMedia } from "@/hooks/media";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { Check, Image as ImageIcon, Music } from "lucide-react";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const KeySchema = (existingHotKeys: string[], excludeHotKey?: string) =>
  z.object({
    name: z.string().min(3, "Name must be at least 3 characters").max(100),
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
        "This hotkey is already used"
      ),
    audioMediaId: z.string().min(1, "Audio is required"),
    imageMediaId: z.string().optional(),
  });

type KeyFormData = z.infer<ReturnType<typeof KeySchema>>;

interface KeyFormProps {
  boardId: string;
  existingKeys: Key[];
  mode: 'add' | 'edit';
  initialData?: Key;
  onSubmit: (data: KeyFormData & { boardId: string }) => void | Promise<void>;
  onCancel: () => void;
  isSubmitting?: boolean;
}

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
    existingKeys.flatMap((k) => k.hotKey ? [k.hotKey.toUpperCase()] : []),
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

  const audioMedia = useMemo(() => (mediaData || []).filter((m: Media) => m.mediaType === "audio"), [mediaData]);
  const imageMedia = useMemo(() => (mediaData || []).filter((m: Media) => m.mediaType === "image"), [mediaData]);

  const handleAudioSelect = (id: string) => {
    if (selectedAudio === id) {
      setSelectedAudio(null);
      form.setValue("audioMediaId", "");
    } else {
      setSelectedAudio(id);
      form.setValue("audioMediaId", id);
    }
  };

  const handleImageSelect = (id: string) => {
    if (selectedImage === id) {
      setSelectedImage(null);
      form.setValue("imageMediaId", "");
    } else {
      setSelectedImage(id);
      form.setValue("imageMediaId", id);
    }
  };

  const handleFormSubmit = form.handleSubmit((data) => {
    void onSubmit({ ...data, boardId });
  });

  return (
    <form onSubmit={handleFormSubmit} className="space-y-6">
      {/* Basic Info */}
      <div className="grid gap-4 sm:grid-cols-4">
        <div className="sm:col-span-3">
          <Label htmlFor="name">Name</Label>
          <Input
            id="name"
            disabled={isSubmitting}
            placeholder="Key name"
            className="mt-1.5"
            {...form.register("name")}
          />
          {form.formState.errors.name && (
            <p className="mt-1 text-sm text-destructive">{form.formState.errors.name.message}</p>
          )}
        </div>

        <div>
          <Label htmlFor="hotKey">Hotkey</Label>
          <Input
            id="hotKey"
            maxLength={1}
            disabled={isSubmitting}
            placeholder="A"
            className="mt-1.5 text-center font-mono uppercase"
            {...form.register("hotKey")}
          />
          {form.formState.errors.hotKey && (
            <p className="mt-1 text-sm text-destructive">{form.formState.errors.hotKey.message}</p>
          )}
        </div>
      </div>

      <div>
        <Label htmlFor="description">Description (optional)</Label>
        <Textarea
          id="description"
          disabled={isSubmitting}
          placeholder="Brief description"
          rows={2}
          className="mt-1.5"
          {...form.register("description")}
        />
        {form.formState.errors.description && (
          <p className="mt-1 text-sm text-destructive">{form.formState.errors.description.message}</p>
        )}
      </div>

      {/* Audio Selection */}
      <div>
        <Label className="flex items-center gap-2">
          <Music className="h-4 w-4" />
          Audio
        </Label>
        <div className="mt-2">
          {isLoading ? (
            <div className="space-y-2">
              <Skeleton className="h-16 w-full" />
              <Skeleton className="h-16 w-full" />
            </div>
          ) : audioMedia.length === 0 ? (
            <p className="text-sm text-muted-foreground py-4 text-center bg-muted rounded-lg">
              No audio files. Upload some in your library first.
            </p>
          ) : (
            <div className="grid gap-2 max-h-48 overflow-y-auto pr-1">
              {audioMedia.map((audio: Media) => (
                <button
                  type="button"
                  key={audio.id || ''}
                  disabled={isSubmitting}
                  className={cn(
                    "text-left p-3 rounded-lg border transition-colors relative",
                    selectedAudio === audio.id
                      ? "border-primary bg-primary/5"
                      : "border-border hover:border-muted-foreground/50",
                    isSubmitting && "opacity-50 cursor-not-allowed"
                  )}
                  onClick={() => { if (audio.id) handleAudioSelect(audio.id); }}
                >
                  {selectedAudio === audio.id && (
                    <div className="absolute top-2 right-2">
                      <Check className="h-4 w-4 text-primary" />
                    </div>
                  )}
                  <p className="text-sm font-medium truncate pr-6">{audio.fileName}</p>
                  <audio src={audio.fileUrl} controls className="w-full mt-2 h-8" />
                </button>
              ))}
            </div>
          )}
        </div>
        {form.formState.errors.audioMediaId && (
          <p className="mt-1 text-sm text-destructive">{form.formState.errors.audioMediaId.message}</p>
        )}
      </div>

      {/* Image Selection */}
      <div>
        <Label className="flex items-center gap-2">
          <ImageIcon className="h-4 w-4" />
          Image (optional)
        </Label>
        <div className="mt-2">
          {isLoading ? (
            <div className="grid grid-cols-4 gap-2">
              <Skeleton className="aspect-square" />
              <Skeleton className="aspect-square" />
              <Skeleton className="aspect-square" />
              <Skeleton className="aspect-square" />
            </div>
          ) : imageMedia.length === 0 ? (
            <p className="text-sm text-muted-foreground py-4 text-center bg-muted rounded-lg">
              No images. Upload some in your library first.
            </p>
          ) : (
            <div className="grid grid-cols-4 gap-2 max-h-48 overflow-y-auto pr-1">
              {imageMedia.map((img: Media) => (
                <button
                  type="button"
                  key={img.id || ''}
                  disabled={isSubmitting}
                  className={cn(
                    "relative rounded-lg border overflow-hidden aspect-square transition-colors",
                    selectedImage === img.id
                      ? "border-primary ring-2 ring-primary"
                      : "border-border hover:border-muted-foreground/50",
                    isSubmitting && "opacity-50 cursor-not-allowed"
                  )}
                  onClick={() => { if (img.id) handleImageSelect(img.id); }}
                >
                  <img src={img.fileUrl} alt={img.fileName} className="w-full h-full object-cover" />
                  {selectedImage === img.id && (
                    <div className="absolute inset-0 bg-primary/20 flex items-center justify-center">
                      <Check className="h-5 w-5 text-primary" />
                    </div>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="flex gap-3 pt-4 border-t">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting}
          className="flex-1"
        >
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={isSubmitting}
          className="flex-1"
        >
          {isSubmitting
            ? (mode === 'add' ? 'Adding...' : 'Saving...')
            : (mode === 'add' ? 'Add Key' : 'Save Changes')
          }
        </Button>
      </div>
    </form>
  );
}
