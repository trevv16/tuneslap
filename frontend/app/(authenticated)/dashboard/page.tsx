import { Metadata } from "next";

import DashboardWrapper from "./DashboardWrapper";

export const metadata: Metadata = {
  title: "Dashboard",
  description: "Dashboard",
};

export default function DashboardPage() {
  return <DashboardWrapper />;
}
