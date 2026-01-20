"use client";

import type { KeyResponse as BoardKey } from "@/api/models";
import { useAudio } from "@/hooks/useAudio";
import { useEffect, useRef } from "react";

export default function SoundKey({ boardKey }: { boardKey: BoardKey }) {
  const buttonRef = useRef<HTMLButtonElement | null>(null);
  const { bindPressHandlers } = useAudio(boardKey.audioUrl || '', boardKey.hotKey || '');

  useEffect(() => {
    bindPressHandlers(buttonRef.current);
  }, [bindPressHandlers]);

  return (
    <div className="group overflow-hidden rounded-lg focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 focus-within:ring-offset-background relative">
      <img
        alt="Key Image"
        src={(boardKey.imageUrl && boardKey.imageUrl !== "") ? boardKey.imageUrl : "/defaultKey.png"}
        className="pointer-events-none aspect-square object-cover group-hover:opacity-75 w-full"
      />
      
      {/* Overlay with name and hotkey */}
      <div className="absolute inset-0 flex flex-col justify-end bg-gradient-to-t from-black/70 to-transparent p-3">
        <p className="text-sm font-medium text-white truncate">{boardKey.name}</p>
        <p className="text-xs text-white/80 font-mono uppercase">{boardKey.hotKey}</p>
      </div>
      
      <button type="button" ref={buttonRef} className="absolute inset-0 focus:outline-none">
        <span className="sr-only">Play {boardKey.name}</span>
      </button>
    </div>
  )
}
