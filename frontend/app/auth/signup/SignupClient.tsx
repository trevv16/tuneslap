'use client'

import { useSignUp } from '@/hooks/auth'
import { zodResolver } from '@hookform/resolvers/zod'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { toast } from 'react-hot-toast'
import { z } from 'zod'
import { usePublicPageRedirect } from '../../../hooks/useAuthRedirect'

const signupSchema = z.object({
  name: z.string().min(3, 'Name must be at least 3 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters')
})

type SignupFormValues = z.infer<typeof signupSchema>

export default function SignupClient() {
  // Redirect authenticated users to dashboard
  usePublicPageRedirect()

  const router = useRouter()
  const { mutate: signup, isPending } = useSignUp()

  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<SignupFormValues>({
    resolver: zodResolver(signupSchema)
  })

  const onSubmit = (data: SignupFormValues) => {
    signup(data, {
      onSuccess: () => {
        toast.success('Account created successfully! Please sign in.')
        router.push('/auth/signin')
      },
      onError: (error) => {
        toast.error(error.message || 'Failed to create account')
      }
    })
  }

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
          Create your account
        </h2>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div>
            <label htmlFor="name" className="block text-sm/6 font-medium text-highlight">
              Name
            </label>
            <div className="mt-2">
              <input
                id="name"
                type="text"
                autoComplete="name"
                className="block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6"
                {...register('name')}
              />
              {errors.name && (
                <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>
          </div>

          <div>
            <label htmlFor="email" className="block text-sm/6 font-medium text-highlight">
              Email address
            </label>
            <div className="mt-2">
              <input
                id="email"
                type="email"
                autoComplete="email"
                className="block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6"
                {...register('email')}
              />
              {errors.email && (
                <p className="mt-1 text-sm text-red-500">{errors.email.message}</p>
              )}
            </div>
          </div>

          <div>
            <div className="flex items-center justify-between">
              <label htmlFor="password" className="block text-sm/6 font-medium text-highlight">
                Password
              </label>
            </div>
            <div className="mt-2">
              <input
                id="password"
                type="password"
                autoComplete="new-password"
                className="block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6"
                {...register('password')}
              />
              {errors.password && (
                <p className="mt-1 text-sm text-red-500">{errors.password.message}</p>
              )}
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isPending}
              className="flex w-full justify-center rounded-md bg-accent px-3 py-1.5 text-sm/6 font-semibold text-inverted shadow-xs hover:bg-accent focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isPending ? 'Creating account...' : 'Sign up'}
            </button>
          </div>
        </form>

        <p className="mt-10 text-center text-sm/6 text-base">
          Already have an account?{' '}
          <Link href="/auth/signin" className="font-semibold text-accent hover:text-accent">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
