import { createContext, useContext, useEffect, useState } from "react";

export const ThemeContext = createContext<ThemeContextType | null>(null);

type Theme = "light" | "dark";

type ThemeContextType = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

// Client-side function to get theme from localStorage
export function getClientTheme(): Theme {
  if (typeof window === 'undefined') return 'dark';

  try {
    const stored = localStorage.getItem('theme');
    return (stored as Theme) || 'dark';
  } catch {
    return 'dark';
  }
}

// Function to set theme in localStorage and cookies
export function setClientTheme(theme: Theme): void {
  if (typeof window === 'undefined') return;

  try {
    localStorage.setItem('theme', theme);
    // Also set a cookie for server-side access
    document.cookie = `theme=${theme}; path=/; max-age=31536000`; // 1 year
  } catch {
    // Ignore errors
  }
}

// Function to apply theme classes to DOM elements
export function applyThemeClasses(theme: Theme): void {
  if (typeof window === 'undefined') return;

  try {
    const html = document.documentElement;
    const body = document.body;

    // Remove existing theme classes
    html.classList.remove('light', 'dark');
    body.classList.remove('light', 'dark');

    // Add new theme class
    html.classList.add(theme);
    body.classList.add(theme);

    // Set data-theme attribute
    html.setAttribute('data-theme', theme);
    body.setAttribute('data-theme', theme);
  } catch {
    // Ignore errors
  }
}

export const ThemeContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [theme, setTheme] = useState<Theme>(() => {
    // Initialize with client-side theme on mount
    if (typeof window !== 'undefined') {
      return getClientTheme();
    }
    return 'dark';
  });

  // Update localStorage, cookie, and HTML classes when theme changes
  useEffect(() => {
    setClientTheme(theme);
    applyThemeClasses(theme);
  }, [theme]);

  // Initialize theme classes on mount
  useEffect(() => {
    if (typeof window !== 'undefined') {
      applyThemeClasses(theme);
    }
  }, []); // Only run on mount

  return <ThemeContext.Provider value={{ theme, setTheme }}>{children}</ThemeContext.Provider>;
};

export const useThemeContext = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useThemeContext must be used within a ThemeProvider");
  }
  return context;
};