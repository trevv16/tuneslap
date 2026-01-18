"use client"

import Link from "next/link"

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen px-4 py-12">
      <h2 className="text-2xl font-bold mb-4">404 - Page Not Found</h2>
      <p className="mb-4">Sorry, we couldn&apos;t find the page you&apos;re looking for.</p>
      <Link href="/" className="px-4 py-2 bg-primary-500 text-light-50 rounded hover:bg-primary-600">
        Go back home
      </Link>
    </div>
  )
}
