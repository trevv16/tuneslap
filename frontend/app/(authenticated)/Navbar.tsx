'use client'

import { Logo } from "@/components/Logo"
import { ThemeToggle } from "@/components/ThemeToggle"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { useAuthContext } from "@/contexts/AuthContext"
import { cn } from "@/lib/utils"
import { Menu, User } from "lucide-react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { useState } from "react"

export default function Navbar() {
  const pathname = usePathname()
  const { user, signOut } = useAuthContext()
  const [sheetOpen, setSheetOpen] = useState(false)

  const navigation = [
    { name: 'Dashboard', href: '/dashboard' },
    { name: 'Library', href: '/library' },
  ]

  const isActive = (href: string) => pathname === href

  const getUserInitials = () => {
    if (!user?.name) return 'U'
    const parts = user.name.split(' ')
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase()
    }
    return parts[0][0].toUpperCase()
  }

  return (
    <nav className="border-b bg-background">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo and Desktop Nav */}
          <div className="flex items-center gap-8">
            <Link href="/dashboard" className="shrink-0">
              <Logo />
            </Link>

            {/* Desktop Navigation */}
            <div className="hidden md:flex md:items-center md:gap-1">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={cn(
                    "px-3 py-2 text-sm font-medium rounded-md transition-colors",
                    isActive(item.href)
                      ? "bg-accent text-accent-foreground"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  )}
                >
                  {item.name}
                </Link>
              ))}
            </div>
          </div>

          {/* Desktop Right Side */}
          <div className="hidden md:flex md:items-center md:gap-2">
            <ThemeToggle />

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-9 w-9 rounded-full" data-testid="user-menu-trigger">
                  <Avatar className="h-9 w-9">
                    <AvatarImage
                      src={user?.imageUrl && user.imageUrl !== "" ? user.imageUrl : undefined}
                      alt={user?.name || "User"}
                    />
                    <AvatarFallback>
                      {getUserInitials()}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                {(user?.name || user?.email) && (
                  <>
                    <div className="px-2 py-1.5">
                      {user?.name && (
                        <p className="text-sm font-medium">{user.name}</p>
                      )}
                      {user?.email && (
                        <p className="text-xs text-muted-foreground truncate">
                          {user.email}
                        </p>
                      )}
                    </div>
                    <DropdownMenuSeparator />
                  </>
                )}
                <DropdownMenuItem asChild>
                  <Link href="/account" className="cursor-pointer">
                    <User className="mr-2 h-4 w-4" />
                    Account
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={signOut}
                  className="cursor-pointer text-destructive focus:text-destructive"
                  data-testid="sign-out-button"
                >
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {/* Mobile Menu Button */}
          <div className="flex items-center gap-2 md:hidden">
            <ThemeToggle />

            <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon">
                  <Menu className="h-5 w-5" />
                  <span className="sr-only">Open menu</span>
                </Button>
              </SheetTrigger>
              <SheetContent side="right" className="w-[300px] sm:w-[350px]">
                <SheetHeader>
                  <SheetTitle className="sr-only">Navigation Menu</SheetTitle>
                </SheetHeader>

                <div className="px-4">
                  {/* User Info */}
                  <div className="flex items-center gap-3 py-4 border-b">
                    <Avatar className="h-10 w-10">
                      <AvatarImage
                        src={user?.imageUrl && user.imageUrl !== "" ? user.imageUrl : undefined}
                        alt={user?.name || "User"}
                      />
                      <AvatarFallback>
                        {getUserInitials()}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex flex-col">
                      {user?.name && (
                        <p className="font-medium">{user.name}</p>
                      )}
                      {user?.email && (
                        <p className="text-sm text-muted-foreground truncate max-w-[200px]">
                          {user.email}
                        </p>
                      )}
                    </div>
                  </div>

                  {/* Navigation Links */}
                  <div className="flex flex-col gap-1 py-4">
                    {navigation.map((item) => (
                      <Link
                        key={item.name}
                        href={item.href}
                        onClick={() => { setSheetOpen(false); }}
                        className={cn(
                          "px-3 py-2 text-sm font-medium rounded-md transition-colors",
                          isActive(item.href)
                            ? "bg-accent text-accent-foreground"
                            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                        )}
                      >
                        {item.name}
                      </Link>
                    ))}
                    <Link
                      href="/account"
                      onClick={() => { setSheetOpen(false); }}
                      className={cn(
                        "px-3 py-2 text-sm font-medium rounded-md transition-colors",
                        isActive('/account')
                          ? "bg-accent text-accent-foreground"
                          : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                      )}
                    >
                      Account
                    </Link>
                  </div>

                  {/* Sign Out */}
                  <div className="border-t pt-4">
                    <Button
                      variant="ghost"
                      className="w-full justify-start text-destructive hover:text-destructive hover:bg-destructive/10"
                      onClick={() => {
                        setSheetOpen(false)
                        signOut()
                      }}
                      data-testid="sign-out-button-mobile"
                    >
                      Sign out
                    </Button>
                  </div>
                </div>
              </SheetContent>
            </Sheet>
          </div>
        </div>
      </div>
    </nav>
  )
}
