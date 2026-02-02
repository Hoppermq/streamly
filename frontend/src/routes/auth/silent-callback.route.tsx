import { createRoute } from "@tanstack/react-router";
import { rootRoute } from "@/routes/root.route.tsx";
import { useEffect } from "react";
import { userManager } from "@/hooks/auth.oidc-config.ts";

const SilentCallbackPage = () => {
  useEffect(() => {
    userManager.signinSilentCallback().catch((error) => {
      console.error("Silent callback error:", error);
    });
  }, []);

  return null;
};

export const silentCallbackRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/auth/silent-callback",
  component: SilentCallbackPage,
});
