"use client";

import type { BoardResponse as Board } from '@/api/models';
import { usePreloadBoard } from '@/hooks/usePreloadBoard';
import Link from 'next/link';

type BoardsListProps = {
  boards: Board[];
}

export default function BoardsList({ boards }: BoardsListProps) {
  const { preloadBoard, isPreloading } = usePreloadBoard();

  return (
    <ul className="grid grid-cols-1 gap-x-6 gap-y-8 lg:grid-cols-3 xl:gap-x-8">
      {boards.map((board) => (
        <Link
          key={board.id}
          href={`/boards/${board.id}`}
          onMouseEnter={() => board.id && preloadBoard(board.id)}
        >
          <li
            className="overflow-hidden rounded-xl"
          >
            <div className="flex items-center gap-x-4 border-b border-gray-900/5 bg-dark-300 bottom-1 p-6">
              <img
                alt={board.name}
                src={(board.imageUrl && board.imageUrl !== "") ? board.imageUrl : "/defaultBoard.png"}
                className="size-12 flex-none rounded-lg bg-white object-cover ring-1 ring-gray-900/10"
              />
              <div className="text-sm/6 font-medium text-white">{board.name}</div>
              {board.id && isPreloading(board.id) && (
                <div className="absolute top-2 right-2">
                  <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
                </div>
              )}
            </div>
            <dl className="divide-y bg-dark-200 divide-gray-100 px-6 pb-4 text-sm/6">
              <div className="flex justify-between gap-x-4 py-3 text-base">
                {board.description}
              </div>
            </dl>
          </li>
        </Link>
      ))}
    </ul>
  )
}
