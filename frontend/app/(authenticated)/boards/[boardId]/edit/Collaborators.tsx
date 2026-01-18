"use client"

import { useDeleteCollaborator, useGetCollaborators, useUpdateCollaborator } from '@/hooks/collaborators';
import type { UpdateCollaboratorRequestRoleEnum } from '@/api/models';
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle, Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { CheckIcon, EllipsisVerticalIcon, ExclamationTriangleIcon, PlusIcon } from '@heroicons/react/20/solid';
import { zodResolver } from '@hookform/resolvers/zod';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import toast from 'react-hot-toast';
import { z } from 'zod';
import AddCollaboratorForm from './AddCollaboratorForm';

// Validation schema for role change
const roleChangeSchema = z.object({
  role: z.string().min(1, 'Role is required'),
});

type RoleChangeFormData = z.infer<typeof roleChangeSchema>;

type CollaboratorsProps = {
  boardId: string;
}

export default function Collaborators({ boardId }: CollaboratorsProps) {
  const [open, setOpen] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [roleModalOpen, setRoleModalOpen] = useState(false);
  const [selectedCollaborator, setSelectedCollaborator] = useState<{ id: string; name: string; currentRole: string } | null>(null);

  // Use the dedicated collaborators API instead of getting from board data
  const collaboratorsQuery = useGetCollaborators(boardId);
  const collaborators = collaboratorsQuery.data?.data || [];


  // Mutations
  const deleteCollaboratorMutation = useDeleteCollaborator(boardId);
  const updateCollaboratorMutation = useUpdateCollaborator(boardId);

  // Form for role change
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<RoleChangeFormData>({
    resolver: zodResolver(roleChangeSchema),
  });

  // Handle delete confirmation
  const handleDeleteConfirm = async () => {
    if (!selectedCollaborator) return;

    try {
      await deleteCollaboratorMutation.mutateAsync(selectedCollaborator.id);
      toast.success('Collaborator removed successfully!');
      setDeleteModalOpen(false);
      setSelectedCollaborator(null);
    } catch (error) {
      toast.error('Failed to remove collaborator. Please try again.');
      console.error('Delete collaborator error:', error);
    }
  };

  // Handle role change
  const handleRoleChange = async (data: RoleChangeFormData) => {
    if (!selectedCollaborator) return;

    try {
      await updateCollaboratorMutation.mutateAsync({
        collaboratorId: selectedCollaborator.id,
        data: { role: data.role as UpdateCollaboratorRequestRoleEnum },
      });
      toast.success('Collaborator role updated successfully!');
      setRoleModalOpen(false);
      setSelectedCollaborator(null);
      reset();
    } catch (error) {
      toast.error('Failed to update collaborator role. Please try again.');
      console.error('Update collaborator error:', error);
    }
  };

  // Open delete modal
  const openDeleteModal = (collaborator: { id: string; name: string; currentRole: string }) => {
    setSelectedCollaborator(collaborator);
    setDeleteModalOpen(true);
  };

  // Open role change modal
  const openRoleModal = (collaborator: { id: string; name: string; currentRole: string }) => {
    setSelectedCollaborator(collaborator);
    setRoleModalOpen(true);
  };

  return (
    <>
      <div className="mt-8 lg:flex lg:items-center lg:justify-between">
        <div className="min-w-0 flex-1">
          <h3 className="mb-4 text-lg font-semibold text-base">Collaborators</h3>
        </div>
        <div className="my-5 flex lg:mt-0 lg:ml-4">
          <button
            onClick={() => setOpen(true)}
            className="inline-flex items-center rounded-md px-3 py-2 text-sm font-semibold shadow-xs focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent bg-accent text-inverted hover:bg-highlight"
          >
            <PlusIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5" />
            Invite Collaborator
          </button>
        </div>
      </div>

      {collaborators.length === 0 && (
        <div className="bg-elevated p-4 rounded-lg divide-y divide-dark-700">
          <div className="text-center text-sm text-gray-500">No collaborators found</div>
        </div>
      )}

      {collaborators.length > 0 && <ul className="bg-elevated p-4 rounded-lg divide-y divide-dark-700">
        {collaborators.map((person) => (
          <li key={person.id} className="flex justify-between gap-x-6 py-5">
            <div className="flex min-w-0 gap-x-4">
              <img alt="" src={(person.imageUrl && person.imageUrl !== "") ? person.imageUrl : "/defaultUser.jpg"} className="size-12 flex-none rounded-full bg-gray-50" />
              <div className="min-w-0 flex-auto">
                <p className="text-sm/6 font-semibold text-base">
                  {person.name}
                </p>
                <p className="mt-1 flex text-xs/5 text-highlight">
                  {person.role}
                </p>
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
                      onClick={() => openRoleModal({ id: person.id || '', name: person.name || '', currentRole: person.role || 'viewer' })}
                      className="w-full text-left block px-3 py-1 text-sm/6 text-base data-focus:text-inverted data-focus:outline-hidden"
                    >
                      Edit Role<span className="sr-only"> for, {person.name}</span>
                    </button>
                  </MenuItem>
                  <MenuItem>
                    <button
                      onClick={() => openDeleteModal({ id: person.id || '', name: person.name || '', currentRole: person.role || 'viewer' })}
                      className="w-full text-left block px-3 py-1 text-sm/6 text-base data-focus:text-inverted data-focus:outline-hidden"
                    >
                      Remove<span className="sr-only">, {person.name}</span>
                    </button>
                  </MenuItem>
                </MenuItems>
              </Menu>
            </div>
          </li>
        ))}
      </ul>}

      {/* Add Collaborator Modal */}
      <Dialog open={open} onClose={setOpen} className="relative z-10">
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
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-success">
                  <CheckIcon aria-hidden="true" className="size-6 text-success" />
                </div>
                <div className=" text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-gray-900">
                    Invite a collaborator
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Send an invitation to collaborate on this board. If the user doesn&apos;t have an account, they&apos;ll be prompted to create one.
                    </p>
                  </div>
                </div>
                <AddCollaboratorForm boardId={boardId} setOpen={setOpen} />
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
                    Remove collaborator
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Are you sure you want to remove <strong>{selectedCollaborator?.name}</strong> from this board? This action cannot be undone.
                    </p>
                  </div>
                </div>
                <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                  <button
                    type="button"
                    onClick={() => setDeleteModalOpen(false)}
                    disabled={deleteCollaboratorMutation.isPending}
                    className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    onClick={handleDeleteConfirm}
                    disabled={deleteCollaboratorMutation.isPending}
                    className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {deleteCollaboratorMutation.isPending ? 'Removing...' : 'Remove'}
                  </button>
                </div>
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>

      {/* Role Change Modal */}
      <Dialog open={roleModalOpen} onClose={() => setRoleModalOpen(false)} className="relative z-10">
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
                <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-accent">
                  <CheckIcon aria-hidden="true" className="size-6 text-accent" />
                </div>
                <div className="text-center sm:mt-5">
                  <DialogTitle as="h3" className="text-base font-semibold text-gray-900">
                    Change collaborator role
                  </DialogTitle>
                  <div className="mt-2 mb-4">
                    <p className="text-sm text-gray-500">
                      Update the role for <strong>{selectedCollaborator?.name}</strong>. Current role: <strong>{selectedCollaborator?.currentRole}</strong>
                    </p>
                  </div>
                </div>
                <form onSubmit={handleSubmit(handleRoleChange)}>
                  <div className="my-4">
                    <label htmlFor="role" className="block text-sm/6 font-medium text-gray-900">
                      New Role
                    </label>
                    <div className="mt-2">
                      <select
                        id="role"
                        disabled={updateCollaboratorMutation.isPending}
                        className={`block w-full rounded-md bg-white h-8 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-gray-300 sm:text-sm/6 ${errors.role ? 'outline-red-500 focus:outline-red-500' : ''}`}
                        {...register('role')}
                        defaultValue={selectedCollaborator?.currentRole}
                      >
                        <option value="viewer">Viewer</option>
                        <option value="editor">Editor</option>
                        <option value="owner">Owner</option>
                      </select>
                    </div>
                    {errors.role && (
                      <p className="mt-1 text-sm text-red-500">{errors.role.message}</p>
                    )}
                  </div>
                  <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                    <button
                      type="button"
                      onClick={() => setRoleModalOpen(false)}
                      disabled={updateCollaboratorMutation.isPending}
                      className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={updateCollaboratorMutation.isPending}
                      className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {updateCollaboratorMutation.isPending ? 'Updating...' : 'Update Role'}
                    </button>
                  </div>
                </form>
              </div>
            </DialogPanel>
          </div>
        </div>
      </Dialog>
    </>
  )
} 