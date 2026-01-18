import { Metadata } from "next";
import PageTemplate from "../PageTemplate";
import AccountClient from "./AccountClient";
import AccountMenu from "./AccountMenu";

export const metadata: Metadata = {
  title: "Account",
  description: "Account",
};

export default function AccountPage() {
  return (
    <PageTemplate>
      <AccountMenu />
      <AccountClient />
    </PageTemplate>
  );
}
