import {
  createRootRouteWithContext,
  createRoute,
  createRouter,
  Navigate,
  Outlet,
  redirect,
} from "@tanstack/react-router";
import { z } from "zod";
import { QueryClient } from "@tanstack/react-query";
import { APIClient } from "@api/APIClient";
import { Login, Onboarding } from "@screens/auth";
import { Layout } from "@components/layout";
import { NotFound } from "@components/alerts/NotFound";
import { Dashboard } from "@screens/Dashboard";
import { Settings } from "@screens/Settings";
import { AuthContext, SettingsContext } from "@utils/Context";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { queryClient } from "@api/QueryClient";

const DashboardRoute = createRoute({
  getParentRoute: () => AuthIndexRoute,
  path: "/",
  loader: () => {
    // https://tanstack.com/router/v1/docs/guide/deferred-data-loading#deferred-data-loading-with-defer-and-await
    // TODO load stats

    // TODO load recent releases

    return {};
  },
  component: Dashboard,
});

export const OnboardRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "onboard",
  beforeLoad: async () => {
    // Check if onboarding is available for this instance
    // and redirect if needed
    try {
      await APIClient.auth.canOnboard();
    } catch {
      console.error("onboarding not available, redirect to login");

      throw redirect({
        to: LoginRoute.to,
      });
    }
  },
  component: Onboarding,
});

export const LoginRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "login",
  validateSearch: z.object({
    redirect: z.string().optional(),
  }),
  beforeLoad: ({ navigate }) => {
    // handle canOnboard
    APIClient.auth
      .canOnboard()
      .then(() => {
        console.info("onboarding available, redirecting");

        navigate({ to: OnboardRoute.to });
      })
      .catch(() => {
        console.info("onboarding not available, please login");
      });
  },
}).update({ component: Login });

export const AuthRoute = createRoute({
  getParentRoute: () => RootRoute,
  id: "auth",
  // Before loading, authenticate the user via our auth context
  // This will also happen during prefetching (e.g. hovering over links, etc.)
  beforeLoad: ({ context, location }) => {
    // If the user is not logged in, check for item in localStorage
    if (!AuthContext.get().isLoggedIn) {
      throw redirect({
        to: LoginRoute.to,
        search: {
          // Use the current location to power a redirect after login
          // (Do not use `router.state.resolvedLocation` as it can
          // potentially lag behind the actual current location)
          redirect: location.href,
        },
      });
    }

    // Otherwise, return the user in context
    return context;
  },
});

function AuthenticatedLayout() {
  const isLoggedIn = AuthContext.useSelector((s) => s.isLoggedIn);
  if (!isLoggedIn) {
    const redirect =
      location.pathname.length > 1
        ? { redirect: location.pathname }
        : undefined;
    return <Navigate to="/login" search={redirect} />;
  }

  return (
    <div className="full-height-center">
      <Layout />
    </div>
  );
}

export const SettingsRoute = createRoute({
  getParentRoute: () => AuthIndexRoute,
  path: "settings",
  pendingMs: 3000,
  component: Settings,
});

export const AuthIndexRoute = createRoute({
  getParentRoute: () => AuthRoute,
  component: AuthenticatedLayout,
  id: "authenticated-routes",
});

// Root component that renders the outlet for child routes
export const RootComponent = () => {
  const settings = SettingsContext.useValue();
  return (
    <div className="full-height-center">
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

export const RootRoute = createRootRouteWithContext<{
  queryClient: QueryClient;
}>()({
  component: RootComponent,
  notFoundComponent: NotFound,
});

const authenticatedTree = AuthRoute.addChildren([
  AuthIndexRoute.addChildren([DashboardRoute, SettingsRoute]),
]);
const routeTree = RootRoute.addChildren([
  LoginRoute,
  OnboardRoute,
  authenticatedTree,
]);

export const Router = createRouter({
  routeTree,
  context: {
    queryClient,
  },
});
