import { UsersApi } from "@/api";
import { getApiConfig } from "@/api/config";
import type { GetMeResponse, UpdateMeRequest, UpdateMeResponse } from "@/api/models";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

// Query key factory for users
export const userKeys = {
  me: () => ["me"] as const,
};

export const useGetMe = (authToken: string) => {
  return useQuery<GetMeResponse, Error>({
    queryKey: userKeys.me(),
    queryFn: async () => {
      const usersApi = new UsersApi(getApiConfig());
      return await usersApi.getMe();
    },
    enabled: !!authToken && authToken !== "",
    staleTime: 10 * 60 * 1000, // 10 minutes
    gcTime: 15 * 60 * 1000, // 15 minutes
  });
};

export const useUpdateMe = () => {
  const queryClient = useQueryClient();

  return useMutation<UpdateMeResponse, Error, UpdateMeRequest>({
    mutationFn: async (request: UpdateMeRequest) => {
      const usersApi = new UsersApi(getApiConfig());
      const response = await usersApi.updateMe({ updateMeRequest: request });
      return response;
    },
    onSuccess: () => {
      // Invalidate the user data query to refetch updated user info
      queryClient.invalidateQueries({ queryKey: userKeys.me() });
    },
  });
};
