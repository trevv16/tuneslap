"use client"

import { useGetBoardById } from "@/hooks/boards";
import Collaborators from "./Collaborators";
import DeleteBoardSection from "./DeleteBoardSection";
import EditBoardForm from "./EditBoardForm";
import Keys from "./Keys";

type EditBoardClientProps = {
  boardId: string;
}

export default function EditBoardClient({ boardId }: EditBoardClientProps) {
  const boardQuery = useGetBoardById(boardId)

  return (
    <>
      <EditBoardForm boardQuery={boardQuery} boardId={boardId} />
      <Collaborators boardId={boardId} />
      <Keys boardId={boardId} />
      <DeleteBoardSection boardId={boardId} />
    </>
  );
}