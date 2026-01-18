'use client'

import { useCreateBoard } from '@/hooks/boards';
import type { CreateBoardRequestLayoutEnum } from '@/api/models';
import { zodResolver } from '@hookform/resolvers/zod';
import { Dispatch, SetStateAction } from "react";
import { useForm } from 'react-hook-form';
import toast from 'react-hot-toast';
import { z } from 'zod';

// Validation schema
const createBoardSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  description: z.string().min(1, 'Description is required'),
  layout: z.string().min(1, 'Layout is required'),
})

type CreateBoardFormData = z.infer<typeof createBoardSchema>

type CreateBoardFormProps = {
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export default function CreateBoardForm({ setOpen }: CreateBoardFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateBoardFormData>({
    resolver: zodResolver(createBoardSchema),
    defaultValues: {
      layout: 'grid',
    },
  })

  const createBoardMutation = useCreateBoard()

  const onSubmit = async (data: CreateBoardFormData) => {
    try {
      await createBoardMutation.mutateAsync({
        name: data.name,
        description: data.description,
        imageUrl: '', // TODO: Add image upload functionality
        layout: data.layout as CreateBoardRequestLayoutEnum,
      })

      toast.success('Board created successfully!')
      reset()
      setOpen(false)
    } catch (error) {
      toast.error('Failed to create board. Please try again.')
      console.error('Create board error:', error)
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="my-12">
        <div>
          <label htmlFor="name" className="block text-sm/6 font-medium text-gray-900">
            Name
          </label>
          <div className="mt-2">
            <input
              id="name"
              type="text"
              disabled={createBoardMutation.isPending}
              className={`block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-gray-300 sm:text-sm/6 ${errors.name ? 'outline-red-500 focus:outline-red-500' : ''}`}
              placeholder="Enter board name"
              {...register('name')}
            />
          </div>
          {errors.name && (
            <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
          )}
        </div>
        <div className="mt-4">
          <label htmlFor="description" className="block text-sm/6 font-medium text-gray-900">
            Description
          </label>
          <div className="mt-2">
            <input
              id="description"
              type="text"
              disabled={createBoardMutation.isPending}
              className={`block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-gray-300 sm:text-sm/6 ${errors.description ? 'outline-red-500 focus:outline-red-500' : ''}`}
              placeholder="Enter board description"
              {...register('description')}
            />
          </div>
          {errors.description && (
            <p className="mt-1 text-sm text-red-500">{errors.description.message}</p>
          )}
        </div>
        <div className="mt-4">
          <label htmlFor="layout" className="block text-sm/6 font-medium text-gray-900">
            Layout
          </label>
          <div className="mt-2">
            <select
              id="layout"
              disabled={createBoardMutation.isPending}
              className={`block w-full rounded-md bg-white h-8 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-gray-300 sm:text-sm/6 ${errors.layout ? 'outline-red-500 focus:outline-red-500' : ''}`}
              {...register('layout')}
            >
              <option value="grid">Grid</option>
            </select>
          </div>
          {errors.layout && (
            <p className="mt-1 text-sm text-red-500">{errors.layout.message}</p>
          )}
        </div>
      </div>
      <div className="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
        <button
          type="button"
          onClick={() => setOpen(false)}
          disabled={createBoardMutation.isPending}
          className="inline-flex w-full justify-center rounded-md bg-error border-2 border-error px-3 py-2 text-sm font-semibold text-error shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={createBoardMutation.isPending}
          className="inline-flex w-full justify-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:border-accent disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {createBoardMutation.isPending ? 'Creating...' : 'Create board'}
        </button>
      </div>
    </form>
  )
}