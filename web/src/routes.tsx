import {
  createRootRouteWithContext,
  createRoute,
  createRouter,
  redirect,
  Outlet,
} from "@tanstack/react-router";
import { z } from "zod";
import { QueryClient } from "@tanstack/react-query";
import { APIClient } from "@api/APIClient";
import { Login, Onboarding } from "@screens/auth";
import { NotFound } from "@components/alerts/NotFound";
import { queryClient } from "@api/QueryClient";
import { SettingsContext } from "@utils/Context";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

// Onboarding route configuration
export const OnboardRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "onboard",
  beforeLoad: async () => {
    try {
      await APIClient.auth.canOnboard();
    } catch (e) {
      console.error("Onboarding not available, redirecting to login");
      throw redirect({ to: LoginRoute.to });
    }
  },
  component: Onboarding,
});

// Login route configuration
export const LoginRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "login",
  validateSearch: z.object({
    redirect: z.string().optional(),
  }),
  beforeLoad: ({ navigate }) => {
    APIClient.auth
      .canOnboard()
      .then(() => {
        console.info("Onboarding available, redirecting");
        navigate({ to: OnboardRoute.to });
      })
      .catch(() => {
        console.info("Onboarding not available, please login");
      });
  },
}).update({ component: Login });

// Root component that renders the outlet for child routes
export const RootComponent = () => {
  const settings = SettingsContext.useValue();
  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Outlet />
      {settings.debug ? (
        <>
          <TanStackRouterDevtools />
          <ReactQueryDevtools initialIsOpen={false} />
        </>
      ) : null}
    </div>
  );
};

// Root route configuration
export const RootRoute = createRootRouteWithContext<{
  queryClient: QueryClient;
}>()({
  component: RootComponent,
  notFoundComponent: NotFound,
});

// Create the route tree with only Login and Onboard routes
const routeTree = RootRoute.addChildren([LoginRoute, OnboardRoute]);

// Export the router with the simplified route tree
export const Router = createRouter({
  routeTree,
  context: {
    queryClient,
  },
});
