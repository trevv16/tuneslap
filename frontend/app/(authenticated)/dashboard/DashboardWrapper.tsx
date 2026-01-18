'use client';

import { useState } from 'react';
import Header from "../Header";
import PageTemplate from "../PageTemplate";
import DashboardClient from "./DashboardClient";
import HeaderActions from "./HeaderActions";

export default function DashboardWrapper() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  return (
    <PageTemplate>
      <Header
        pageTitle="Dashboard"
        headerActions={<HeaderActions open={isCreateModalOpen} setOpen={setIsCreateModalOpen} />}
      />
      <DashboardClient onOpenCreateModal={() => setIsCreateModalOpen(true)} />
    </PageTemplate>
  );
}
