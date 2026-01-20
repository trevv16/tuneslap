"use client"

import { useDeleteCollaborator, useGetCollaborators, useUpdateCollaborator } from '@/hooks/collaborators'
import type { UpdateCollaboratorRequestRoleEnum } from '@/api/models'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
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
import { AlertTriangle, MoreVertical, Plus, UserPlus } from 'lucide-react'
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
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-muted-foreground">
          {collaborators.length === 0 
            ? 'No collaborators yet' 
            : `${collaborators.length} collaborator${collaborators.length === 1 ? '' : 's'}`
          }
        </p>
        <Button size="sm" onClick={() => setOpen(true)}>
          <Plus className="mr-1.5 h-4 w-4" />
          Invite
        </Button>
      </div>

      {collaborators.length > 0 && (
        <div className="space-y-2">
          {collaborators.map((person) => (
            <div key={person.id} className="flex items-center justify-between p-3 rounded-lg bg-muted/50">
              <div className="flex items-center gap-3">
                <Avatar className="h-9 w-9">
                  <AvatarImage src={person.imageUrl && person.imageUrl !== "" ? person.imageUrl : undefined} />
                  <AvatarFallback className="text-xs">{person.name?.charAt(0).toUpperCase() || 'U'}</AvatarFallback>
                </Avatar>
                <div>
                  <p className="text-sm font-medium">{person.name}</p>
                  <Badge variant="outline" className="text-xs capitalize">{person.role}</Badge>
                </div>
              </div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-8 w-8">
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => openRoleModal({ id: person.id || '', name: person.name || '', currentRole: person.role || 'viewer' })}>
                    Change Role
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
          ))}
        </div>
      )}

      {collaborators.length === 0 && (
        <div className="flex flex-col items-center justify-center py-8 text-center">
          <UserPlus className="h-10 w-10 text-muted-foreground mb-3" />
          <p className="text-sm text-muted-foreground">Invite collaborators to work on this board together.</p>
        </div>
      )}

      {/* Add Collaborator Modal */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Invite Collaborator</DialogTitle>
            <DialogDescription>
              Send an invitation to collaborate on this board.
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
              Are you sure you want to remove <strong>{selectedCollaborator?.name}</strong>? They will lose access to this board.
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
            <DialogTitle>Change Role</DialogTitle>
            <DialogDescription>
              Update the role for <strong>{selectedCollaborator?.name}</strong>.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit(handleRoleChange)} className="space-y-4">
            <div>
              <Label htmlFor="role">Role</Label>
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
                    <SelectTrigger className="mt-1.5">
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
            <DialogFooter>
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
                {updateCollaboratorMutation.isPending ? 'Updating...' : 'Update'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  )
}
