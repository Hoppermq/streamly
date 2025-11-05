import {createRoute, type RootRoute} from "@tanstack/react-router";
import EventPage from "@/pages/events";
import {ROUTES} from "@/lib/constants/routes.ts";

export default (parentRoute: RootRoute) =>
  createRoute({
    path: ROUTES.EVENTS,
    component: EventPage,
    getParentRoute: () => parentRoute
  })
