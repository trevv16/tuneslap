'use client'

import { useCreateBoard } from '@/hooks/boards'
import type { CreateBoardRequestLayoutEnum } from '@/api/models'
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

const createBoardSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  description: z.string().min(1, 'Description is required'),
  layout: z.string().min(1, 'Layout is required'),
})

type CreateBoardFormData = z.infer<typeof createBoardSchema>

type CreateBoardFormProps = {
  setOpen: Dispatch<SetStateAction<boolean>>
}

export default function CreateBoardForm({ setOpen }: CreateBoardFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
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
        imageUrl: '',
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
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="name">Name</Label>
        <Input
          id="name"
          disabled={createBoardMutation.isPending}
          placeholder="Enter board name"
          className="mt-2"
          {...register('name')}
        />
        {errors.name && (
          <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>
        )}
      </div>

      <div>
        <Label htmlFor="description">Description</Label>
        <Input
          id="description"
          disabled={createBoardMutation.isPending}
          placeholder="Enter board description"
          className="mt-2"
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
              disabled={createBoardMutation.isPending}
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

      <div className="flex gap-3 pt-4">
        <Button
          type="button"
          variant="outline"
          onClick={() => setOpen(false)}
          disabled={createBoardMutation.isPending}
          className="flex-1"
        >
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={createBoardMutation.isPending}
          className="flex-1"
        >
          {createBoardMutation.isPending ? 'Creating...' : 'Create board'}
        </Button>
      </div>
    </form>
  )
}
