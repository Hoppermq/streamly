import HomePage from "@/pages/home";
import {createRoute} from "@tanstack/react-router";
import {ROUTES} from "@/lib/constants/routes.ts";
import type {rootRoute} from "@/routes/root.route.tsx";

export const homeRoute =  (parentRoute: typeof rootRoute) =>
  createRoute({
    path: ROUTES.HOME,
    component: HomePage,
    getParentRoute: () => parentRoute
  })
