import {createRoute} from "@tanstack/react-router";
import EventPage from "@/pages/events";
import {ROUTES} from "@/lib/constants/routes.ts";
import type {rootRoute} from "@/routes/root.route.tsx";

export default (parentRoute: typeof rootRoute) =>
  createRoute({
    path: ROUTES.EVENTS,
    component: EventPage,
    getParentRoute: () => parentRoute
  })
