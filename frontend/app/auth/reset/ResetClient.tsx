'use client'

import { Logo } from "@/components/Logo"
import Link from 'next/link'
import { usePublicPageRedirect } from '../../../hooks/useAuthRedirect'

export default function ResetClient() {
  // Redirect authenticated users to dashboard
  usePublicPageRedirect()

  return (
    <div className="bg-base flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm flex flex-col items-center">
        <Link href="/">
          <Logo />
        </Link>
        <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-highlight">
          Reset your password
        </h2>
        <p className="mt-2 text-center text-sm text-base">
          Enter your new password below.
        </p>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form action="#" method="POST" className="space-y-6">
          <div>
            <label htmlFor="password" className="block text-sm/6 font-medium text-highlight">
              New Password
            </label>
            <div className="mt-2">
              <input
                id="password"
                name="password"
                type="password"
                required
                autoComplete="new-password"
                className="block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6"
              />
            </div>
          </div>

          <div>
            <label htmlFor="confirmPassword" className="block text-sm/6 font-medium text-highlight">
              Confirm New Password
            </label>
            <div className="mt-2">
              <input
                id="confirmPassword"
                name="confirmPassword"
                type="password"
                required
                autoComplete="new-password"
                className="block w-full rounded-md bg-inverted px-3 py-1.5 text-inverted outline-1 -outline-offset-1 outline-muted placeholder:text-muted focus:outline-2 focus:-outline-offset-2 focus:outline-accent sm:text-sm/6"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              className="flex w-full justify-center rounded-md bg-accent px-3 py-1.5 text-sm/6 font-semibold text-inverted shadow-xs hover:bg-accent focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            >
              Reset password
            </button>
          </div>
        </form>

        <p className="mt-10 text-center text-sm/6 text-base">
          Remember your password?{' '}
          <Link href="/auth/signin" className="font-semibold text-accent hover:text-accent">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
