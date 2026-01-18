"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

export default function AccountMenu() {
  const pathname = usePathname();
  const secondaryNavigation = [
    { name: 'Account', href: '/account', current: pathname === '/account' },
  ]

  return (
    <header className="border-b border-white/5">
      {/* Secondary navigation */}
      <nav className="flex overflow-x-auto py-4">
        <ul
          role="list"
          className="flex min-w-full flex-none gap-x-6 px-4 text-sm/6 font-semibold text-base sm:px-6 lg:px-8"
        >
          {secondaryNavigation.map((item) => (
            <li key={item.name}>
              <Link href={item.href} className={item.current ? 'text-accent' : ''}>
                {item.name}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </header>
  )
}