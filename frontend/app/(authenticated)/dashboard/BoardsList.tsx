"use client"

import type { BoardResponse as Board } from '@/api/models'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { usePreloadBoard } from '@/hooks/usePreloadBoard'
import { Keyboard, LayoutGrid, Users } from 'lucide-react'
import Link from 'next/link'

type BoardsListProps = {
  boards: Board[]
}

export default function BoardsList({ boards }: BoardsListProps) {
  const { preloadBoard, isPreloading } = usePreloadBoard()

  const formatDate = (date: Date | undefined) => {
    if (!date) return 'Unknown'
    return new Date(date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  }


  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
      {boards.map((board) => (
        <Link
          key={board.id}
          href={`/boards/${board.id}`}
          onMouseEnter={() => board.id && preloadBoard(board.id)}
          className="group"
        >
          <Card className="h-full transition-all hover:shadow-lg hover:border-primary/50 relative overflow-hidden">
            {board.id && isPreloading(board.id) && (
              <div className="absolute top-2 right-2 z-10">
                <div className="w-2 h-2 bg-primary rounded-full animate-pulse" />
              </div>
            )}

            <CardHeader className="p-3 pb-2">
              <div className="aspect-square w-full rounded-md bg-muted flex items-center justify-center overflow-hidden mb-2">
                {board.imageUrl && board.imageUrl !== "" ? (
                  <img
                    alt={board.name}
                    src={board.imageUrl}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <LayoutGrid className="h-8 w-8 text-muted-foreground" />
                )}
              </div>
              <CardTitle className="text-sm font-semibold line-clamp-2 group-hover:text-primary transition-colors">
                {board.name}
              </CardTitle>
              <CardDescription className="text-xs">
                {formatDate(board.createdAt)}
              </CardDescription>
            </CardHeader>

            <CardContent className="p-3 pt-0">
              <p className="text-xs text-muted-foreground line-clamp-2 min-h-[2rem] mb-2">
                {board.description || "No description"}
              </p>

              <div className="flex items-center justify-between text-xs text-muted-foreground">
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-0.5">
                    <Keyboard className="h-3 w-3" />
                    <span>{board.keys?.length || 0}</span>
                  </div>
                  <div className="flex items-center gap-0.5">
                    <Users className="h-3 w-3" />
                    <span>{board.collaborators?.length || 0}</span>
                  </div>
                </div>
                {board.layout && (
                  <Badge variant="secondary" className="text-[10px] px-1 py-0 h-4">
                    {board.layout.toUpperCase()}
                  </Badge>
                )}
              </div>
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  )
}
