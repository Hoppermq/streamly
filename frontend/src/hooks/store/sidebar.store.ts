import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type SidebarState = {
  expanded: boolean
  toggle: () => void
  selectedPath: string
  setSelectedPath: (path: string) => void
}

export const useSidebarStore = create<SidebarState>()(
  persist(
    (set) => ({
      expanded: true,
      toggle: () => set((state) => ({ expanded: !state.expanded })),
      selectedPath: '/',
      setSelectedPath: (path) => set({ selectedPath: path }),
    }),
    { name: 'streamly-sidebar' }
  )
)
