import { BoardsApi } from "@/api";
import { getApiConfig } from "@/api/config";
import type { BoardResponse } from "@/api/models";
import {
  extractMediaUrls,
  preloadAudios,
  preloadImages,
} from "@/utils/preloadUtils";
import { getStoredToken } from "@/utils/token";
import { useQueryClient } from "@tanstack/react-query";
import { useCallback, useRef } from "react";
import { boardKeys } from "./boards";

interface UsePreloadBoardOptions {
  debounceMs?: number;
  preloadTimeout?: number;
}

export const usePreloadBoard = (options: UsePreloadBoardOptions = {}) => {
  const { debounceMs = 300, preloadTimeout = 10000 } = options;
  const queryClient = useQueryClient();
  const preloadingRef = useRef<Set<string>>(new Set());
  const timeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);

  const preloadBoardData = useCallback(
    async (boardId: string) => {
      // Prevent duplicate preloading
      if (preloadingRef.current.has(boardId)) {
        return;
      }

      preloadingRef.current.add(boardId);

      try {
        const authToken = getStoredToken();
        if (!authToken) return;

        // Prefetch board data using React Query
        await queryClient.prefetchQuery({
          queryKey: boardKeys.detail(boardId),
          queryFn: async () => {
            const boardsApi = new BoardsApi(getApiConfig());
            return await boardsApi.getBoardById({ boardId });
          },
          staleTime: 5 * 60 * 1000, // 5 minutes
          gcTime: 5 * 60 * 1000, // 5 minutes
        });

        // Get the board data to extract media URLs
        const board = queryClient.getQueryData<BoardResponse>(boardKeys.detail(boardId));
        if (!board) return;

        // Extract media URLs and preload them
        const { images, audios } = extractMediaUrls(board);

        // Preload media files concurrently
        await Promise.allSettled([
          preloadImages(images, { timeout: preloadTimeout }),
          preloadAudios(audios, { timeout: preloadTimeout }),
        ]);
      } catch (error) {
        console.warn(`Failed to preload board ${boardId}:`, error);
      } finally {
        preloadingRef.current.delete(boardId);
      }
    },
    [queryClient, preloadTimeout]
  );

  // Debounced version to prevent excessive calls
  const debouncedPreload = useCallback(
    (boardId: string) => {
      // Check if already preloading before setting timeout
      if (preloadingRef.current.has(boardId)) {
        return;
      }

      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      timeoutRef.current = setTimeout(() => {
        void preloadBoardData(boardId);
      }, debounceMs);
    },
    [preloadBoardData, debounceMs]
  );

  return {
    preloadBoard: debouncedPreload,
    isPreloading: (boardId: string) => preloadingRef.current.has(boardId),
  };
};
