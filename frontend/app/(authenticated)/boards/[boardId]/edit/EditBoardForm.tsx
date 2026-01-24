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

interface EditBoardFormProps {
  boardId: string
  boardQuery: UseQueryResult<Board>
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
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="sm:col-span-2">
          <Label htmlFor="name">Board Name</Label>
          <Input
            id="name"
            disabled={updateBoardMutation.isPending}
            placeholder="Enter board name"
            className="mt-1.5"
            {...register('name')}
          />
          {errors.name && (
            <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>
          )}
        </div>

        <div className="sm:col-span-2">
          <Label htmlFor="description">Description</Label>
          <Textarea
            id="description"
            rows={3}
            disabled={updateBoardMutation.isPending}
            placeholder="Enter board description"
            className="mt-1.5"
            {...register('description')}
          />
          {errors.description && (
            <p className="mt-1 text-sm text-destructive">{errors.description.message}</p>
          )}
        </div>

        <div>
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
                <SelectTrigger className="mt-1.5">
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

      <div className="flex justify-end pt-2">
        <Button
          type="submit"
          disabled={updateBoardMutation.isPending}
        >
          {updateBoardMutation.isPending ? 'Saving...' : 'Save Changes'}
        </Button>
      </div>
    </form>
  )
}
