"use client"

import type { KeyResponse as Key } from '@/api/models'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useCreateKey, useDeleteKey, useGetBoardKeys, useUpdateKey } from '@/hooks/keys'
import { AlertTriangle, Check, MoreVertical, Plus } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'
import KeyForm from '../KeyForm'

type KeysProps = {
  boardId: string
}

export default function Keys({ boardId }: KeysProps) {
  const [open, setOpen] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [selectedKey, setSelectedKey] = useState<Key | null>(null)

  const keysQuery = useGetBoardKeys(boardId)
  const keys = keysQuery.data?.data || []

  const createKeyMutation = useCreateKey(boardId)
  const deleteKeyMutation = useDeleteKey(boardId)
  const updateKeyMutation = useUpdateKey(boardId)

  const handleDeleteConfirm = async () => {
    if (!selectedKey) return

    try {
      if (!selectedKey.id) return
      await deleteKeyMutation.mutateAsync(selectedKey.id)
      toast.success('Key removed successfully!')
      setDeleteModalOpen(false)
      setSelectedKey(null)
    } catch (error) {
      toast.error('Failed to remove key. Please try again.')
      console.error('Delete key error:', error)
    }
  }

  const handleCreateKey = async (data: { boardId: string; name: string; description?: string; hotKey: string; audioMediaId: string; imageMediaId?: string }) => {
    try {
      await createKeyMutation.mutateAsync(data)
      toast.success('Key created successfully!')
      setOpen(false)
    } catch (error) {
      toast.error('Failed to create key. Please try again.')
      console.error('Add key error:', error)
    }
  }

  const handleUpdateKey = async (data: { boardId: string; name: string; description?: string; hotKey: string; audioMediaId: string; imageMediaId?: string }) => {
    if (!selectedKey) return

    try {
      if (!selectedKey.id) return
      await updateKeyMutation.mutateAsync({
        keyId: selectedKey.id,
        data: data,
      })
      toast.success('Key updated successfully!')
      setEditModalOpen(false)
      setSelectedKey(null)
    } catch (error) {
      toast.error('Failed to update key. Please try again.')
      console.error('Update key error:', error)
    }
  }

  const openDeleteModal = (key: Key) => {
    setSelectedKey(key)
    setDeleteModalOpen(true)
  }

  const openEditModal = (key: Key) => {
    setSelectedKey(key)
    setEditModalOpen(true)
  }

  return (
    <>
      <div className="mt-8 lg:flex lg:items-center lg:justify-between">
        <div className="min-w-0 flex-1">
          <h3 className="mb-4 text-lg font-semibold text-foreground">Keys</h3>
        </div>
        <div className="my-5 flex lg:mt-0 lg:ml-4">
          <Button onClick={() => setOpen(true)}>
            <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
            Add Key
          </Button>
        </div>
      </div>

      {keys.length === 0 && (
        <div className="bg-card p-4 rounded-lg border">
          <div className="text-center text-sm text-muted-foreground">No keys found</div>
        </div>
      )}

      {keys.length > 0 && (
        <ul className="bg-card p-4 rounded-lg border divide-y">
          {keys.map((key, index) => (
            <li key={(key?.id || 'key-') + index} className="flex justify-between gap-x-6 py-5">
              <div className="flex min-w-0 gap-x-4">
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 bg-primary rounded-lg flex items-center justify-center">
                    <span className="text-primary-foreground font-bold text-lg">{key?.hotKey || '?'}</span>
                  </div>
                </div>
                <div className="min-w-0 flex-auto">
                  <p className="text-sm font-semibold text-foreground">
                    {key?.name || 'No name'}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground">
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
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <span className="sr-only">Open options</span>
                      <MoreVertical className="h-5 w-5" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => openEditModal(key)}>
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem 
                      onClick={() => openDeleteModal(key)}
                      className="text-destructive focus:text-destructive"
                    >
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </li>
          ))}
        </ul>
      )}

      {/* Add Key Modal */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Check className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle className="text-center">Add a new key</DialogTitle>
            <DialogDescription className="text-center">
              Create a new key for this board. Select audio and optionally an image, then assign a hotkey.
            </DialogDescription>
          </DialogHeader>
          <KeyForm
            boardId={boardId}
            existingKeys={keys}
            mode="add"
            onSubmit={handleCreateKey}
            onCancel={() => setOpen(false)}
            isSubmitting={createKeyMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog open={deleteModalOpen} onOpenChange={setDeleteModalOpen}>
        <DialogContent>
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>
            <DialogTitle className="text-center">Delete key</DialogTitle>
            <DialogDescription className="text-center">
              Are you sure you want to delete <strong>{selectedKey?.name}</strong>? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="sm:grid sm:grid-cols-2 sm:gap-3">
            <Button
              variant="outline"
              onClick={() => setDeleteModalOpen(false)}
              disabled={deleteKeyMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteConfirm}
              disabled={deleteKeyMutation.isPending}
            >
              {deleteKeyMutation.isPending ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Key Modal */}
      <Dialog open={editModalOpen} onOpenChange={setEditModalOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Check className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle className="text-center">Edit key</DialogTitle>
            <DialogDescription className="text-center">
              Update the key <strong>{selectedKey?.name}</strong> with new settings.
            </DialogDescription>
          </DialogHeader>
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
        </DialogContent>
      </Dialog>
    </>
  )
}
