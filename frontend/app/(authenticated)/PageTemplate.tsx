import { cn } from "@/lib/utils"

interface PageTemplateProps {
  children: React.ReactNode
  className?: string
}

export default function PageTemplate({ children, className }: PageTemplateProps) {
  return (
    <main className={cn("mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8", className)}>
      {children}
    </main>
  )
}
