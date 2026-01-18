'use client'

import { useAuthRedirect } from "../../hooks/useAuthRedirect";
import Navbar from "./Navbar";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  // Ensure user is authenticated for all pages in this route group
  useAuthRedirect();

  return (
    <>
      <Navbar />
      <div className="min-h-full bg-base">
        {children}
      </div>
    </>
  )
}
