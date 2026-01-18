"use client";

import type { KeyResponse as BoardKey } from "@/api/models";
import { useAudio } from "@/hooks/useAudio";
import { useEffect, useRef } from "react";

export default function SoundKey({ boardKey }: { boardKey: BoardKey }) {
  const buttonRef = useRef<HTMLButtonElement | null>(null);
  const { bindPressHandlers } = useAudio(boardKey.audioUrl || '', boardKey.hotKey || ''); // play on button press or "a" key

  useEffect(() => {
    bindPressHandlers(buttonRef.current);
  }, [bindPressHandlers]);

  return (
    <>
      <div className="group overflow-hidden rounded-lg focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 focus-within:ring-offset-gray-100">
        <img
          alt="Key Image"
          src={(boardKey.imageUrl && boardKey.imageUrl !== "") ? boardKey.imageUrl : "/defaultKey.png"}
          className="pointer-events-none aspect-10/7 object-cover group-hover:opacity-75"
        />
        <button type="button" ref={buttonRef} className="absolute inset-0 focus:outline-hidden">
          <span className="sr-only">View details for {boardKey.name}</span>
        </button>
      </div>
      <p className="pointer-events-none mt-2 block truncate text-sm font-medium text-base">{boardKey.name}</p>
      <p className="pointer-events-none block text-sm font-medium text-muted">{boardKey.hotKey}</p>
    </>
  )
}