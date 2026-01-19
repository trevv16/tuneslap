'use client'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useSignUp } from '@/hooks/auth'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'
import { usePublicPageRedirect } from '../../../hooks/useAuthRedirect'

const signupSchema = z.object({
  name: z.string().min(3, 'Name must be at least 3 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters')
})

type SignupFormValues = z.infer<typeof signupSchema>

export default function SignupClient() {
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
          Create your account
        </h2>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div>
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              type="text"
              autoComplete="name"
              className="mt-2"
              {...register('name')}
            />
            {errors.name && (
              <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          <div>
            <Label htmlFor="email">Email address</Label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              className="mt-2"
              {...register('email')}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>

          <div>
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              className="mt-2"
              {...register('password')}
            />
            {errors.password && (
              <p className="mt-1 text-sm text-destructive">{errors.password.message}</p>
            )}
          </div>

          <Button
            type="submit"
            disabled={isPending}
            className="w-full"
          >
            {isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Creating account...
              </>
            ) : (
              'Sign up'
            )}
          </Button>
        </form>

        <p className="mt-10 text-center text-sm text-muted-foreground">
          Already have an account?{' '}
          <Link href="/auth/signin" className="font-semibold text-primary hover:text-primary/80">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
