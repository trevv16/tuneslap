import { Metadata } from "next";
import PageTemplate from "../../../PageTemplate";
import EditBoardClient from "./EditBoardClient";

export const metadata: Metadata = {
  title: "Edit Board",
  description: "Edit Board",
};

export default async function EditBoardPage({
  params,
}: {
  params: Promise<{ boardId: string }>;
}) {
  const { boardId } = await params;
  return (
    <PageTemplate>
      <EditBoardClient boardId={boardId} />
    </PageTemplate>
  );
}
