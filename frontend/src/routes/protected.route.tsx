import {createRoute, redirect} from "@tanstack/react-router";
import {rootRoute} from "@/routes/root.route.tsx";
import HomePage from "@/pages/home";
import {useAuthStore} from "@/hooks/store/auth.store.ts";

export const protectedRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/dashboard", //should be org here
  beforeLoad: async () => {
    const { isAuthenticated } = useAuthStore.getState();
    if (!isAuthenticated) { throw redirect({href: '/login'})}
  },
  component: HomePage
})
