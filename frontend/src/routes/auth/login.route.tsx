import {userManager} from "@/hooks/auth.oidc-config.ts";
import type {rootRoute} from "@/routes/root.route.tsx";
import {Button} from "@/components/ui/button.tsx";
import {createRoute} from "@tanstack/react-router";

const LoginPage = () => {
  const handleLogin = async () => {
    try {
      await userManager.signinRedirect();
    } catch (error) { console.error('login failed:', error)}
  }

  return (
    <div className={"flex min-h-screen items-center justify-center"}>
      <div className={"w-full max-w-md space-y-8 rounded-lg border p-8 shadow-lg"}>
        <div className={'text-center'}>
          <h1 className={"text-4xl font-semibold text-gray-900"}>Welcome to streamly</h1>
          <p className={"text-lg font-semibold text-gray-900"}>signin to continue</p>
        </div>
        <Button onClick={handleLogin} className={'w-full'} size="lg">
          Sign in
        </Button>
      </div>
    </div>
  )
}

export const loginRoute = (parent: typeof rootRoute) => createRoute({
  getParentRoute: () => parent,
  path: '/login',
  component: LoginPage
})
