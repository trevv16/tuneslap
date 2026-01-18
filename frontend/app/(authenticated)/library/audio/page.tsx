import { Metadata } from "next";
import PageTemplate from "../../PageTemplate";
import LibraryClient from "../LibraryClient";


export const metadata: Metadata = {
  title: "Library - Audio",
  description: "Audio Library",
};

export default function AudioLibraryPage() {
  return (
    <PageTemplate>
      <LibraryClient mediaType="audio" />
    </PageTemplate>
  );
}
