import { Metadata } from "next";
import PageTemplate from "../PageTemplate";
import LibraryClient from "./LibraryClient";


export const metadata: Metadata = {
  title: "Library",
  description: "Library",
};

export default function LibraryPage() {
  return (
    <PageTemplate>
      <LibraryClient />
    </PageTemplate>
  );
}
