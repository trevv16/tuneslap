import PageTemplate from "../../PageTemplate";
import BoardDetailClient from "./BoardDetailClient";

// export async function generateMetadata({ params }: { params: { boardId: string } }): Promise<Metadata> {
//   const { boardId } = params;
//   const authToken = getStoredToken() || "";

//   const board = await getBoardById(authToken, { boardId: boardId as string });

//   return {
//     title: board?.data?.name || "Board Detail",
//     description: board?.data?.description || "Board Detail",
//   };
// }

export default function BoardDetailPage() {
  return (
    <PageTemplate>
      <BoardDetailClient />
    </PageTemplate>
  );
}
