import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type SidebarState = {
  expanded: boolean
  toggle: () => void
}

export const useSidebarStore = create<SidebarState>()(
  persist(
    (set) => ({
      expanded: true,
      toggle: () => set((state) => ({ expanded: !state.expanded })),
    }),
    { name: 'streamly-sidebar' }
  )
)
