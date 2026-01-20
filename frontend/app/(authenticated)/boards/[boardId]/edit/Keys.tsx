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
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { useCreateKey, useDeleteKey, useGetBoardKeys, useUpdateKey } from '@/hooks/keys'
import { AlertTriangle, Keyboard, MoreVertical, Plus } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'
import KeyForm from '../KeyForm'

type KeysProps = {
  boardId: string
}

export default function Keys({ boardId }: KeysProps) {
  const [sheetOpen, setSheetOpen] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [selectedKey, setSelectedKey] = useState<Key | null>(null)
  const [mode, setMode] = useState<'add' | 'edit'>('add')

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
      setSheetOpen(false)
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
      setSheetOpen(false)
      setSelectedKey(null)
    } catch (error) {
      toast.error('Failed to update key. Please try again.')
      console.error('Update key error:', error)
    }
  }

  const openAddSheet = () => {
    setMode('add')
    setSelectedKey(null)
    setSheetOpen(true)
  }

  const openEditSheet = (key: Key) => {
    setMode('edit')
    setSelectedKey(key)
    setSheetOpen(true)
  }

  const openDeleteModal = (key: Key) => {
    setSelectedKey(key)
    setDeleteModalOpen(true)
  }

  return (
    <>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-muted-foreground">
          {keys.length === 0 
            ? 'No keys configured' 
            : `${keys.length} key${keys.length === 1 ? '' : 's'}`
          }
        </p>
        <Button size="sm" onClick={openAddSheet}>
          <Plus className="mr-1.5 h-4 w-4" />
          Add Key
        </Button>
      </div>

      {keys.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
          {keys.map((key, index) => (
            <div
              key={(key?.id || 'key-') + index}
              className="relative group rounded-lg border bg-muted/30 p-3 hover:bg-muted/50 transition-colors"
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex h-8 w-8 items-center justify-center rounded bg-primary text-primary-foreground font-mono font-bold text-sm">
                  {key?.hotKey?.toUpperCase() || '?'}
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button 
                      variant="ghost" 
                      size="icon" 
                      className="h-7 w-7 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <MoreVertical className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => openEditSheet(key)}>
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
              
              {key?.imageUrl && (
                <div className="mb-2 rounded overflow-hidden aspect-square">
                  <img 
                    src={key.imageUrl} 
                    alt={key?.name || 'Key image'} 
                    className="w-full h-full object-cover"
                  />
                </div>
              )}
              
              <p className="text-sm font-medium truncate">{key?.name || 'Unnamed'}</p>
              {key?.description && (
                <p className="text-xs text-muted-foreground truncate mt-0.5">{key.description}</p>
              )}
            </div>
          ))}
        </div>
      )}

      {keys.length === 0 && (
        <div className="flex flex-col items-center justify-center py-8 text-center">
          <Keyboard className="h-10 w-10 text-muted-foreground mb-3" />
          <p className="text-sm text-muted-foreground">Add keys to create your soundboard.</p>
        </div>
      )}

      {/* Add/Edit Key Sheet */}
      <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
        <SheetContent className="sm:max-w-lg">
          <SheetHeader>
            <SheetTitle>{mode === 'add' ? 'Add Key' : 'Edit Key'}</SheetTitle>
            <SheetDescription>
              {mode === 'add' 
                ? 'Create a new key for this board.' 
                : `Update settings for "${selectedKey?.name}".`
              }
            </SheetDescription>
          </SheetHeader>
          <div className="flex-1 overflow-y-auto px-4 pb-4">
            {mode === 'add' ? (
              <KeyForm
                boardId={boardId}
                existingKeys={keys}
                mode="add"
                onSubmit={handleCreateKey}
                onCancel={() => setSheetOpen(false)}
                isSubmitting={createKeyMutation.isPending}
              />
            ) : selectedKey && (
              <KeyForm
                boardId={boardId}
                existingKeys={keys}
                mode="edit"
                initialData={selectedKey}
                onSubmit={handleUpdateKey}
                onCancel={() => setSheetOpen(false)}
                isSubmitting={updateKeyMutation.isPending}
              />
            )}
          </div>
        </SheetContent>
      </Sheet>

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
    </>
  )
}
