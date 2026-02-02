import {homeRoute} from "@/routes/home/home.route.tsx";
import {rootRoute} from "@/routes/root.route.tsx";
import eventRoute from "@/routes/events/event.route.tsx";
import {loginRoute} from "@/routes/auth/login.route.tsx";
import {callbackRoute} from "@/routes/auth/callback.route.tsx";
import {silentCallbackRoute} from "@/routes/auth/silent-callback.route.tsx";

const HomeRoute = homeRoute(rootRoute);
const EventRoute = eventRoute(rootRoute);
const LoginRoute = loginRoute(rootRoute);
const CallBackRoute = callbackRoute();
const SilentCallbackRoute = silentCallbackRoute;

export const routes = [
  HomeRoute,
  EventRoute,
  LoginRoute,
  CallBackRoute,
  SilentCallbackRoute,
]
