/**
 * Utility functions for preloading media files
 */

interface PreloadOptions {
  timeout?: number;
  onSuccess?: () => void;
  onError?: (error: Error) => void;
}

/**
 * Preload an image from a URL
 */
export const preloadImage = (
  url: string,
  options: PreloadOptions = {}
): Promise<void> => {
  return new Promise((resolve, reject) => {
    const { timeout = 10000, onSuccess, onError } = options;

    if (!url) {
      resolve();
      return;
    }

    const img = new Image();
    const timeoutId = setTimeout(() => {
      img.src = "";
      reject(new Error(`Image preload timeout: ${url}`));
    }, timeout);

    img.onload = () => {
      clearTimeout(timeoutId);
      onSuccess?.();
      resolve();
    };

    img.onerror = () => {
      clearTimeout(timeoutId);
      const error = new Error(`Failed to preload image: ${url}`);
      onError?.(error);
      reject(error);
    };

    img.src = url;
  });
};

/**
 * Preload an audio file from a URL
 */
export const preloadAudio = (
  url: string,
  options: PreloadOptions = {}
): Promise<void> => {
  return new Promise((resolve, reject) => {
    const { timeout = 10000, onSuccess, onError } = options;

    if (!url) {
      resolve();
      return;
    }

    const audio = new Audio();
    const timeoutId = setTimeout(() => {
      audio.src = "";
      reject(new Error(`Audio preload timeout: ${url}`));
    }, timeout);

    audio.oncanplaythrough = () => {
      clearTimeout(timeoutId);
      onSuccess?.();
      resolve();
    };

    audio.onerror = () => {
      clearTimeout(timeoutId);
      const error = new Error(`Failed to preload audio: ${url}`);
      onError?.(error);
      reject(error);
    };

    // Start loading the audio
    audio.src = url;
    audio.load();
  });
};

/**
 * Preload multiple images concurrently
 */
export const preloadImages = async (
  urls: string[],
  options: PreloadOptions = {}
): Promise<void> => {
  const validUrls = urls.filter((url) => url && url.trim() !== "");
  if (validUrls.length === 0) return;

  const promises = validUrls.map((url) =>
    preloadImage(url, options).catch((error: unknown) => {
      console.warn(`Failed to preload image: ${url}`, error);
      // Don't fail the entire batch for one failed image
    })
  );

  await Promise.allSettled(promises);
};

/**
 * Preload multiple audio files concurrently
 */
export const preloadAudios = async (
  urls: string[],
  options: PreloadOptions = {}
): Promise<void> => {
  const validUrls = urls.filter((url) => url && url.trim() !== "");
  if (validUrls.length === 0) return;

  const promises = validUrls.map((url) =>
    preloadAudio(url, options).catch((error: unknown) => {
      console.warn(`Failed to preload audio: ${url}`, error);
      // Don't fail the entire batch for one failed audio
    })
  );

  await Promise.allSettled(promises);
};

import type { BoardResponse as Board, KeyResponse as BoardKey } from "@/api/models";

/**
 * Extract all media URLs from a board's keys
 */
export const extractMediaUrls = (
  board: Board
): { images: string[]; audios: string[] } => {
  const images: string[] = [];
  const audios: string[] = [];

  // Add board image if it exists
  if (board.imageUrl && board.imageUrl.trim() !== "") {
    images.push(board.imageUrl);
  }

  // Extract media from board keys
  if (board.keys && Array.isArray(board.keys)) {
    board.keys.forEach((key: BoardKey) => {
      if (key.imageUrl && key.imageUrl.trim() !== "") {
        images.push(key.imageUrl);
      }
      if (key.audioUrl && key.audioUrl.trim() !== "") {
        audios.push(key.audioUrl);
      }
    });
  }

  return { images, audios };
};
