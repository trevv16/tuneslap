'use client'

import { useThemeContext } from "@/contexts/ThemeContext";
import { MoonIcon, SunIcon } from "@heroicons/react/24/outline";

export default function ThemeSection() {
  const { theme, setTheme } = useThemeContext();

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : 'light');
  };

  return (
    <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
      <div>
        <h2 className="text-base/7 font-semibold text-base">Theme Settings</h2>
        <p className="mt-1 text-sm/6 text-base">Choose your preferred theme for the application.</p>
      </div>

      <div className="md:col-span-2">
        <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
          <div className="col-span-full">
            <label className="block text-sm/6 font-medium text-base">
              Theme
            </label>
            <div className="mt-2">
              <div className="flex items-center gap-x-4">
                <button
                  type="button"
                  onClick={toggleTheme}
                  className="flex items-center gap-x-3 rounded-md bg-white/5 px-4 py-3 text-sm font-medium text-base shadow-xs hover:bg-white/10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                >
                  {theme === 'light' ? (
                    <>
                      <SunIcon className="size-5" />
                      Light Mode
                    </>
                  ) : (
                    <>
                      <MoonIcon className="size-5" />
                      Dark Mode
                    </>
                  )}
                </button>
                <span className="text-sm text-gray-400">
                  Currently using {theme === 'light' ? 'light' : 'dark'} theme
                </span>
              </div>
            </div>
            <p className="mt-2 text-xs/5 text-gray-400">
              Switch between light and dark themes to match your preference.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
} 