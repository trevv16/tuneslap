"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useDeleteBoard } from "@/hooks/boards";
import { AlertTriangle } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";

interface DeleteBoardSectionProps {
  boardId: string;
};

export default function DeleteBoardSection({ boardId }: DeleteBoardSectionProps) {
  const router = useRouter();
  const deleteBoardMutation = useDeleteBoard();
  const [isDeleting, setIsDeleting] = useState(false);
  const [open, setOpen] = useState(false);

  const confirmDelete = async () => {
    setIsDeleting(true);
    try {
      await deleteBoardMutation.mutateAsync(boardId);
      setOpen(false);
      router.push("/dashboard");
    } catch (error) {
      console.error("Failed to delete board:", error);
      setIsDeleting(false);
      setOpen(false);
    }
  };

  return (
    <>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium">Delete this board</p>
          <p className="text-sm text-muted-foreground">
            Once deleted, all data will be permanently removed.
          </p>
        </div>
        <Button
          variant="destructive"
          onClick={() => { setOpen(true); }}
        >
          Delete Board
        </Button>
      </div>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>
            <DialogTitle className="text-center">Delete board</DialogTitle>
            <DialogDescription className="text-center">
              Are you sure? This action cannot be undone. All keys, collaborators, and settings will be permanently deleted.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="sm:grid sm:grid-cols-2 sm:gap-3">
            <Button
              variant="outline"
              onClick={() => { setOpen(false); }}
              disabled={isDeleting}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={confirmDelete}
              disabled={isDeleting}
            >
              {isDeleting ? "Deleting..." : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
