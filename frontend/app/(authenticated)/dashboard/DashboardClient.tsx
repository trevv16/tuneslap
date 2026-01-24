"use client";

import EmptyState from "@/components/EmptyState";
import { useGetBoards } from "@/hooks/boards";
import BoardsList from "./BoardsList";

interface DashboardClientProps {
  onOpenCreateModal: () => void;
}

export default function DashboardClient({ onOpenCreateModal }: DashboardClientProps) {
  const { data: boards } = useGetBoards();

  const boardsArray = Array.isArray(boards) ? boards : [];
  if (!boardsArray || boardsArray.length === 0) {
    return (
      <EmptyState
        title="Create Your First Soundboard"
        description="Get started by creating a new board."
        buttonText="New Board"
        buttonOnClick={onOpenCreateModal}
      />
    )
  }

  return (
    <BoardsList boards={boardsArray} />
  )
}
