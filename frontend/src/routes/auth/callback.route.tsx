import {createRoute, useNavigate} from "@tanstack/react-router";
import {rootRoute} from "@/routes/root.route.tsx";
import {useAuthStore} from "@/hooks/store/auth.store.ts";
import {useEffect} from "react";
import {userManager} from "@/hooks/auth.oidc-config.ts";

export const callbackRoute = (parent: typeof rootRoute) => createRoute({
  getParentRoute: () => parent,
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
        navigate({to: '/login'})
      })
    }, [navigate, setUser])

    return <div>Processing login</div>
  }
})
