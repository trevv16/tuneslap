import type { UserResponse } from "@/api/models";
import { getStoredToken, removeStoredToken } from "@/utils/token";
import { useRouter } from "next/navigation";
import { createContext, Dispatch, SetStateAction, useContext, useEffect, useState } from "react";
import { useGetMe } from "../hooks/users";

type AuthContextType = {
  isLoading: boolean;
  isAuthenticated: boolean;
  user: UserResponse | null;
  setUser: Dispatch<SetStateAction<UserResponse | null>>;
  signOut: () => void;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const router = useRouter();
  const token = getStoredToken();
  const { data: userData, isLoading, error } = useGetMe(token || "");
  const [user, setUser] = useState<UserResponse | null>(null);

  const signOut = () => {
    removeStoredToken();
    setUser(null);
    router.push("/auth/signin");
  }

  useEffect(() => {
    if (userData?.data) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setUser(userData.data);
    } else if (error || !token) {
      // Clear user if there's an error or no token
      setUser(null);
    }
  }, [userData, error, token]);

  // Consider loading if we have a token but haven't resolved the user yet
  // (and there's no error that would indicate auth failure)
  const isActuallyLoading = !!token && !user && !error;

  // Consider authenticated if we have both token and user
  const isAuthenticated = !!token && !!user;

  return (
    <AuthContext.Provider
      value={{
        isLoading: isActuallyLoading,
        isAuthenticated,
        user,
        setUser,
        signOut
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export function useAuthContext() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('Must be a child of AuthContextProvider');
  }
  return context;
}