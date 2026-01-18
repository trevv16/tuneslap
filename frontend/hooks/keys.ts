import { KeysApi } from "@/api";
import { getApiConfig } from "@/api/config";
import type {
  CreateKeyRequest,
  CreateKeyResponse,
  DeleteKeyResponse,
  GetBoardKeysResponse,
  GetKeyByIdResponse,
  UpdateKeyRequest,
  UpdateKeyResponse,
} from "@/api/models";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { boardKeys } from "./boards";

// Query key factory for keys
export const keyKeys = {
  all: (boardId: string) => ["board", "keys", boardId] as const,
  detail: (boardId: string, keyId: string) => ["board", "key", boardId, keyId] as const,
};

// Create a key for a board
export function useCreateKey(boardId: string) {
  const queryClient = useQueryClient();

  return useMutation<CreateKeyResponse, Error, CreateKeyRequest>({
    mutationFn: async (data: CreateKeyRequest) => {
      const keysApi = new KeysApi(getApiConfig());
      return await keysApi.createKey({ boardId, createKeyRequest: data });
    },
    onSuccess: () => {
      // Invalidate and refetch board data and board keys
      queryClient.invalidateQueries({ queryKey: boardKeys.detail(boardId) });
      queryClient.invalidateQueries({ queryKey: keyKeys.all(boardId) });
      queryClient.invalidateQueries({ queryKey: boardKeys.all() });
    },
  });
}

// Get all keys for a board
export function useGetBoardKeys(boardId: string) {
  return useQuery<GetBoardKeysResponse, Error>({
    queryKey: keyKeys.all(boardId),
    queryFn: async () => {
      const keysApi = new KeysApi(getApiConfig());
      return await keysApi.getBoardKeys({ boardId });
    },
    enabled: !!boardId,
  });
}

// Get a key by ID
export function useGetKeyById(boardId: string, keyId: string) {
  return useQuery<GetKeyByIdResponse, Error>({
    queryKey: keyKeys.detail(boardId, keyId),
    queryFn: async () => {
      const keysApi = new KeysApi(getApiConfig());
      return await keysApi.getKeyById({ boardId, keyId });
    },
    enabled: !!boardId && !!keyId,
  });
}

// Update a key
export function useUpdateKey(boardId: string) {
  const queryClient = useQueryClient();

  return useMutation<UpdateKeyResponse, Error, { keyId: string; data: UpdateKeyRequest }>({
    mutationFn: async ({ keyId, data }: { keyId: string; data: UpdateKeyRequest }) => {
      const keysApi = new KeysApi(getApiConfig());
      return await keysApi.updateKey({ boardId, keyId, updateKeyRequest: data });
    },
    onSuccess: () => {
      // Invalidate and refetch board data
      queryClient.invalidateQueries({ queryKey: boardKeys.detail(boardId) });
      queryClient.invalidateQueries({ queryKey: keyKeys.all(boardId) });
      queryClient.invalidateQueries({ queryKey: boardKeys.all() });
    },
  });
}

// Delete a key
export function useDeleteKey(boardId: string) {
  const queryClient = useQueryClient();

  return useMutation<DeleteKeyResponse, Error, string>({
    mutationFn: async (keyId: string) => {
      const keysApi = new KeysApi(getApiConfig());
      return await keysApi.deleteKey({ boardId, keyId });
    },
    onSuccess: () => {
      // Invalidate and refetch board data
      queryClient.invalidateQueries({ queryKey: boardKeys.detail(boardId) });
      queryClient.invalidateQueries({ queryKey: keyKeys.all(boardId) });
      queryClient.invalidateQueries({ queryKey: boardKeys.all() });
    },
  });
}
