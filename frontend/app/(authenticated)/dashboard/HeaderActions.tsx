'use client'

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Check, Plus } from "lucide-react"
import { Dispatch, SetStateAction } from "react"
import CreateBoardForm from "./CreateBoardForm"

type HeaderActionsProps = {
  open: boolean
  setOpen: Dispatch<SetStateAction<boolean>>
}

export default function HeaderActions({ open, setOpen }: HeaderActionsProps) {
  return (
    <>
      <Button onClick={() => setOpen(true)}>
        <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
        Create
      </Button>
      
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Check className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle className="text-center">Create a new board</DialogTitle>
            <DialogDescription className="text-center">
              Create a new board to get started.
            </DialogDescription>
          </DialogHeader>
          <CreateBoardForm setOpen={setOpen} />
        </DialogContent>
      </Dialog>
    </>
  )
}
