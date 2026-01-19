"use client"

import SoundBoard from "@/components/SoundBoard"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { useGetBoardById } from "@/hooks/boards"
import { useCreateKey } from "@/hooks/keys"
import { Pencil, Plus } from "lucide-react"
import Link from "next/link"
import { useParams } from "next/navigation"
import { useState } from "react"
import { toast } from "sonner"
import Header from "../../Header"
import KeyForm from "./KeyForm"

type HeaderActionsProps = {
  boardId: string
}

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

  const HeaderActions = ({ boardId }: HeaderActionsProps) => {
    return (
      <div className="flex items-center gap-3">
        <Button variant="outline" asChild>
          <Link href={`/boards/${boardId}/edit`}>
            <Pencil className="mr-1.5 -ml-0.5 h-5 w-5" />
            Edit
          </Link>
        </Button>
        <Button onClick={() => setAddKeyOpen(true)}>
          <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
          Add Key
        </Button>
      </div>
    )
  }

  return (
    <>
      <Header pageTitle={board?.name || "Board Detail"} headerActions={<HeaderActions boardId={boardId as string} />} />
      <SoundBoard keys={board?.keys || []} />

      {/* Add Key Button */}
      <div className="mt-6 flex justify-center">
        <Button onClick={() => setAddKeyOpen(true)}>
          <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
          Add Key
        </Button>
      </div>

      {/* Add Key Modal */}
      <Dialog open={addKeyOpen} onOpenChange={setAddKeyOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Plus className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle className="text-center">Add a new key</DialogTitle>
            <DialogDescription className="text-center">
              Create a new key for this board. Select audio (required) and optionally an image, then assign a hotkey.
            </DialogDescription>
          </DialogHeader>
          <KeyForm
            boardId={boardId as string}
            existingKeys={board?.keys || []}
            mode="add"
            onSubmit={handleAddKey}
            onCancel={() => setAddKeyOpen(false)}
            isSubmitting={createKeyMutation.isPending}
          />
        </DialogContent>
      </Dialog>
    </>
  )
}
