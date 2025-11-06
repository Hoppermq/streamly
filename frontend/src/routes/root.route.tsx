import {createRootRoute, Outlet} from "@tanstack/react-router";
import RootContainer from "@/containers/root.container.tsx";
import {TanStackRouterDevtools} from "@tanstack/react-router-devtools";

export const rootRoute = createRootRoute({
  component: () => (
    <RootContainer>
      <>
        <Outlet />
        <TanStackRouterDevtools />
      </>
    </RootContainer>
  ),
})
