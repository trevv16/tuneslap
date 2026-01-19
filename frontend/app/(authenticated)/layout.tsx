'use client'

import { useAuthRedirect } from "../../hooks/useAuthRedirect"
import Navbar from "./Navbar"

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  useAuthRedirect()

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      {children}
    </div>
  )
}
