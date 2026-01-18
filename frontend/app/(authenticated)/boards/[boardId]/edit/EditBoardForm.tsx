'use client'

import type { BoardResponse as Board } from '@/api/models'
import { UpdateBoardRequestLayoutEnum } from '@/api/models'
import { useUpdateBoard } from '@/hooks/boards'
import { zodResolver } from '@hookform/resolvers/zod'
import { UseQueryResult } from '@tanstack/react-query'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { z } from 'zod'

// Validation schema
const updateBoardSchema = z.object({
  name: z.string().min(1, 'Name is required').min(3, 'Name must be at least 3 characters'),
  description: z.string().min(1, 'Description is required'),
  layout: z.string().min(1, 'Layout is required'),
})

type UpdateBoardFormData = z.infer<typeof updateBoardSchema>

type EditBoardFormProps = {
  boardId: string;
  boardQuery: UseQueryResult<Board, Error>;
}

export default function EditBoardForm({ boardQuery, boardId }: EditBoardFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateBoardFormData>({
    resolver: zodResolver(updateBoardSchema),
    defaultValues: {
      name: '',
      description: '',
      layout: 'grid',
    },
  })

  const { data: board } = boardQuery;
  const updateBoardMutation = useUpdateBoard()

  // Reset form when board data becomes available
  useEffect(() => {
    if (board) {
      reset({
        name: board.name || '',
        description: board.description || '',
        layout: (board.layout === 'grid' || board.layout === 'list' 
          ? board.layout 
          : 'grid') as UpdateBoardRequestLayoutEnum,
      })
    }
  }, [board, reset])

  const onSubmit = async (data: UpdateBoardFormData) => {
    try {
      await updateBoardMutation.mutateAsync({
        boardId,
        name: data.name,
        description: data.description,
        layout: data.layout as UpdateBoardRequestLayoutEnum,
        imageUrl: '', // Set to empty string for now, will be updated with GCS upload later
      })

      toast.success('Board updated successfully!')
    } catch (error) {
      toast.error('Failed to update board. Please try again.')
      console.error('Update board error:', error)
    }
  }

  // if (isLoadingBoard) {
  //   return (
  //     <div className="divide-y divide-white/5">
  //       <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
  //         <div>
  //           <h2 className="text-base/7 font-semibold text-base">Board Information</h2>
  //           <p className="mt-1 text-sm/6 text-base">Loading board information...</p>
  //         </div>
  //       </div>
  //     </div>
  //   )
  // }

  return (
    <div className="divide-y divide-white/5">
      <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
        <div>
          <h2 className="text-base/7 font-semibold text-base">Board Information</h2>
          <p className="mt-1 text-sm/6 text-base">Update your board details and settings.</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="md:col-span-2">
          <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
            <div className="col-span-full flex items-center gap-x-8">
              <div className="size-24 flex-none rounded-lg bg-white/5 flex items-center justify-center">
                <svg
                  className="h-12 w-12 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                  />
                </svg>
              </div>
              <div>
                <button
                  type="button"
                  className="rounded-md bg-white/10 px-3 py-2 text-sm font-semibold text-base shadow-xs hover:bg-white/20"
                >
                  Change Image
                </button>
                <p className="mt-2 text-xs/5 text-gray-400">JPG, GIF or PNG. 1MB max.</p>
              </div>
            </div>

            <div className="col-span-full">
              <label htmlFor="name" className="block text-sm/6 font-medium text-base">
                Board Name
              </label>
              <div className="mt-2">
                <input
                  id="name"
                  type="text"
                  disabled={updateBoardMutation.isPending}
                  className={`block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6 ${errors.name ? 'outline-red-500 focus:outline-red-500' : ''}`}
                  placeholder="Enter board name"
                  {...register('name')}
                />
              </div>
              {errors.name && (
                <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>

            <div className="col-span-full">
              <label htmlFor="description" className="block text-sm/6 font-medium text-base">
                Description
              </label>
              <div className="mt-2">
                <textarea
                  id="description"
                  rows={3}
                  disabled={updateBoardMutation.isPending}
                  className={`block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6 ${errors.description ? 'outline-red-500 focus:outline-red-500' : ''}`}
                  placeholder="Enter board description"
                  {...register('description')}
                />
              </div>
              {errors.description && (
                <p className="mt-1 text-sm text-red-500">{errors.description.message}</p>
              )}
            </div>

            <div className="col-span-full">
              <label htmlFor="layout" className="block text-sm/6 font-medium text-base">
                Layout
              </label>
              <div className="mt-2">
                <select
                  id="layout"
                  disabled={updateBoardMutation.isPending}
                  className={`block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6 ${errors.layout ? 'outline-red-500 focus:outline-red-500' : ''}`}
                  {...register('layout')}
                >
                  <option value="grid">Grid</option>
                  {/* <option value="list">List</option> */}
                </select>
              </div>
              {errors.layout && (
                <p className="mt-1 text-sm text-red-500">{errors.layout.message}</p>
              )}
            </div>
          </div>

          <div className="mt-8 flex">
            <button
              type="submit"
              disabled={updateBoardMutation.isPending}
              className="rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight hover:text-inverted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {updateBoardMutation.isPending ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
