'use client'

import { useGetMe } from '@/hooks/users'
import { zodResolver } from '@hookform/resolvers/zod'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { z } from 'zod'
import { useSignIn } from '../../../hooks/auth'
import { usePublicPageRedirect } from '../../../hooks/useAuthRedirect'

// Validation schema
const signInSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(1, 'Password is required'),
})

type SignInFormData = z.infer<typeof signInSchema>

export default function SignInClient() {
  const router = useRouter()
  const signInMutation = useSignIn()
  useGetMe(signInMutation.data?.data?.token || "");

  // Redirect authenticated users to dashboard
  usePublicPageRedirect()

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<SignInFormData>({
    resolver: zodResolver(signInSchema),
  })

  const onSubmit = async (data: SignInFormData) => {
    try {
      await signInMutation.mutateAsync(data, {
        onSuccess: () => {
          // Toast is handled by hook
          reset()
          // Router push is handled by usePublicPageRedirect implicitly when auth state changes
          // But to be explicit and immediate:
          router.push('/dashboard')
        }
      })
    } catch {
      toast.error('Sign in failed. Please try again.')
    }
  }

  // Handle success and error states
  useEffect(() => {
    if (signInMutation.isSuccess) {
      toast.success('Sign in successful!')
    }
  }, [signInMutation.isSuccess])

  return (
    <div className="bg-base flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <Link href="/">
          <img
            alt="TuneSlap"
            src="/logo.png"
            className="mx-auto h-10 w-auto"
          />
        </Link>
        <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-highlight">
          Sign in to your account
        </h2>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div>
            <label htmlFor="email" className="block text-sm/6 font-medium text-highlight">
              Email address
            </label>
            <div className="mt-2">
              <input
                id="email"
                type="email"
                autoComplete="email"
                disabled={signInMutation.isPending}
                className={`block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6 ${errors.email ? 'outline-red-500 focus:outline-red-500' : ''
                  }`}
                placeholder="Enter your email"
                {...register('email')}
              />
            </div>
            {errors.email && (
              <p className="mt-1 text-sm text-red-500">{errors.email.message}</p>
            )}
          </div>

          <div>
            <div className="flex items-center justify-between">
              <label htmlFor="password" className="block text-sm/6 font-medium text-highlight">
                Password
              </label>
              <div className="text-sm">
                <Link href="/auth/forgot" className="font-semibold text-accent hover:text-accent">
                  Forgot password?
                </Link>
              </div>
            </div>
            <div className="mt-2">
              <input
                id="password"
                type="password"
                autoComplete="current-password"
                disabled={signInMutation.isPending}
                className={`block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6 ${errors.password ? 'outline-red-500 focus:outline-red-500' : ''
                  }`}
                placeholder="Enter your password"
                {...register('password')}
              />
            </div>
            {errors.password && (
              <p className="mt-1 text-sm text-red-500">{errors.password.message}</p>
            )}
          </div>

          <div>
            <button
              type="submit"
              disabled={signInMutation.isPending}
              className="flex w-full justify-center rounded-md bg-accent px-3 py-1.5 text-sm/6 font-semibold text-inverted shadow-xs hover:bg-accent focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {signInMutation.isPending ? (
                <div className="flex items-center">
                  <svg
                    className="animate-spin -ml-1 mr-3 h-4 w-4 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Signing in...
                </div>
              ) : (
                'Sign in'
              )}
            </button>
          </div>
        </form>

        <p className="mt-10 text-center text-sm/6 text-base">
          Not a member?{' '}
          <Link href="/auth/signup" className="font-semibold text-accent hover:text-accent">
            Start a 14 day free trial
          </Link>
        </p>
      </div>
    </div>
  )
}
