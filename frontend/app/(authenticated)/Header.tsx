'use client'

import { useMediaStats } from "@/hooks/media";
import { formatBytes } from "@/utils/helpers";
import { DocumentChartBarIcon, PhotoIcon, SpeakerWaveIcon } from "@heroicons/react/20/solid";

type HeaderProps = {
  pageTitle: string;
  headerActions: React.ReactNode;
}

export default function Header({ pageTitle, headerActions }: HeaderProps) {
  const { data: mediaStats } = useMediaStats();

  return (
    <div className="bg-dark">
      <header className="py-10">
        <div className="lg:flex lg:items-center lg:justify-between">
          <div className="min-w-0 flex-1">
            <h2 className="text-2xl/7 font-bold sm:truncate sm:text-3xl sm:tracking-tight text-base">
              {pageTitle}
            </h2>
            {mediaStats && (
              <div className="mt-1 flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap sm:space-x-6">
                <div className="mt-2 flex items-center text-sm text-highlight">
                  <SpeakerWaveIcon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
                  {mediaStats.data?.audioCount || 0} Audio
                </div>
                <div className="mt-2 flex items-center text-sm text-highlight">
                  <PhotoIcon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
                  {mediaStats.data?.imageCount || 0} Images
                </div>
                <div className="mt-2 flex items-center text-sm text-highlight">
                  <DocumentChartBarIcon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
                  {formatBytes(mediaStats.data?.usedStorage || 0)} Used
                </div>
                <div className="mt-2 flex items-center text-sm text-highlight">
                  <DocumentChartBarIcon aria-hidden="true" className="mr-1.5 size-5 shrink-0 text-highlight" />
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
    </div>
  )
}