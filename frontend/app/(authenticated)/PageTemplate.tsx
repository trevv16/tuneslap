"use client"

import { useThemeContext } from "@/contexts/ThemeContext";

type PageTemplateProps = {
  children: React.ReactNode;
}

export default function PageTemplate({ children }: PageTemplateProps) {
  const { theme } = useThemeContext();

  return (
    <main className="px-4 pb-12 sm:px-6 lg:px-8" data-theme={theme}>
      {children}
    </main>
  )
}