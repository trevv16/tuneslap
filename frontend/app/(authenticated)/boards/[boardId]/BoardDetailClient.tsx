"use client";

import SoundBoard from "@/components/SoundBoard";
import { useGetBoardById } from "@/hooks/boards";
import { useCreateKey } from "@/hooks/keys";
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle } from "@headlessui/react";
import { PencilIcon, PlusIcon } from "@heroicons/react/20/solid";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useState } from "react";
import toast from "react-hot-toast";
import Header from "../../Header";
import KeyForm from "./KeyForm";

type HeaderActionsProps = {
  boardId: string;
}


export default function BoardDetailClient() {
  const { boardId } = useParams();
  const { data: board } = useGetBoardById(boardId as string);
  const createKeyMutation = useCreateKey(boardId as string);

  const [addKeyOpen, setAddKeyOpen] = useState(false);

  const handleAddKey = async (data: { name: string; description?: string; hotKey: string; audioMediaId: string; imageMediaId?: string; boardId: string }) => {
    try {
      await createKeyMutation.mutateAsync(data);
      toast.success('Key added successfully!');
      setAddKeyOpen(false);
    } catch (error) {
      toast.error('Failed to add key. Please try again.');
      console.error('Add key error:', error);
    }
  };

  const HeaderActions = ({ boardId }: HeaderActionsProps) => {
    return (
      <>
        <span className="hidden sm:block group">
          <Link
            href={`/boards/${boardId}/edit`}
            className="inline-flex items-center rounded-md bg-elevated px-3 py-2 text-sm font-semibold text-highlight shadow-xs ring-1 border border-muted ring-inset group-hover:bg-highlight group-hover:text-inverted"
          >
            <PencilIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-highlight group-hover:text-inverted" />
            Edit
          </Link>
        </span>
        <span className="sm:ml-3">
          <button
            onClick={() => setAddKeyOpen(true)}
            className="inline-flex items-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <PlusIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-inverted group-hover:text-inverted" />
            Add Key
          </button>
        </span>
      </>
    )
  }

  return (
    <>
      <Header pageTitle={board?.name || "Board Detail"} headerActions={<HeaderActions boardId={boardId as string} />} />
      <SoundBoard keys={board?.keys || []} />

      {/* Add Key Button */}
      <div className="mt-6 flex justify-center">
        <button
          className="inline-flex items-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          onClick={() => setAddKeyOpen(true)}
        >
          <PlusIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-inverted group-hover:text-inverted" />
          Add Key
        </button>
      </div>

      {/* Add Key Modal */}
      <Dialog open={addKeyOpen} onClose={() => setAddKeyOpen(false)} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-elevated px-4 pt-5 pb-4 text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-sm sm:p-6 data-closed:sm:translate-y-0 data-closed:sm:scale-95"
            >
              <div>
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-accent">
                  <PlusIcon aria-hidden="true" className="size-6 text-accent" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-highlight">
                    Add a new key
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-muted">
                      Create a new key for this board. Select audio (required) and optionally an image, then assign a hotkey.
                    </p>
                  </div>
                </div>
                <KeyForm
                  boardId={boardId as string}
                  existingKeys={board?.keys || []}
                  mode="add"
                  onSubmit={handleAddKey}
                  onCancel={() => setAddKeyOpen(false)}
                  isSubmitting={createKeyMutation.isPending}
                />
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>
    </>
  )
}