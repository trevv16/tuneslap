import { useAuthContext } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export function useAuthRedirect(redirectTo = "/auth/signin") {
  const { isAuthenticated, isLoading } = useAuthContext();
  const router = useRouter();

  useEffect(() => {
    // Only redirect if we're not loading and not authenticated
    // This ensures we wait for the auth state to be fully determined
    if (!isLoading && !isAuthenticated) {
      router.push(redirectTo);
    }
  }, [isAuthenticated, isLoading, redirectTo, router]);
}

export function usePublicPageRedirect(redirectTo = "/dashboard") {
  const { isAuthenticated, isLoading, user } = useAuthContext();
  const router = useRouter();

  useEffect(() => {
    // Only redirect if we're not loading and authenticated
    // This ensures we wait for the auth state to be fully determined
    if (!isLoading && isAuthenticated && user) {
      router.push(redirectTo);
    }
  }, [isAuthenticated, isLoading, redirectTo, router, user]);
}
