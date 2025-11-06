import type { Theme } from "@/lib/types"
import { create } from 'zustand'

type ThemeState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  initTheme: () => void
}

export const useThemeStore = create<ThemeState>((set) => ({
  theme: 'system',
  setTheme: (theme: Theme) => {
    localStorage.setItem('streamly-theme', theme)
    set({ theme })
    applyTheme(theme)
  },
  initTheme: () => {
    const stored = localStorage.getItem('streamly-theme') as Theme
    const initialTheme = stored || 'system'
    set({ theme: initialTheme })
    applyTheme(initialTheme)
  }
}))

const applyTheme = (theme: Theme) => {
  const root = window.document
  .documentElement
  root.classList.remove('light', 'dark')

  if (theme === 'system') {
    const systemTheme: Theme = window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light'

    root.classList.add(systemTheme)
  } else {
    root.classList.add(theme)
  }
}
