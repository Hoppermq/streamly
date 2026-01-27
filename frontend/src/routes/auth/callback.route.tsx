import {createRoute, useNavigate} from "@tanstack/react-router";
import {rootRoute} from "@/routes/root.route.tsx";
import {useAuthStore} from "@/hooks/store/auth.store.ts";
import {useEffect} from "react";
import {userManager} from "@/hooks/auth.oidc-config.ts";

export const callbackRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/callback',
  component: () => {
    const navigate = useNavigate()
    const setUser = useAuthStore((state) => state.setUser)

    useEffect(() => {
      userManager.signinRedirectCallback().then((user) => {
        setUser(user)
        navigate({ to: '/' })
      }).catch((err) => {
        console.error('login error:', err)
        navigate({href: '/login'})
      })
    })

    return <div>Processing login</div>
  }
})
