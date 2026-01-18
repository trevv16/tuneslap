"use client"

import type { KeyResponse as Key } from '@/api/models';
import { useCreateKey, useDeleteKey, useGetBoardKeys, useUpdateKey } from '@/hooks/keys';
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle, Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { CheckIcon, EllipsisVerticalIcon, ExclamationTriangleIcon, PlusIcon } from '@heroicons/react/20/solid';
import { useState } from 'react';
import toast from 'react-hot-toast';
import KeyForm from '../KeyForm';

type KeysProps = {
  boardId: string;
}

export default function Keys({ boardId }: KeysProps) {
  const [open, setOpen] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [selectedKey, setSelectedKey] = useState<Key | null>(null);

  // Use the dedicated keys API
  const keysQuery = useGetBoardKeys(boardId);
  const keys = keysQuery.data?.data || [];

  // Mutations
  const createKeyMutation = useCreateKey(boardId);
  const deleteKeyMutation = useDeleteKey(boardId);
  const updateKeyMutation = useUpdateKey(boardId);

  // Handle delete confirmation
  const handleDeleteConfirm = async () => {
    if (!selectedKey) return;

    try {
      if (!selectedKey.id) return;
      await deleteKeyMutation.mutateAsync(selectedKey.id);
      toast.success('Key removed successfully!');
      setDeleteModalOpen(false);
      setSelectedKey(null);
    } catch (error) {
      toast.error('Failed to remove key. Please try again.');
      console.error('Delete key error:', error);
    }
  };

  // Handle key creation
  const handleCreateKey = async (data: { boardId: string; name: string; description?: string; hotKey: string; audioMediaId: string; imageMediaId?: string }) => {
    try {
      await createKeyMutation.mutateAsync(data);
      toast.success('Key created successfully!');
      setOpen(false);
    } catch (error) {
      toast.error('Failed to create key. Please try again.');
      console.error('Add key error:', error);
    }
  };

  // Handle key update
  const handleUpdateKey = async (data: { boardId: string; name: string; description?: string; hotKey: string; audioMediaId: string; imageMediaId?: string }) => {
    if (!selectedKey) return;

    try {
      if (!selectedKey.id) return;
      await updateKeyMutation.mutateAsync({
        keyId: selectedKey.id,
        data: data,
      });
      toast.success('Key updated successfully!');
      setEditModalOpen(false);
      setSelectedKey(null);
    } catch (error) {
      toast.error('Failed to update key. Please try again.');
      console.error('Update key error:', error);
    }
  };

  // Open delete modal
  const openDeleteModal = (key: Key) => {
    setSelectedKey(key);
    setDeleteModalOpen(true);
  };

  // Open edit modal
  const openEditModal = (key: Key) => {
    setSelectedKey(key);
    setEditModalOpen(true);
  };

  return (
    <>
      <div className="mt-8 lg:flex lg:items-center lg:justify-between">
        <div className="min-w-0 flex-1">
          <h3 className="mb-4 text-lg font-semibold text-base">Keys</h3>
        </div>
        <div className="my-5 flex lg:mt-0 lg:ml-4">
          <button
            onClick={() => setOpen(true)}
            className="inline-flex items-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <PlusIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-inverted group-hover:text-inverted" />
            Add Key
          </button>
        </div>
      </div>

      {keys.length === 0 && (
        <div className="bg-elevated p-4 rounded-lg divide-y divide-dark-700">
          <div className="text-center text-sm text-gray-500">No keys found</div>
        </div>
      )}

      {keys.length > 0 && (
        <ul className="bg-elevated p-4 rounded-lg divide-y divide-dark-700">
          {keys.map((key, index) => (
            <li key={(key?.id || 'key-') + index} className="flex justify-between gap-x-6 py-5">
              <div className="flex min-w-0 gap-x-4">
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 bg-accent rounded-lg flex items-center justify-center">
                    <span className="text-inverted font-bold text-lg">{key?.hotKey || '?'}</span>
                  </div>
                </div>
                <div className="min-w-0 flex-auto">
                  <p className="text-sm/6 font-semibold text-base">
                    {key?.name || 'No name'}
                  </p>
                  <p className="mt-1 flex text-xs/5 text-highlight">
                    {key?.description || 'No description'}
                  </p>
                  <div className="mt-2 flex gap-2">
                    {key?.audioUrl && (
                      <audio src={key.audioUrl} controls className="h-8 text-xs" />
                    )}
                    {key?.imageUrl && (
                      <img src={key.imageUrl} alt={key?.name || 'No name'} className="w-8 h-8 object-cover rounded" />
                    )}
                  </div>
                </div>
              </div>
              <div className="flex shrink-0 items-center gap-x-6">
                <Menu as="div" className="relative flex-none">
                  <MenuButton className="-m-2.5 block p-2.5 text-gray-500 hover:text-gray-900">
                    <span className="sr-only">Open options</span>
                    <EllipsisVerticalIcon aria-hidden="true" className="size-5" />
                  </MenuButton>
                  <MenuItems
                    transition
                    className="absolute right-0 z-10 mt-2 w-32 origin-top-right rounded-md bg-highlight py-2 shadow-lg ring-1 ring-gray-900/5 transition focus:outline-hidden data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-100 data-enter:ease-out data-leave:duration-75 data-leave:ease-in"
                  >
                    <MenuItem>
                      <button
                        onClick={() => openEditModal(key)}
                        className="w-full text-left block px-3 py-1 text-sm/6 text-base data-focus:text-inverted data-focus:outline-hidden"
                      >
                        Edit<span className="sr-only">, {key?.name || 'No name'}</span>
                      </button>
                    </MenuItem>
                    <MenuItem>
                      <button
                        onClick={() => openDeleteModal(key)}
                        className="w-full text-left block px-3 py-1 text-sm/6 text-base data-focus:text-inverted data-focus:outline-hidden"
                      >
                        Delete<span className="sr-only">, {key?.name || 'No name'}</span>
                      </button>
                    </MenuItem>
                  </MenuItems>
                </Menu>
              </div>
            </li>
          ))}
        </ul>
      )}

      {/* Add Key Modal */}
      <Dialog open={open} onClose={setOpen} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-lg sm:p-6 data-closed:sm:translate-y-0 data-closed:sm:scale-95"
            >
              <div>
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-success">
                  <CheckIcon aria-hidden="true" className="size-6 text-success" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-gray-900">
                    Add a new key
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Create a new key for this board. Select audio and optionally an image, then assign a hotkey.
                    </p>
                  </div>
                </div>
                <KeyForm
                  boardId={boardId}
                  existingKeys={keys}
                  mode="add"
                  onSubmit={handleCreateKey}
                  onCancel={() => setOpen(false)}
                  isSubmitting={createKeyMutation.isPending}
                />
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog open={deleteModalOpen} onClose={() => setDeleteModalOpen(false)} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-sm sm:p-6 data-closed:sm:translate-y-0 data-closed:sm:scale-95"
            >
              <div>
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-error">
                  <ExclamationTriangleIcon aria-hidden="true" className="size-6 text-error" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-gray-900">
                    Delete key
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Are you sure you want to delete <strong>{selectedKey?.name}</strong>? This action cannot be undone.
                    </p>
                  </div>
                </div>
                <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                  <button
                    type="button"
                    onClick={() => setDeleteModalOpen(false)}
                    disabled={deleteKeyMutation.isPending}
                    className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    onClick={handleDeleteConfirm}
                    disabled={deleteKeyMutation.isPending}
                    className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {deleteKeyMutation.isPending ? 'Deleting...' : 'Delete'}
                  </button>
                </div>
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>

      {/* Edit Key Modal */}
      <Dialog open={editModalOpen} onClose={() => setEditModalOpen(false)} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500/75 transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-lg sm:p-6 data-closed:sm:translate-y-0 data-closed:sm:scale-95"
            >
              <div>
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-accent">
                  <CheckIcon aria-hidden="true" className="size-6 text-accent" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-gray-900">
                    Edit key
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Update the key <strong>{selectedKey?.name}</strong> with new settings.
                    </p>
                  </div>
                </div>
                {selectedKey && (
                  <KeyForm
                    boardId={boardId}
                    existingKeys={keys}
                    mode="edit"
                    initialData={selectedKey}
                    onSubmit={handleUpdateKey}
                    onCancel={() => setEditModalOpen(false)}
                    isSubmitting={updateKeyMutation.isPending}
                  />
                )}
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>
    </>
  )
} 