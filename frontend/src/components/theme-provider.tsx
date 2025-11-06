import { useThemeStore } from "@/hooks/store/theme.store"
import type { Theme } from "@/lib/types"
import { createContext, useContext, useEffect } from "react"
import type { FC } from "react"

type ThemeProviderProps = {
  children: React.ReactNode
  defaultTheme?: Theme
  storageKey?: string
}

const ThemeCtx = createContext<{
  theme: Theme
  setTheme: (theme: Theme) => void
} | null>(null)


export const ThemeProvider : FC<ThemeProviderProps> = ({
  children,
  defaultTheme = 'system',
  storageKey = 'streamly-theme',
}) => {
  const { theme, setTheme, initTheme } = useThemeStore()

  useEffect(() => {
    initTheme()
  }, [initTheme])

  useEffect(() => {
    if (theme !== 'system') return

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = () => {
      const root = document.documentElement
      root.classList.remove('light', 'dark')
      root.classList.add(mediaQuery.matches ? 'dark' : 'light')
    }

    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [theme])

  return (
    <ThemeCtx.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeCtx.Provider>
  )
}

export const useTheme = () => {
  const ctx = useContext(ThemeCtx)
  if (!ctx) {
    throw new Error('useTheme must be used within a theme provider')
  }

  return ctx
}
