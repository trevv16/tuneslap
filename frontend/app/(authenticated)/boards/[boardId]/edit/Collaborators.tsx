"use client"

import { useDeleteCollaborator, useGetCollaborators, useUpdateCollaborator } from '@/hooks/collaborators'
import type { UpdateCollaboratorRequestRoleEnum } from '@/api/models'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
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
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { zodResolver } from '@hookform/resolvers/zod'
import { AlertTriangle, Check, MoreVertical, Plus } from 'lucide-react'
import { useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'
import AddCollaboratorForm from './AddCollaboratorForm'

const roleChangeSchema = z.object({
  role: z.string().min(1, 'Role is required'),
})

type RoleChangeFormData = z.infer<typeof roleChangeSchema>

type CollaboratorsProps = {
  boardId: string
}

export default function Collaborators({ boardId }: CollaboratorsProps) {
  const [open, setOpen] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [roleModalOpen, setRoleModalOpen] = useState(false)
  const [selectedCollaborator, setSelectedCollaborator] = useState<{ id: string; name: string; currentRole: string } | null>(null)

  const collaboratorsQuery = useGetCollaborators(boardId)
  const collaborators = collaboratorsQuery.data?.data || []

  const deleteCollaboratorMutation = useDeleteCollaborator(boardId)
  const updateCollaboratorMutation = useUpdateCollaborator(boardId)

  const {
    handleSubmit,
    formState: { errors },
    reset,
    control,
  } = useForm<RoleChangeFormData>({
    resolver: zodResolver(roleChangeSchema),
  })

  const handleDeleteConfirm = async () => {
    if (!selectedCollaborator) return

    try {
      await deleteCollaboratorMutation.mutateAsync(selectedCollaborator.id)
      toast.success('Collaborator removed successfully!')
      setDeleteModalOpen(false)
      setSelectedCollaborator(null)
    } catch (error) {
      toast.error('Failed to remove collaborator. Please try again.')
      console.error('Delete collaborator error:', error)
    }
  }

  const handleRoleChange = async (data: RoleChangeFormData) => {
    if (!selectedCollaborator) return

    try {
      await updateCollaboratorMutation.mutateAsync({
        collaboratorId: selectedCollaborator.id,
        data: { role: data.role as UpdateCollaboratorRequestRoleEnum },
      })
      toast.success('Collaborator role updated successfully!')
      setRoleModalOpen(false)
      setSelectedCollaborator(null)
      reset()
    } catch (error) {
      toast.error('Failed to update collaborator role. Please try again.')
      console.error('Update collaborator error:', error)
    }
  }

  const openDeleteModal = (collaborator: { id: string; name: string; currentRole: string }) => {
    setSelectedCollaborator(collaborator)
    setDeleteModalOpen(true)
  }

  const openRoleModal = (collaborator: { id: string; name: string; currentRole: string }) => {
    setSelectedCollaborator(collaborator)
    setRoleModalOpen(true)
  }

  return (
    <>
      <div className="mt-8 lg:flex lg:items-center lg:justify-between">
        <div className="min-w-0 flex-1">
          <h3 className="mb-4 text-lg font-semibold text-foreground">Collaborators</h3>
        </div>
        <div className="my-5 flex lg:mt-0 lg:ml-4">
          <Button onClick={() => setOpen(true)}>
            <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
            Invite Collaborator
          </Button>
        </div>
      </div>

      {collaborators.length === 0 && (
        <div className="bg-card p-4 rounded-lg border">
          <div className="text-center text-sm text-muted-foreground">No collaborators found</div>
        </div>
      )}

      {collaborators.length > 0 && (
        <ul className="bg-card p-4 rounded-lg border divide-y">
          {collaborators.map((person) => (
            <li key={person.id} className="flex justify-between gap-x-6 py-5">
              <div className="flex min-w-0 gap-x-4">
                <Avatar>
                  <AvatarImage src={person.imageUrl && person.imageUrl !== "" ? person.imageUrl : undefined} />
                  <AvatarFallback>{person.name?.charAt(0).toUpperCase() || 'U'}</AvatarFallback>
                </Avatar>
                <div className="min-w-0 flex-auto">
                  <p className="text-sm font-semibold text-foreground">
                    {person.name}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground capitalize">
                    {person.role}
                  </p>
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
                    <DropdownMenuItem onClick={() => openRoleModal({ id: person.id || '', name: person.name || '', currentRole: person.role || 'viewer' })}>
                      Edit Role
                    </DropdownMenuItem>
                    <DropdownMenuItem 
                      onClick={() => openDeleteModal({ id: person.id || '', name: person.name || '', currentRole: person.role || 'viewer' })}
                      className="text-destructive focus:text-destructive"
                    >
                      Remove
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </li>
          ))}
        </ul>
      )}

      {/* Add Collaborator Modal */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Check className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle className="text-center">Invite a collaborator</DialogTitle>
            <DialogDescription className="text-center">
              Send an invitation to collaborate on this board. If the user doesn&apos;t have an account, they&apos;ll be prompted to create one.
            </DialogDescription>
          </DialogHeader>
          <AddCollaboratorForm boardId={boardId} setOpen={setOpen} />
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog open={deleteModalOpen} onOpenChange={setDeleteModalOpen}>
        <DialogContent>
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>
            <DialogTitle className="text-center">Remove collaborator</DialogTitle>
            <DialogDescription className="text-center">
              Are you sure you want to remove <strong>{selectedCollaborator?.name}</strong> from this board? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="sm:grid sm:grid-cols-2 sm:gap-3">
            <Button
              variant="outline"
              onClick={() => setDeleteModalOpen(false)}
              disabled={deleteCollaboratorMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteConfirm}
              disabled={deleteCollaboratorMutation.isPending}
            >
              {deleteCollaboratorMutation.isPending ? 'Removing...' : 'Remove'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Role Change Modal */}
      <Dialog open={roleModalOpen} onOpenChange={setRoleModalOpen}>
        <DialogContent>
          <DialogHeader>
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Check className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle className="text-center">Change collaborator role</DialogTitle>
            <DialogDescription className="text-center">
              Update the role for <strong>{selectedCollaborator?.name}</strong>. Current role: <strong>{selectedCollaborator?.currentRole}</strong>
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit(handleRoleChange)} className="space-y-4">
            <div>
              <Label htmlFor="role">New Role</Label>
              <Controller
                name="role"
                control={control}
                defaultValue={selectedCollaborator?.currentRole}
                render={({ field }) => (
                  <Select
                    value={field.value}
                    onValueChange={field.onChange}
                    disabled={updateCollaboratorMutation.isPending}
                  >
                    <SelectTrigger className="mt-2">
                      <SelectValue placeholder="Select role" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="viewer">Viewer</SelectItem>
                      <SelectItem value="editor">Editor</SelectItem>
                      <SelectItem value="owner">Owner</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.role && (
                <p className="mt-1 text-sm text-destructive">{errors.role.message}</p>
              )}
            </div>
            <DialogFooter className="sm:grid sm:grid-cols-2 sm:gap-3">
              <Button
                type="button"
                variant="outline"
                onClick={() => setRoleModalOpen(false)}
                disabled={updateCollaboratorMutation.isPending}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={updateCollaboratorMutation.isPending}
              >
                {updateCollaboratorMutation.isPending ? 'Updating...' : 'Update Role'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  )
}
