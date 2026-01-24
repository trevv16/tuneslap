import { MediaApi } from "@/api";
import { GetAllMediaMediaTypeEnum } from "@/api/apis/MediaApi";
import { getApiConfig } from "@/api/config";
import type {
  CreateMediaRequest,
  CreateMediaResponse,
  DeleteMediaResponse,
  GetMyMediaStatsResponse,
  MediaListItem,
  MediaProcessingParams,
  MediaResponse,
  ProcessMediaResponse,
  UpdateMediaRequest,
  UpdateMediaResponse,
} from "@/api/models";
import {
  CreateMediaRequestContentTypeEnum,
  CreateMediaRequestMediaTypeEnum,
  MediaProcessingParamsAudioContentTypeEnum,
  MediaProcessingParamsAudioOutputFormatsEnum,
  MediaProcessingParamsImageApplyFiltersEnum,
  MediaProcessingParamsImageAspectRatioEnum,
  MediaProcessingParamsImageFormatEnum,
} from "@/api/models";
import { generateUploadUrl } from "@/api/uploadUrl";
import { getStoredToken } from "@/utils/token";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

// Query key factory for media
export const mediaKeys = {
  all: (params?: { page?: number; limit?: number; mediaType?: "image" | "audio" }) =>
    ["media", "list", params?.mediaType ?? "all", params?.page ?? 1, params?.limit ?? 25] as const,
  detail: (mediaId: string) => ["media", "detail", mediaId] as const,
  stats: () => ["media", "stats"] as const,
};

// Validation schema for media creation
const mediaCreateSchema = z.object({
  fileName: z
    .string()
    .min(3, "File name must be at least 3 characters")
    .max(100, "File name must be less than 100 characters"),
  description: z
    .string()
    .max(1000, "Description must be less than 1000 characters")
    .optional(),
  file: z.instanceof(File, { message: "Please select a file" }),
});

export type MediaCreateFormData = z.infer<typeof mediaCreateSchema>;

// Helper function to convert MIME type string to CreateMediaRequestContentTypeEnum
function convertMimeTypeToContentTypeEnum(mimeType: string): CreateMediaRequestContentTypeEnum | undefined {
  const mimeToEnum: Record<string, CreateMediaRequestContentTypeEnum> = {
    'audio/mp3': CreateMediaRequestContentTypeEnum.AudioMp3,
    'audio/mpeg': CreateMediaRequestContentTypeEnum.AudioMp3,
    'audio/wav': CreateMediaRequestContentTypeEnum.AudioWav,
    'audio/webm': CreateMediaRequestContentTypeEnum.AudioWebm,
    'audio/ogg': CreateMediaRequestContentTypeEnum.AudioOgg,
    'audio/aac': CreateMediaRequestContentTypeEnum.AudioAac,
    'image/jpeg': CreateMediaRequestContentTypeEnum.ImageJpeg,
    'image/png': CreateMediaRequestContentTypeEnum.ImagePng,
    'image/gif': CreateMediaRequestContentTypeEnum.ImageGif,
    'image/webp': CreateMediaRequestContentTypeEnum.ImageWebp,
    'image/svg+xml': CreateMediaRequestContentTypeEnum.ImageSvgxml,
  };
  return mimeToEnum[mimeType];
}

// Helper function to convert media type string to enum
function convertMediaTypeToEnum(mediaType: 'image' | 'audio'): CreateMediaRequestMediaTypeEnum {
  return mediaType === 'image' ? CreateMediaRequestMediaTypeEnum.Image : CreateMediaRequestMediaTypeEnum.Audio;
}

// Helper function to convert media type string to GetAllMediaMediaTypeEnum
function convertMediaTypeToGetAllEnum(mediaType: 'image' | 'audio' | undefined): GetAllMediaMediaTypeEnum | undefined {
  if (!mediaType) return undefined;
  return mediaType === 'image' ? GetAllMediaMediaTypeEnum.Image : GetAllMediaMediaTypeEnum.Audio;
}

// Custom hook for media creation
export function useMediaCreate() {
  const createMediaMutation = useCreateMedia();
  const { uploadFile, isUploading: isFileUploading } = useFileUpload();

  const form = useForm<MediaCreateFormData>({
    resolver: zodResolver(mediaCreateSchema),
    defaultValues: {
      fileName: "",
      description: "",
    },
  });

  // Get file extension from file type
  const getFileExtension = (file: File): string => {
    const mimeToExtension: Record<string, string> = {
      "audio/mp3": ".mp3",
      "audio/mpeg": ".mp3",
      "audio/wav": ".wav",
      "audio/webm": ".webm",
      "audio/ogg": ".ogg",
      "audio/aac": ".aac",
      "image/jpeg": ".jpg",
      "image/jpg": ".jpg",
      "image/png": ".png",
      "image/gif": ".gif",
      "image/webp": ".webp",
      "image/svg+xml": ".svg",
    };

    return mimeToExtension[file.type] || "";
  };

  const handleSubmit = async (data: MediaCreateFormData) => {
    try {
      // Get file extension and append it to fileName
      const fileExtension = getFileExtension(data.file);
      const fullFileName = data.fileName + fileExtension;

      // Step 1: Upload file to GCS
      const uploadResult = await uploadFile(data.file, fullFileName);
      if (!uploadResult.success || !uploadResult.fileUrl) {
        return { success: false, error: uploadResult.error };
      }

      // Step 2: Create media record
      const mediaType = data.file.type.startsWith("image/") ? "image" : "audio";

      const mediaRequest: CreateMediaRequest = {
        mediaType: convertMediaTypeToEnum(mediaType),
        fileName: fullFileName,
        description: data.description,
        fileUrl: uploadResult.fileUrl,
        contentType: convertMimeTypeToContentTypeEnum(data.file.type),
        fileSize: data.file.size,
        // TODO: set dimensions and duration after processing
      };

      const result = await createMediaMutation.mutateAsync(mediaRequest);

      if (result.success) {
        return { success: true, data: result.data };
      } else {
        return { success: false, error: result.data };
      }
    } catch (error) {
      return { success: false, error };
    }
  };

  return {
    form,
    handleSubmit,
    isSubmitting: createMediaMutation.isPending || isFileUploading,
    isSuccess: createMediaMutation.isSuccess,
    error: createMediaMutation.error,
  };
}

// Validation schema for media editing
const mediaEditSchema = z.object({
  description: z
    .string()
    .max(1000, "Description must be less than 1000 characters")
    .optional(),
});

export type MediaEditFormData = z.infer<typeof mediaEditSchema>;

// Custom hook for media editing
export function useMediaEdit(
  mediaId: string,
  initialData?: { fileName: string; description?: string }
) {
  const updateMediaMutation = useUpdateMedia();

  const form = useForm<MediaEditFormData>({
    resolver: zodResolver(mediaEditSchema),
    defaultValues: {
      description: initialData?.description || "",
    },
  });

  const handleSubmit = async (data: MediaEditFormData) => {
    try {
      await updateMediaMutation.mutateAsync({
        mediaId,
        data: {
          description: data.description,
        } as UpdateMediaRequest,
      });
      return { success: true };
    } catch (error) {
      return { success: false, error };
    }
  };

  return {
    form,
    handleSubmit,
    isSubmitting: updateMediaMutation.isPending,
    isSuccess: updateMediaMutation.isSuccess,
    error: updateMediaMutation.error,
  };
}

// Get media stats
export function useMediaStats() {
  const token = getStoredToken() || "";

  return useQuery<GetMyMediaStatsResponse, Error>({
    queryKey: mediaKeys.stats(),
    queryFn: async () => {
      const mediaApi = new MediaApi(getApiConfig());
      return await mediaApi.getMyMediaStats();
    },
    enabled: !!token && token !== "",
  });
}

// Get all media
export function useAllMedia(params?: {
  page?: number;
  limit?: number;
  mediaType?: "image" | "audio";
}) {
  const token = getStoredToken() || "";
  const queryKey = mediaKeys.all(params);

  return useQuery<MediaListItem[], Error>({
    queryKey,
    queryFn: async () => {
      console.log('[useAllMedia] Fetching with queryKey:', queryKey, 'mediaType:', params?.mediaType);
      const mediaApi = new MediaApi(getApiConfig());
      const response = await mediaApi.getAllMedia({
        page: params?.page,
        limit: params?.limit,
        mediaType: convertMediaTypeToGetAllEnum(params?.mediaType),
      });
      console.log('[useAllMedia] Got response with', response.data?.media?.length, 'items');
      return response.data?.media || [];
    },
    enabled: !!token && token !== "",
    staleTime: 30 * 1000, // 30 seconds - data is fresh for 30s before refetching
  });
}

// Get media by ID
export function useMediaById(mediaId: string) {
  const token = getStoredToken() || "";

  return useQuery<MediaResponse, Error>({
    queryKey: mediaKeys.detail(mediaId),
    queryFn: async () => {
      const mediaApi = new MediaApi(getApiConfig());
      return await mediaApi.getMediaById({ mediaId });
    },
    enabled: !!token && token !== "" && !!mediaId,
  });
}

// Create media
export function useCreateMedia() {
  const queryClient = useQueryClient();

  return useMutation<CreateMediaResponse, Error, CreateMediaRequest>({
    mutationFn: async (data: CreateMediaRequest) => {
      const mediaApi = new MediaApi(getApiConfig());
      return await mediaApi.createMedia({ createMediaRequest: data });
    },
    onSuccess: () => {
      // Invalidate and refetch media lists
      queryClient.invalidateQueries({ queryKey: mediaKeys.all() });
      queryClient.invalidateQueries({ queryKey: mediaKeys.stats() });
    },
  });
}

// Update media
export function useUpdateMedia() {
  const queryClient = useQueryClient();

  return useMutation<UpdateMediaResponse, Error, { mediaId: string; data: UpdateMediaRequest }>({
    mutationFn: async ({ mediaId, data }: { mediaId: string; data: UpdateMediaRequest }) => {
      const mediaApi = new MediaApi(getApiConfig());
      const response = await mediaApi.updateMedia({ mediaId, updateMediaRequest: data });
      return response;
    },
    onSuccess: (data, variables) => {
      // Update the specific media item in cache
      if (data.data) {
        queryClient.setQueryData(
          mediaKeys.detail(variables.mediaId),
          data.data
        );
      }
      // Invalidate and refetch media lists
      queryClient.invalidateQueries({ queryKey: mediaKeys.all() });
    },
  });
}

// Delete media
export function useDeleteMedia() {
  const queryClient = useQueryClient();

  return useMutation<DeleteMediaResponse, Error, string>({
    mutationFn: async (mediaId: string) => {
      const mediaApi = new MediaApi(getApiConfig());
      return await mediaApi.deleteMedia({ mediaId });
    },
    onSuccess: (data, mediaId) => {
      // Remove the deleted media from cache
      queryClient.removeQueries({ queryKey: mediaKeys.detail(mediaId) });
      // Invalidate and refetch media lists
      queryClient.invalidateQueries({ queryKey: mediaKeys.all() });
      queryClient.invalidateQueries({ queryKey: mediaKeys.stats() });
    },
  });
}

// Process media
export function useProcessMedia() {
  const queryClient = useQueryClient();

  return useMutation<ProcessMediaResponse, Error, { mediaId: string; processingParams: MediaProcessingParams }>({
    mutationFn: async ({
      mediaId,
      processingParams,
    }: {
      mediaId: string;
      processingParams: MediaProcessingParams;
    }) => {
      const mediaApi = new MediaApi(getApiConfig());
      return await mediaApi.processMedia({
        mediaId,
        mediaProcessingParams: processingParams,
      });
    },
    onSuccess: (data, variables) => {
      // Invalidate and refetch media lists
      queryClient.invalidateQueries({ queryKey: mediaKeys.all() });
      queryClient.invalidateQueries({ queryKey: mediaKeys.detail(variables.mediaId) });
    },
  });
}

// Hook for uploading files to GCS
export function useFileUpload() {
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);

  const uploadFile = async (
    file: File,
    fileName: string
  ): Promise<{ success: boolean; fileUrl?: string; error?: string }> => {
    setIsUploading(true);
    setUploadProgress(0);

    try {
      const token = getStoredToken();
      if (!token) {
        throw new Error("No authentication token available");
      }

      // Step 1: Generate signed upload URL
      const uploadUrlResult = await generateUploadUrl(token, {
        fileName,
        contentType: file.type,
        fileSize: file.size,
      });

      if (!uploadUrlResult.success || !uploadUrlResult.data) {
        throw new Error(
          uploadUrlResult.error || "Failed to generate upload URL"
        );
      }

      const { signedUrl, fileUrl } = uploadUrlResult.data;

      if (!fileUrl) {
        throw new Error("File URL not provided in upload URL response");
      }

      // Step 2: Upload file directly to storage
      const uploadResponse = await fetch(signedUrl, {
        method: "PUT",
        body: file,
        headers: {
          "Content-Type": file.type,
        },
        mode: "cors", // Explicitly set CORS mode
      });

      if (!uploadResponse.ok) {
        throw new Error(
          `Upload failed: ${uploadResponse.status} ${uploadResponse.statusText}`
        );
      }

      // Step 3: Use the file URL returned from the backend (no hardcoded URL construction)
      setUploadProgress(100);
      return { success: true, fileUrl };
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Upload failed";
      return { success: false, error: errorMessage };
    } finally {
      setIsUploading(false);
      setUploadProgress(0);
    }
  };

  return {
    uploadFile,
    uploadProgress,
    isUploading,
  };
}

// Validation schema for audio processing
const audioProcessingSchema = z.object({
  contentType: z.nativeEnum(MediaProcessingParamsAudioContentTypeEnum).optional(),
  trimStart: z.number().min(0, "Trim start must be non-negative").optional(),
  trimEnd: z.number().min(0, "Trim end must be non-negative").optional(),
  normalize: z.boolean().optional(),
  fadeIn: z.number().min(0, "Fade in must be non-negative").max(10, "Fade in must be at most 10 seconds").optional(),
  fadeOut: z.number().min(0, "Fade out must be non-negative").max(10, "Fade out must be at most 10 seconds").optional(),
  speed: z.number().min(0.5, "Speed must be at least 0.5x").max(2, "Speed must be at most 2x").optional(),
  pitch: z.number().min(-12, "Pitch must be at least -12 semitones").max(12, "Pitch must be at most +12 semitones").optional(),
  outputFormats: z.array(z.nativeEnum(MediaProcessingParamsAudioOutputFormatsEnum)).optional(),
});

export type AudioProcessingFormData = z.infer<typeof audioProcessingSchema>;

// Custom hook for audio processing form
export function useAudioProcessingForm(mediaId: string) {
  const processMediaMutation = useProcessMedia();

  const form = useForm<AudioProcessingFormData>({
    resolver: zodResolver(audioProcessingSchema),
    defaultValues: {
      contentType: MediaProcessingParamsAudioContentTypeEnum.AudioWebm,
      trimStart: 0,
      trimEnd: 0,
      normalize: true,
      fadeIn: 0,
      fadeOut: 0,
      speed: 1,
      pitch: 0,
    },
  });

  const handleSubmit = async (data: AudioProcessingFormData) => {
    try {
      const processingParams: MediaProcessingParams = {
        audio: {
          contentType: data.contentType,
          trimStart: data.trimStart,
          trimEnd: data.trimEnd,
          normalize: data.normalize,
          fadeIn: data.fadeIn,
          fadeOut: data.fadeOut,
          speed: data.speed,
          pitch: data.pitch,
          outputFormats: data.outputFormats ?? [MediaProcessingParamsAudioOutputFormatsEnum.Webm],
        },
      };

      await processMediaMutation.mutateAsync({
        mediaId,
        processingParams,
      });

      return { success: true };
    } catch (error) {
      return { success: false, error };
    }
  };

  return {
    form,
    handleSubmit,
    isSubmitting: processMediaMutation.isPending,
    isSuccess: processMediaMutation.isSuccess,
    error: processMediaMutation.error,
  };
}

// Validation schema for image processing
const imageProcessingSchema = z.object({
  format: z.nativeEnum(MediaProcessingParamsImageFormatEnum).optional(),
  aspectRatio: z.nativeEnum(MediaProcessingParamsImageAspectRatioEnum).optional(),
  resizeWidth: z.number().min(1).max(4096).optional(),
  resizeHeight: z.number().min(1).max(4096).optional(),
  applyFilters: z.array(z.nativeEnum(MediaProcessingParamsImageApplyFiltersEnum)).optional(),
});

export type ImageProcessingFormData = z.infer<typeof imageProcessingSchema>;

// Custom hook for image processing form
export function useImageProcessingForm(mediaId: string) {
  const processMediaMutation = useProcessMedia();

  const form = useForm<ImageProcessingFormData>({
    resolver: zodResolver(imageProcessingSchema),
    defaultValues: {
      format: MediaProcessingParamsImageFormatEnum.Webp,
      applyFilters: [],
    },
  });

  const handleSubmit = async (data: ImageProcessingFormData) => {
    try {
      const processingParams: MediaProcessingParams = {
        image: {
          format: data.format ?? MediaProcessingParamsImageFormatEnum.Webp,
          aspectRatio: data.aspectRatio,
          resizeTo: data.resizeWidth && data.resizeHeight
            ? [data.resizeWidth, data.resizeHeight]
            : undefined,
          applyFilters: data.applyFilters && data.applyFilters.length > 0 ? data.applyFilters : undefined,
        },
      };

      await processMediaMutation.mutateAsync({
        mediaId,
        processingParams,
      });

      return { success: true };
    } catch (error) {
      return { success: false, error };
    }
  };

  return {
    form,
    handleSubmit,
    isSubmitting: processMediaMutation.isPending,
    isSuccess: processMediaMutation.isSuccess,
    error: processMediaMutation.error,
  };
}
