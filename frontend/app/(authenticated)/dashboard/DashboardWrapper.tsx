'use client';

import { useState } from 'react';
import DemoBanner from "@/components/DemoBanner";
import Header from "../Header";
import PageTemplate from "../PageTemplate";
import DashboardClient from "./DashboardClient";
import HeaderActions from "./HeaderActions";

export default function DashboardWrapper() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  return (
    <PageTemplate>
      <DemoBanner message="This is a demo environment. All data will be deleted within one hour." />
      <Header
        pageTitle="Dashboard"
        headerActions={<HeaderActions open={isCreateModalOpen} setOpen={setIsCreateModalOpen} />}
      />
      <DashboardClient onOpenCreateModal={() => setIsCreateModalOpen(true)} />
    </PageTemplate>
  );
}
