import { BoardsApi } from "@/api";
import { getApiConfig } from "@/api/config";
import type {
  BoardResponse,
  CreateBoardRequest,
  CreateBoardResponse,
  UpdateBoardRequest,
  UpdateBoardResponse,
} from "@/api/models";
import { getStoredToken } from "@/utils/token";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

// Query key factory for boards
export const boardKeys = {
  all: () => ["boards"] as const,
  detail: (boardId: string) => ["board", boardId] as const,
};

export const useGetBoards = () => {
  const authToken = getStoredToken() || "";
  return useQuery<BoardResponse[]>({
    queryKey: boardKeys.all(),
    queryFn: async () => {
      const boardsApi = new BoardsApi(getApiConfig());
      const response = await boardsApi.getAllBoards({});
      return response.data?.boards || [];
    },
    enabled: !!authToken && authToken !== "",
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useGetBoardById = (boardId: string) => {
  const authToken = getStoredToken() || "";

  return useQuery<BoardResponse>({
    queryKey: boardKeys.detail(boardId),
    queryFn: async () => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.getBoardById({ boardId });
    },
    enabled: !!authToken && authToken !== "",
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCreateBoard = () => {
  const queryClient = useQueryClient();
  return useMutation<CreateBoardResponse, Error, CreateBoardRequest>({
    mutationFn: async (request: CreateBoardRequest) => {
      const boardsApi = new BoardsApi(getApiConfig());
      return await boardsApi.createBoard({ createBoardRequest: request });
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: boardKeys.all() });
    },
  });
};

export const useUpdateBoard = () => {
  const queryClient = useQueryClient();

  return useMutation<UpdateBoardResponse, Error, UpdateBoardRequest & { boardId: string }>({
    mutationFn: async (request: UpdateBoardRequest & { boardId: string }) => {
      const boardsApi = new BoardsApi(getApiConfig());
      const { boardId, ...updateRequest } = request;
      return await boardsApi.updateBoard({ boardId, updateBoardRequest: updateRequest });
    },
    onSuccess: (_data, variables) => {
      // Invalidate both the specific board and the boards list
      void queryClient.invalidateQueries({ queryKey: boardKeys.detail(variables.boardId) });
      void queryClient.invalidateQueries({ queryKey: boardKeys.all() });
    },
  });
};

export const useDeleteBoard = () => {
  const queryClient = useQueryClient();

  return useMutation<undefined, Error, string>({
    mutationFn: async (boardId: string) => {
      const boardsApi = new BoardsApi(getApiConfig());
      await boardsApi.deleteBoard({ boardId });
      return undefined;
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: boardKeys.all() });
    },
  });
};
