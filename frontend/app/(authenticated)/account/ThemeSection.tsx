'use client'

import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"

export default function ThemeSection() {
  const { theme, setTheme } = useTheme()

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : 'light')
  }

  return (
    <div data-testid="theme-section" className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
      <div>
        <h2 className="text-base font-semibold text-foreground">Theme Settings</h2>
        <p className="mt-1 text-sm text-muted-foreground">Choose your preferred theme for the application.</p>
      </div>

      <div className="md:col-span-2">
        <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
          <div className="col-span-full">
            <Label>Theme</Label>
            <div className="mt-2">
              <div className="flex items-center gap-x-4">
                <Button
                  variant="outline"
                  onClick={toggleTheme}
                  className="gap-3"
                  data-testid="theme-toggle"
                >
                  {theme === 'light' ? (
                    <>
                      <Sun className="h-5 w-5" />
                      Light Mode
                    </>
                  ) : (
                    <>
                      <Moon className="h-5 w-5" />
                      Dark Mode
                    </>
                  )}
                </Button>
                <span className="text-sm text-muted-foreground">
                  Currently using {theme === 'light' ? 'light' : 'dark'} theme
                </span>
              </div>
            </div>
            <p className="mt-2 text-xs text-muted-foreground">
              Switch between light and dark themes to match your preference.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
