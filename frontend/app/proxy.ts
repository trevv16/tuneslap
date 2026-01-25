import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function proxy(request: NextRequest) {
  const isDemoMode = process.env.NEXT_PUBLIC_DEMO_MODE === 'true';
  const pathname = request.nextUrl.pathname;

  // If not in demo mode and accessing the homepage, redirect to sign in
  if (!isDemoMode && pathname === '/') {
    return NextResponse.redirect(new URL('/auth/signin', request.url));
  }

  return NextResponse.next();
}

export const config = {
  // Only run middleware on the homepage
  matcher: '/',
};
