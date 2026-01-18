import { Metadata } from "next";
import PageTemplate from "../../PageTemplate";
import LibraryClient from "../LibraryClient";


export const metadata: Metadata = {
  title: "Library - Images",
  description: "Image Library",
};

export default function ImageLibraryPage() {
  return (
    <PageTemplate>
      <LibraryClient mediaType="image" />
    </PageTemplate>
  );
}
