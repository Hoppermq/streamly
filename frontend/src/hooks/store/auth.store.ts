import type {User} from "oidc-client-ts";
import {create} from "zustand";
import {persist} from "zustand/middleware";

type AuthState = {
  user: User | null,
  isAuthenticated: boolean,
  isLoading: boolean,
  setUser: (user: User | null) => void,
  setLoading: (loading: boolean) => void,
  logout: () => void,
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      setUser: (user: User | null) => set({ user, isAuthenticated: !!user, isLoading: false}),
      setLoading: (loading: boolean) => set({ isLoading: loading }),
      logout: () => set({ user: null, isAuthenticated: false }),
    }),
    { name: 'streamly-auth' }
  )
)
