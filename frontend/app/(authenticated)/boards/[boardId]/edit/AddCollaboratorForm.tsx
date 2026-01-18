'use client'

import { useCreateCollaborator } from '@/hooks/collaborators';
import type { CreateCollaboratorRequestRoleEnum } from '@/api/models';
import { zodResolver } from '@hookform/resolvers/zod';
import { Dispatch, SetStateAction } from "react";
import { useForm } from 'react-hook-form';
import toast from 'react-hot-toast';
import { z } from 'zod';

// Validation schema
const addCollaboratorSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  role: z.string().min(1, 'Role is required'),
})

type AddCollaboratorFormData = z.infer<typeof addCollaboratorSchema>

type AddCollaboratorFormProps = {
  boardId: string;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export default function AddCollaboratorForm({ boardId, setOpen }: AddCollaboratorFormProps) {

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
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
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="my-12">
        <div>
          <label htmlFor="email" className="block text-sm/6 font-medium text-gray-900">
            Email Address
          </label>
          <div className="mt-2">
            <input
              id="email"
              type="email"
              disabled={createCollaboratorMutation.isPending}
              className={`block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-gray-300 sm:text-sm/6 ${errors.email ? 'outline-red-500 focus:outline-red-500' : ''}`}
              placeholder="Enter collaborator's email address"
              {...register('email')}
            />
          </div>
          {errors.email && (
            <p className="mt-1 text-sm text-red-500">{errors.email.message}</p>
          )}
        </div>
        <div className="mt-4">
          <label htmlFor="role" className="block text-sm/6 font-medium text-gray-900">
            Role
          </label>
          <div className="mt-2">
            <select
              id="role"
              disabled={createCollaboratorMutation.isPending}
              className={`block w-full rounded-md bg-white h-8 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-gray-300 sm:text-sm/6 ${errors.role ? 'outline-red-500 focus:outline-red-500' : ''}`}
              {...register('role')}
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
      </div>
      <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
        <button
          type="button"
          onClick={() => setOpen(false)}
          disabled={createCollaboratorMutation.isPending}
          className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={createCollaboratorMutation.isPending}
          className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {createCollaboratorMutation.isPending ? 'Sending invitation...' : 'Send invitation'}
        </button>
      </div>
    </form>
  )
} 