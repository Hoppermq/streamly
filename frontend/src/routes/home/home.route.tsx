import HomePage from "@/pages/home";
import {createRoute, type RootRoute} from "@tanstack/react-router";
import {ROUTES} from "@/lib/constants/routes.ts";

export const homeRoute =  (parentRoute: RootRoute) =>
  createRoute({
    path: ROUTES.HOME,
    component: HomePage,
    getParentRoute: () => parentRoute
  })
