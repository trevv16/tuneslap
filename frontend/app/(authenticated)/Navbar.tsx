'use client'

import { useAuthContext } from "@/contexts/AuthContext";
import { classNames } from "@/utils/helpers";
import { Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { Bars3Icon, MagnifyingGlassIcon, XMarkIcon } from "@heroicons/react/20/solid";
import Link from "next/link";
import { usePathname } from "next/navigation";


export default function Navbar() {
  const pathname = usePathname();
  const { user, signOut } = useAuthContext();
  const navigation = [
    { name: 'Dashboard', href: '/dashboard', current: pathname === '/dashboard' },
    { name: 'Library', href: '/library', current: pathname === '/library' },
  ];
  const userNavigation = [
    { name: 'Your Account', href: '/account' },
  ];

  return (
    <Disclosure as="nav" className="bg-base border-b border-white lg:border-none">
      <div className="mx-auto max-w-7xl px-2 sm:px-4 lg:px-8">
        <div className="relative flex h-16 items-center justify-between lg:border-b lg:border-white">
          <div className="flex items-center px-2 lg:px-0">
            <div className="shrink-0">
              <Link href="/dashboard">
                <img
                  alt="TuneSlap"
                  src="/logo.png"
                  className="h-8 w-auto"
                />
              </Link>
            </div>
            <div className="hidden lg:ml-10 lg:block">
              <div className="flex space-x-4">
                {navigation.map((item) => (
                  <Link
                    key={item.name}
                    href={item.href}
                    aria-current={item.current ? 'page' : undefined}
                    className={classNames(
                      item.current ? 'bg-primary-700 text-white' : 'text-white hover:bg-primary-500/75',
                      'rounded-md px-3 py-2 text-sm font-medium',
                    )}
                  >
                    {item.name}
                  </Link>
                ))}
              </div>
            </div>
          </div>
          <div className="flex flex-1 justify-center px-2 lg:ml-6 lg:justify-end">
            <div className="grid w-full max-w-lg grid-cols-1 lg:max-w-xs">
              <input
                name="search"
                type="search"
                placeholder="Search"
                aria-label="Search"
                className="col-start-1 row-start-1 block w-full rounded-md bg-muted py-1.5 pr-3 pl-10 text-base outline-hidden placeholder:text-base placeholder:text-sm focus:outline-2 focus:outline-offset-2 focus:border-accent sm:text-sm/6"
              />
              <MagnifyingGlassIcon
                aria-hidden="true"
                className="pointer-events-none col-start-1 row-start-1 ml-3 size-5 self-center text-gray-400"
              />
            </div>
          </div>
          <div className="flex lg:hidden">
            {/* Mobile menu button */}
            <DisclosureButton className="group relative inline-flex items-center justify-center rounded-md bg-primary-600 p-2 text-primary-200 hover:bg-primary-500/75 hover:text-white focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-primary-600 focus:outline-hidden">
              <span className="absolute -inset-0.5" />
              <span className="sr-only">Open main menu</span>
              <Bars3Icon aria-hidden="true" className="block size-6 group-data-open:hidden" />
              <XMarkIcon aria-hidden="true" className="hidden size-6 group-data-open:block" />
            </DisclosureButton>
          </div>
          <div className="hidden lg:ml-4 lg:block">
            <div className="flex items-center">
              {/* Profile dropdown */}
              <Menu as="div" className="relative ml-3 shrink-0">
                <div>
                  <MenuButton className="relative flex rounded-full bg-primary-600 text-sm text-white focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-primary-600 focus:outline-hidden">
                    <span className="absolute -inset-1.5" />
                    <span className="sr-only">Open user menu</span>
                    <img alt="" src={(user?.imageUrl && user?.imageUrl !== "") ? user?.imageUrl : "/defaultUser.jpg"} className="size-8 rounded-full" />
                  </MenuButton>
                </div>
                <MenuItems
                  transition
                  className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-highlight py-1 shadow-lg ring-1 ring-black/5 transition focus:outline-hidden data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-100 data-enter:ease-out data-leave:duration-75 data-leave:ease-in"
                >
                  {userNavigation.map((item) => (
                    <MenuItem key={item.name}>
                      <Link
                        href={item.href}
                        className="block px-4 py-2 text-sm text-base data-focus:bg-gray-100 data-focus:outline-hidden"
                      >
                        {item.name}
                      </Link>
                    </MenuItem>
                  ))}
                  <MenuItem key="signout">
                    <button
                      onClick={signOut}
                      className="block w-full text-left px-4 py-2 text-sm text-base data-focus:bg-gray-100 data-focus:outline-hidden"
                    >
                      Sign out
                    </button>
                  </MenuItem>
                </MenuItems>
              </Menu>
            </div>
          </div>
        </div>
      </div>

      <DisclosurePanel className="lg:hidden">
        <div className="space-y-1 px-2 pt-2 pb-3">
          {navigation.map((item) => (
            <DisclosureButton
              key={item.name}
              as="a"
              href={item.href}
              aria-current={item.current ? 'page' : undefined}
              className={classNames(
                item.current ? 'bg-primary-700 text-white' : 'text-white hover:bg-primary-500/75',
                'block rounded-md px-3 py-2 text-base font-medium',
              )}
            >
              {item.name}
            </DisclosureButton>
          ))}
        </div>
        <div className="border-t border-primary-700 pt-4 pb-3">
          <div className="flex items-center px-5">
            <div className="shrink-0">
              <img alt="" src={(user?.imageUrl && user?.imageUrl !== "") ? user?.imageUrl : "/defaultUser.jpg"} className="size-10 rounded-full" />
            </div>
            <div className="ml-3">
              <div className="text-base font-medium text-white">{user?.name}</div>
              <div className="text-sm font-medium text-primary-300">{user?.email}</div>
            </div>
          </div>
          <div className="mt-3 space-y-1 px-2">
            {userNavigation.map((item) => (
              <DisclosureButton
                key={item.name}
                as="a"
                href={item.href}
                className="block rounded-md px-3 py-2 text-base font-medium text-white hover:bg-primary-500/75"
              >
                {item.name}
              </DisclosureButton>
            ))}
            <DisclosureButton
              key="signout"
              onClick={signOut}
              className="block w-full text-left rounded-md px-3 py-2 text-base font-medium text-white hover:bg-primary-500/75"
            >
              Sign out
            </DisclosureButton>
          </div>
        </div>
      </DisclosurePanel>
    </Disclosure>
  )
}