'use client'

import type { BoardResponse as Board } from '@/api/models'
import { UpdateBoardRequestLayoutEnum } from '@/api/models'
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
import { Textarea } from '@/components/ui/textarea'
import { useUpdateBoard } from '@/hooks/boards'
import { zodResolver } from '@hookform/resolvers/zod'
import { UseQueryResult } from '@tanstack/react-query'
import { LayoutGrid } from 'lucide-react'
import { useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

const updateBoardSchema = z.object({
  name: z.string().min(1, 'Name is required').min(3, 'Name must be at least 3 characters'),
  description: z.string().min(1, 'Description is required'),
  layout: z.string().min(1, 'Layout is required'),
})

type UpdateBoardFormData = z.infer<typeof updateBoardSchema>

type EditBoardFormProps = {
  boardId: string
  boardQuery: UseQueryResult<Board, Error>
}

export default function EditBoardForm({ boardQuery, boardId }: EditBoardFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
  } = useForm<UpdateBoardFormData>({
    resolver: zodResolver(updateBoardSchema),
    defaultValues: {
      name: '',
      description: '',
      layout: 'grid',
    },
  })

  const { data: board } = boardQuery
  const updateBoardMutation = useUpdateBoard()

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
        imageUrl: '',
      })

      toast.success('Board updated successfully!')
    } catch (error) {
      toast.error('Failed to update board. Please try again.')
      console.error('Update board error:', error)
    }
  }

  return (
    <div className="divide-y divide-border">
      <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
        <div>
          <h2 className="text-base font-semibold text-foreground">Board Information</h2>
          <p className="mt-1 text-sm text-muted-foreground">Update your board details and settings.</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="md:col-span-2">
          <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
            <div className="col-span-full flex items-center gap-x-8">
              <div className="h-24 w-24 flex-none rounded-lg bg-muted flex items-center justify-center">
                <LayoutGrid className="h-12 w-12 text-muted-foreground" />
              </div>
              <div>
                <Button type="button" variant="secondary">
                  Change Image
                </Button>
                <p className="mt-2 text-xs text-muted-foreground">JPG, GIF or PNG. 1MB max.</p>
              </div>
            </div>

            <div className="col-span-full">
              <Label htmlFor="name">Board Name</Label>
              <Input
                id="name"
                disabled={updateBoardMutation.isPending}
                placeholder="Enter board name"
                className="mt-2"
                {...register('name')}
              />
              {errors.name && (
                <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>
              )}
            </div>

            <div className="col-span-full">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                rows={3}
                disabled={updateBoardMutation.isPending}
                placeholder="Enter board description"
                className="mt-2"
                {...register('description')}
              />
              {errors.description && (
                <p className="mt-1 text-sm text-destructive">{errors.description.message}</p>
              )}
            </div>

            <div className="col-span-full">
              <Label htmlFor="layout">Layout</Label>
              <Controller
                name="layout"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value}
                    onValueChange={field.onChange}
                    disabled={updateBoardMutation.isPending}
                  >
                    <SelectTrigger className="mt-2">
                      <SelectValue placeholder="Select layout" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="grid">Grid</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.layout && (
                <p className="mt-1 text-sm text-destructive">{errors.layout.message}</p>
              )}
            </div>
          </div>

          <div className="mt-8">
            <Button
              type="submit"
              disabled={updateBoardMutation.isPending}
            >
              {updateBoardMutation.isPending ? 'Saving...' : 'Save'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
