import {homeRoute} from "@/routes/home/home.route.tsx";
import {rootRoute} from "@/routes/root.route.tsx";
import eventRoute from "@/routes/events/event.route.tsx";

const HomeRoute = homeRoute(rootRoute);
const EventRoute = eventRoute(rootRoute);

export const routes = [
  HomeRoute,
  EventRoute,
]
