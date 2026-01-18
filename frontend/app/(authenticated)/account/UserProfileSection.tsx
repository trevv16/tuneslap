'use client'

import { useAuthContext } from "@/contexts/AuthContext";
import { useUpdateMe } from "@/hooks/users";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import toast from "react-hot-toast";
import z from "zod";

// Validation schema
const updateProfileSchema = z.object({
  name: z.string().min(1, 'Name is required').min(3, 'Name must be at least 3 characters'),
})

type UpdateProfileFormData = z.infer<typeof updateProfileSchema>

export default function UserProfileSection() {
  const authContext = useAuthContext();
  const user = authContext?.user;
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

  // Reset form when user data becomes available
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
        imageUrl: '', // Set to empty string for now, will be updated with GCS upload later
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

  return (
    <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
      <div>
        <h2 className="text-base/7 font-semibold text-base">Personal Information</h2>
        <p className="mt-1 text-sm/6 text-base">Update your account.</p>
      </div>

      {user && (
        <form onSubmit={handleSubmit(onSubmit)} className="md:col-span-2">
          <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
            <div className="col-span-full flex items-center gap-x-8">
              <img
                alt=""
                src={(user?.imageUrl && user?.imageUrl !== "") ? user?.imageUrl : "/defaultUser.jpg"}
                className="size-24 flex-none rounded-lg object-cover"
              />
              <div>
                <button
                  type="button"
                  className="rounded-md bg-white/10 px-3 py-2 text-sm font-semibold text-base shadow-xs hover:bg-white/20"
                >
                  Change avatar
                </button>
                <p className="mt-2 text-xs/5 text-gray-400">JPG, GIF or PNG. 1MB max.</p>
              </div>
            </div>

            <div className="col-span-full">
              <label htmlFor="name" className="block text-sm/6 font-medium text-base">
                Name
              </label>
              <div className="mt-2">
                <input
                  id="name"
                  type="text"
                  disabled={updateMeMutation.isPending}
                  className={`block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6 ${errors.name ? 'outline-red-500 focus:outline-red-500' : ''}`}
                  placeholder="Enter your name"
                  {...register('name')}
                />
              </div>
              {errors.name && (
                <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>

            <div className="col-span-full">
              <label htmlFor="email" className="block text-sm/6 font-medium text-base">
                Email address
              </label>
              <div className="mt-2">
                <input
                  id="email"
                  name="email"
                  type="email"
                  value={user?.email}
                  autoComplete="email"
                  disabled
                  className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6"
                />
              </div>
            </div>
          </div>

          <div className="mt-8 flex">
            <button
              type="submit"
              disabled={updateMeMutation.isPending}
              className="rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight hover:text-inverted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {updateMeMutation.isPending ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      )}
    </div>
  )
}