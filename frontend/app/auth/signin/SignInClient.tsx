'use client'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useGetMe } from '@/hooks/users'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'
import { useSignIn } from '../../../hooks/auth'
import { usePublicPageRedirect } from '../../../hooks/useAuthRedirect'

const signInSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(1, 'Password is required'),
})

type SignInFormData = z.infer<typeof signInSchema>

export default function SignInClient() {
  const router = useRouter()
  const signInMutation = useSignIn()
  useGetMe(signInMutation.data?.data?.token || "")

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
          reset()
          router.push('/dashboard')
        }
      })
    } catch {
      toast.error('Sign in failed. Please try again.')
    }
  }

  useEffect(() => {
    if (signInMutation.isSuccess) {
      toast.success('Sign in successful!')
    }
  }, [signInMutation.isSuccess])

  return (
    <div className="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <Link href="/">
          <img
            alt="TuneSlap"
            src="/logo.png"
            className="mx-auto h-10 w-auto"
          />
        </Link>
        <h2 className="mt-10 text-center text-2xl font-bold tracking-tight text-foreground">
          Sign in to your account
        </h2>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div>
            <Label htmlFor="email">Email address</Label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              disabled={signInMutation.isPending}
              placeholder="Enter your email"
              className="mt-2"
              {...register('email')}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>

          <div>
            <div className="flex items-center justify-between">
              <Label htmlFor="password">Password</Label>
              <div className="text-sm">
                <Link href="/auth/forgot" className="font-semibold text-primary hover:text-primary/80">
                  Forgot password?
                </Link>
              </div>
            </div>
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              disabled={signInMutation.isPending}
              placeholder="Enter your password"
              className="mt-2"
              {...register('password')}
            />
            {errors.password && (
              <p className="mt-1 text-sm text-destructive">{errors.password.message}</p>
            )}
          </div>

          <Button
            type="submit"
            disabled={signInMutation.isPending}
            className="w-full"
          >
            {signInMutation.isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Signing in...
              </>
            ) : (
              'Sign in'
            )}
          </Button>
        </form>

        <p className="mt-10 text-center text-sm text-muted-foreground">
          Not a member?{' '}
          <Link href="/auth/signup" className="font-semibold text-primary hover:text-primary/80">
            Start a 14 day free trial
          </Link>
        </p>
      </div>
    </div>
  )
}
