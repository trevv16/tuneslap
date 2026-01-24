"use client"

import type { BoardResponse as Board } from '@/api/models'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { useGetBoardById } from "@/hooks/boards"
import { Keyboard, LayoutGrid, Pencil, Users } from 'lucide-react'
import Link from 'next/link'
import Collaborators from "./Collaborators"
import DeleteBoardSection from "./DeleteBoardSection"
import EditBoardForm from "./EditBoardForm"
import Keys from "./Keys"

interface EditBoardClientProps {
  boardId: string
}

function BoardHeader({ board, isLoading }: { board: Board | undefined, isLoading: boolean }) {
  if (isLoading) {
    return (
      <div className="flex items-start gap-6 mb-8">
        <Skeleton className="h-20 w-20 rounded-lg" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-4 w-96" />
          <div className="flex gap-4 pt-2">
            <Skeleton className="h-5 w-20" />
            <Skeleton className="h-5 w-20" />
            <Skeleton className="h-5 w-16" />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex items-start gap-6 mb-8">
      <div className="h-20 w-20 shrink-0 rounded-lg bg-muted flex items-center justify-center overflow-hidden">
        {board?.imageUrl && board.imageUrl !== "" ? (
          <img
            alt={board.name}
            src={board.imageUrl}
            className="h-full w-full object-cover"
          />
        ) : (
          <LayoutGrid className="h-10 w-10 text-muted-foreground" />
        )}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-foreground truncate">{board?.name || 'Board'}</h1>
          <Button variant="outline" size="sm" asChild>
            <Link href={`/boards/${board?.id}`}>
              View Board
            </Link>
          </Button>
        </div>
        <p className="text-muted-foreground mt-1 line-clamp-2">{board?.description || 'No description'}</p>
        <div className="flex items-center gap-4 mt-3 text-sm text-muted-foreground">
          <div className="flex items-center gap-1.5">
            <Keyboard className="h-4 w-4" />
            <span>{board?.keys?.length || 0} keys</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Users className="h-4 w-4" />
            <span>{board?.collaborators?.length || 0} collaborators</span>
          </div>
          {board?.layout && (
            <Badge variant="secondary">{board.layout.toUpperCase()}</Badge>
          )}
        </div>
      </div>
    </div>
  )
}

export default function EditBoardClient({ boardId }: EditBoardClientProps) {
  const boardQuery = useGetBoardById(boardId)
  const { data: board, isLoading } = boardQuery

  return (
    <div className="space-y-6">
      {/* Header with board overview */}
      <BoardHeader board={board} isLoading={isLoading} />
      
      <Separator />

      {/* Board Details Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Pencil className="h-5 w-5 text-muted-foreground" />
            <div>
              <CardTitle>Board Details</CardTitle>
              <CardDescription>Update your board name, description, and settings.</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <EditBoardForm boardQuery={boardQuery} boardId={boardId} />
        </CardContent>
      </Card>

      {/* Collaborators Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Users className="h-5 w-5 text-muted-foreground" />
            <div>
              <CardTitle>Collaborators</CardTitle>
              <CardDescription>Manage who has access to this board.</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Collaborators boardId={boardId} />
        </CardContent>
      </Card>

      {/* Keys Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Keyboard className="h-5 w-5 text-muted-foreground" />
            <div>
              <CardTitle>Keys</CardTitle>
              <CardDescription>Configure the sound keys for this board.</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Keys boardId={boardId} />
        </CardContent>
      </Card>

      {/* Danger Zone */}
      <Card className="border-destructive/50">
        <CardHeader>
          <CardTitle className="text-destructive">Danger Zone</CardTitle>
          <CardDescription>Irreversible actions for this board.</CardDescription>
        </CardHeader>
        <CardContent>
          <DeleteBoardSection boardId={boardId} />
        </CardContent>
      </Card>
    </div>
  )
}
