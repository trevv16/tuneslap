import { AuthApi } from "@/api";
import { getApiConfig } from "@/api/config";
import type {
  ForgotRequest,
  ForgotResponse,
  ResetRequest,
  ResetResponse,
  SigninRequest,
  SigninResponse,
  SignupRequest,
  SignupResponse,
} from "@/api/models";
import { setStoredToken } from "@/utils/token";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { userKeys } from "./users";

// Sign Up Hook
export const useSignUp = () => {
  return useMutation<SignupResponse, Error, SignupRequest>({
    mutationFn: async (request: SignupRequest) => {
      const authApi = new AuthApi(getApiConfig());
      return await authApi.signup({ signupRequest: request });
    },
    onSuccess: (data) => {
      // Handle successful signup
      console.log("Signup successful:", data.message);
    },
    onError: (error) => {
      // Handle signup error
      console.error("Signup failed:", error.message);
    },
  });
};

// Sign In Hook
export const useSignIn = () => {
  const queryClient = useQueryClient();

  return useMutation<SigninResponse, Error, SigninRequest>({
    mutationFn: async (request: SigninRequest) => {
      const authApi = new AuthApi(getApiConfig());
      return await authApi.signin({ signinRequest: request });
    },
    onSuccess: (data) => {
      // Save token on successful login
      if (data.success && data.data?.token) {
        setStoredToken(data.data.token);
        // Invalidate the user query to trigger a refetch with the new token
        queryClient.invalidateQueries({ queryKey: userKeys.me() });
      }
    },
    onError: (error) => {
      // Handle signin error
      console.error("Signin failed:", error.message);
    },
  });
};

// Forgot Password Hook
export const useForgotPassword = () => {
  return useMutation<ForgotResponse, Error, ForgotRequest>({
    mutationFn: async (request: ForgotRequest) => {
      const authApi = new AuthApi(getApiConfig());
      return await authApi.forgot({ forgotRequest: request });
    },
    onSuccess: (data) => {
      // Handle successful forgot password request
      console.log("Forgot password email sent:", data.message);
    },
    onError: (error) => {
      // Handle forgot password error
      console.error("Forgot password failed:", error.message);
    },
  });
};

// Reset Password Hook
export const useResetPassword = () => {
  return useMutation<ResetResponse, Error, ResetRequest>({
    mutationFn: async (request: ResetRequest) => {
      const authApi = new AuthApi(getApiConfig());
      return await authApi.reset({ resetRequest: request });
    },
    onSuccess: (data) => {
      // Handle successful password reset
      console.log("Password reset successful:", data.message);
    },
    onError: (error) => {
      // Handle password reset error
      console.error("Password reset failed:", error.message);
    },
  });
};
