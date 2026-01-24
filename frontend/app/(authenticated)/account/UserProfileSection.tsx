'use client'

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuthContext } from "@/contexts/AuthContext"
import { useUpdateMe } from "@/hooks/users"
import { zodResolver } from "@hookform/resolvers/zod"
import { useEffect } from "react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import z from "zod"

const updateProfileSchema = z.object({
  name: z.string().min(1, 'Name is required').min(3, 'Name must be at least 3 characters'),
})

type UpdateProfileFormData = z.infer<typeof updateProfileSchema>

export default function UserProfileSection() {
  const authContext = useAuthContext()
  const user = authContext?.user
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateProfileFormData>({
    resolver: zodResolver(updateProfileSchema),
    defaultValues: {
      name: '',
    },
  })

  useEffect(() => {
    if (user) {
      reset({
        name: user.name || '',
      })
    }
  }, [user, reset])

  const updateMeMutation = useUpdateMe()

  const onSubmit = async (data: UpdateProfileFormData) => {
    try {
      await updateMeMutation.mutateAsync({
        name: data.name,
        imageUrl: '',
      })

      toast.success('Profile updated successfully!')
      reset({
        name: data.name,
      })
    } catch (error) {
      toast.error('Failed to update profile. Please try again.')
      console.error('Update profile error:', error)
    }
  }

  const getUserInitials = () => {
    if (!user?.name) return 'U'
    const parts = user.name.split(' ')
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase()
    }
    return parts[0][0].toUpperCase()
  }

  return (
    <div data-testid="profile-section" className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
      <div>
        <h2 className="text-base font-semibold text-foreground">Personal Information</h2>
        <p className="mt-1 text-sm text-muted-foreground">Update your account.</p>
      </div>

      {user && (
        <form onSubmit={handleSubmit(onSubmit)} className="md:col-span-2">
          <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
            <div className="col-span-full flex items-center gap-x-8">
              <Avatar className="h-24 w-24">
                <AvatarImage src={user.imageUrl && user.imageUrl !== "" ? user.imageUrl : undefined} />
                <AvatarFallback className="text-2xl">{getUserInitials()}</AvatarFallback>
              </Avatar>
              <div>
                <Button type="button" variant="secondary">
                  Change avatar
                </Button>
                <p className="mt-2 text-xs text-muted-foreground">JPG, GIF or PNG. 1MB max.</p>
              </div>
            </div>

            <div className="col-span-full">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                disabled={updateMeMutation.isPending}
                placeholder="Enter your name"
                className="mt-2"
                {...register('name')}
              />
              {errors.name && (
                <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>
              )}
            </div>

            <div className="col-span-full">
              <Label htmlFor="email">Email address</Label>
              <Input
                id="email"
                type="email"
                value={user?.email}
                autoComplete="email"
                disabled
                className="mt-2"
              />
            </div>
          </div>

          <div className="mt-8">
            <Button
              type="submit"
              disabled={updateMeMutation.isPending}
            >
              {updateMeMutation.isPending ? 'Saving...' : 'Save'}
            </Button>
          </div>
        </form>
      )}
    </div>
  )
}
