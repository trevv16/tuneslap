import { BoardsApi } from "@/api";
import { getApiConfig } from "@/api/config";
import type {
  CreateCollaboratorRequest,
  CreateCollaboratorResponse,
  DeleteCollaboratorResponse,
  GetAllCollaboratorsResponse,
  GetCollaboratorByIdResponse,
  UpdateCollaboratorRequest,
  UpdateCollaboratorResponse,
} from "@/api/models";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { boardKeys } from "./boards";

// Query key factory for collaborators
export const collaboratorKeys = {
  all: (boardId: string) => ["collaborators", boardId] as const,
  detail: (boardId: string, collaboratorId: string) => ["collaborator", boardId, collaboratorId] as const,
};

export const useGetCollaborators = (boardId: string) => {
  return useQuery<GetAllCollaboratorsResponse>({
    queryKey: collaboratorKeys.all(boardId),
    queryFn: async () => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.getAllCollaborators({ boardId });
    },
    enabled: !!boardId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useGetCollaboratorById = (boardId: string, collaboratorId: string) => {
  return useQuery<GetCollaboratorByIdResponse>({
    queryKey: collaboratorKeys.detail(boardId, collaboratorId),
    queryFn: async () => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.getCollaboratorById({ boardId, collaboratorId });
    },
    enabled: !!boardId && !!collaboratorId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCreateCollaborator = (boardId: string) => {
  const queryClient = useQueryClient();

  return useMutation<CreateCollaboratorResponse, Error, CreateCollaboratorRequest>({
    mutationFn: async (request: CreateCollaboratorRequest) => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.createCollaborator({ boardId, createCollaboratorRequest: request });
    },
    onSuccess: () => {
      // Invalidate collaborators for the specific board
      queryClient.invalidateQueries({
        queryKey: collaboratorKeys.all(boardId),
      });
      // Also invalidate the board data since collaborators are part of the board
      void queryClient.invalidateQueries({ queryKey: boardKeys.detail(boardId) });
    },
  });
};

export const useUpdateCollaborator = (boardId: string) => {
  const queryClient = useQueryClient();

  return useMutation<UpdateCollaboratorResponse, Error, { collaboratorId: string; data: UpdateCollaboratorRequest }>({
    mutationFn: async (request: {
      collaboratorId: string;
      data: UpdateCollaboratorRequest;
    }) => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.updateCollaborator({
        boardId,
        collaboratorId: request.collaboratorId,
        updateCollaboratorRequest: request.data,
      });
    },
    onSuccess: () => {
      // Invalidate collaborators for the specific board
      queryClient.invalidateQueries({
        queryKey: collaboratorKeys.all(boardId),
      });
      // Also invalidate the board data
      void queryClient.invalidateQueries({ queryKey: boardKeys.detail(boardId) });
    },
  });
};

export const useDeleteCollaborator = (boardId: string) => {
  const queryClient = useQueryClient();

  return useMutation<DeleteCollaboratorResponse, Error, string>({
    mutationFn: async (collaboratorId: string) => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.deleteCollaborator({ boardId, collaboratorId });
    },
    onSuccess: () => {
      // Invalidate collaborators for the specific board
      queryClient.invalidateQueries({
        queryKey: collaboratorKeys.all(boardId),
      });
      // Also invalidate the board data
      void queryClient.invalidateQueries({ queryKey: boardKeys.detail(boardId) });
    },
  });
};
