'use client'

import { useCreateCollaborator } from '@/hooks/collaborators'
import type { CreateCollaboratorRequestRoleEnum } from '@/api/models'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { zodResolver } from '@hookform/resolvers/zod'
import { Dispatch, SetStateAction } from "react"
import { useForm, Controller } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

const addCollaboratorSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  role: z.string().min(1, 'Role is required'),
})

type AddCollaboratorFormData = z.infer<typeof addCollaboratorSchema>

interface AddCollaboratorFormProps {
  boardId: string
  setOpen: Dispatch<SetStateAction<boolean>>
}

export default function AddCollaboratorForm({ boardId, setOpen }: AddCollaboratorFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
  } = useForm<AddCollaboratorFormData>({
    resolver: zodResolver(addCollaboratorSchema),
    defaultValues: {
      role: 'viewer',
    },
  })

  const createCollaboratorMutation = useCreateCollaborator(boardId)

  const onSubmit = async (data: AddCollaboratorFormData) => {
    try {
      await createCollaboratorMutation.mutateAsync({
        email: data.email,
        role: data.role as CreateCollaboratorRequestRoleEnum,
      })

      toast.success('Collaborator invitation sent successfully!')
      reset()
      setOpen(false)
    } catch (error) {
      toast.error('Failed to send collaborator invitation. Please try again.')
      console.error('Add collaborator error:', error)
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="email">Email Address</Label>
        <Input
          id="email"
          type="email"
          disabled={createCollaboratorMutation.isPending}
          placeholder="Enter collaborator's email address"
          className="mt-2"
          {...register('email')}
        />
        {errors.email && (
          <p className="mt-1 text-sm text-destructive">{errors.email.message}</p>
        )}
      </div>

      <div>
        <Label htmlFor="role">Role</Label>
        <Controller
          name="role"
          control={control}
          render={({ field }) => (
            <Select
              value={field.value}
              onValueChange={field.onChange}
              disabled={createCollaboratorMutation.isPending}
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

      <div className="flex gap-3 pt-4">
        <Button
          type="button"
          variant="outline"
          onClick={() => { setOpen(false); }}
          disabled={createCollaboratorMutation.isPending}
          className="flex-1"
        >
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={createCollaboratorMutation.isPending}
          className="flex-1"
        >
          {createCollaboratorMutation.isPending ? 'Sending invitation...' : 'Send invitation'}
        </Button>
      </div>
    </form>
  )
}
