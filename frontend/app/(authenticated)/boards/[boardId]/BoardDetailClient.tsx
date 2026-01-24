"use client"

import SoundBoard from "@/components/SoundBoard"
import { Button } from "@/components/ui/button"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"
import { useGetBoardById } from "@/hooks/boards"
import { useCreateKey } from "@/hooks/keys"
import { Pencil, Plus } from "lucide-react"
import Link from "next/link"
import { useParams } from "next/navigation"
import { useState } from "react"
import { toast } from "sonner"
import Header from "../../Header"
import KeyForm from "./KeyForm"

export default function BoardDetailClient() {
  const { boardId } = useParams()
  const { data: board } = useGetBoardById(boardId as string)
  const createKeyMutation = useCreateKey(boardId as string)

  const [addKeyOpen, setAddKeyOpen] = useState(false)

  const handleAddKey = async (data: { name: string; description?: string; hotKey: string; audioMediaId: string; imageMediaId?: string; boardId: string }) => {
    try {
      await createKeyMutation.mutateAsync(data)
      toast.success('Key added successfully!')
      setAddKeyOpen(false)
    } catch (error) {
      toast.error('Failed to add key. Please try again.')
      console.error('Add key error:', error)
    }
  }

  return (
    <>
      <Header
        pageTitle={board?.name || "Board Detail"}
        headerActions={
          <div className="flex items-center gap-3">
            <Button variant="outline" asChild>
              <Link href={`/boards/${boardId}/edit`}>
                <Pencil className="mr-1.5 -ml-0.5 h-5 w-5" />
                Edit
              </Link>
            </Button>
            <Button onClick={() => { setAddKeyOpen(true); }}>
              <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
              Add Key
            </Button>
          </div>
        }
      />
      <SoundBoard keys={board?.keys || []} onAddKey={() => { setAddKeyOpen(true); }} />

      {/* Add Key Sheet */}
      <Sheet open={addKeyOpen} onOpenChange={setAddKeyOpen}>
        <SheetContent className="sm:max-w-lg">
          <SheetHeader>
            <SheetTitle>Add Key</SheetTitle>
            <SheetDescription>
              Create a new key for this board.
            </SheetDescription>
          </SheetHeader>
          <div className="flex-1 overflow-y-auto px-4 pb-4">
            <KeyForm
              boardId={boardId as string}
              existingKeys={board?.keys || []}
              mode="add"
              onSubmit={handleAddKey}
              onCancel={() => { setAddKeyOpen(false); }}
              isSubmitting={createKeyMutation.isPending}
            />
          </div>
        </SheetContent>
      </Sheet>
    </>
  )
}
