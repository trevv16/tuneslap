'use client'

import { useMediaStats } from "@/hooks/media"
import { formatBytes } from "@/utils/helpers"
import { FileAudio, HardDrive, Image } from "lucide-react"

type HeaderProps = {
  pageTitle: string
  headerActions: React.ReactNode
}

export default function Header({ pageTitle, headerActions }: HeaderProps) {
  const { data: mediaStats } = useMediaStats()

  return (
    <header className="pb-6">
      <div className="lg:flex lg:items-center lg:justify-between">
        <div className="min-w-0 flex-1">
          <h2 className="text-2xl font-bold sm:truncate sm:text-3xl sm:tracking-tight text-foreground">
            {pageTitle}
          </h2>
          {mediaStats && (
            <div className="mt-1 flex flex-col sm:mt-2 sm:flex-row sm:flex-wrap sm:gap-x-6">
              <div className="mt-2 flex items-center text-sm text-muted-foreground">
                <FileAudio className="mr-1.5 h-5 w-5 shrink-0" />
                {mediaStats.data?.audioCount || 0} Audio
              </div>
              <div className="mt-2 flex items-center text-sm text-muted-foreground">
                <Image className="mr-1.5 h-5 w-5 shrink-0" />
                {mediaStats.data?.imageCount || 0} Images
              </div>
              <div className="mt-2 flex items-center text-sm text-muted-foreground">
                <HardDrive className="mr-1.5 h-5 w-5 shrink-0" />
                {formatBytes(mediaStats.data?.usedStorage || 0)} Used
              </div>
              <div className="mt-2 flex items-center text-sm text-muted-foreground">
                <HardDrive className="mr-1.5 h-5 w-5 shrink-0" />
                {formatBytes(mediaStats.data?.availableStorage || 0)} Available
              </div>
            </div>
          )}
        </div>
        <div className="mt-5 flex lg:mt-0 lg:ml-4">
          {headerActions}
        </div>
      </div>
    </header>
  )
}
