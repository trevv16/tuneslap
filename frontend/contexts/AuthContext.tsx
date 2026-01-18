import { getStoredToken, removeStoredToken } from "@/utils/token";
import type { UserResponse } from "@/api/models";
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
  const token = getStoredToken();
  const { data: userData, isLoading, error } = useGetMe(token || "");
  const [user, setUser] = useState<UserResponse | null>(null);

  const signOut = () => {
    removeStoredToken();
    setUser(null);
  }

  useEffect(() => {
    if (userData?.data) {
      setUser(userData.data);
    } else if (error || !token) {
      // Clear user if there's an error or no token
      setUser(null);
    }
  }, [userData, error, token]);

  // Consider loading if:
  // 1. We have a token and the query is loading, OR
  // 2. We have a token but user data is not set yet (waiting for useEffect to run)
  const isActuallyLoading = (isLoading && !!token) || (!!token && !user && !!userData);

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