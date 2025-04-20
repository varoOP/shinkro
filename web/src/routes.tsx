import {
    createRootRouteWithContext,
    createRoute,
    createRouter,
    Navigate,
    Outlet,
    redirect,
} from "@tanstack/react-router";
import {z} from "zod";
import {QueryClient} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient";
import {Login, Onboarding} from "@screens/auth";
import {Layout} from "@components/layout";
import {NotFound} from "@components/alerts/NotFound";
import {Dashboard} from "@screens/Dashboard";
import {Settings} from "@screens/Settings";
import {Application} from "@screens/settings/Application";
import {User} from "@screens/settings/User";
import {Api} from "@screens/settings/Api";
import {Notifications} from "@screens/settings/Notifications";
import {Logs} from "@screens/settings/Logs";
import {Plex} from "@screens/settings/Plex";
import {Mal} from "@screens/settings/Mal";
import {AuthContext, SettingsContext} from "@utils/Context";
import {TanStackRouterDevtools} from "@tanstack/router-devtools";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";
import {queryClient} from "@api/QueryClient";
import {MalAuthCallback} from "@screens/MalAuthCallback.tsx";

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
    beforeLoad: ({navigate}) => {
        // handle canOnboard
        APIClient.auth
            .canOnboard()
            .then(() => {
                console.info("onboarding available, redirecting");

                navigate({to: OnboardRoute.to});
            })
            .catch(() => {
                console.info("onboarding not available, please login");
            });
    },
}).update({component: Login});

export const AuthRoute = createRoute({
    getParentRoute: () => RootRoute,
    id: "auth",
    // Before loading, authenticate the user via our auth context
    // This will also happen during prefetching (e.g. hovering over links, etc.)
    beforeLoad: ({context, location}) => {
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
                ? {redirect: location.pathname}
                : undefined;
        return <Navigate to="/login" search={redirect}/>;
    }

    return (
        <div>
            <Layout/>
        </div>
    );
}

export const MalAuthRoute = createRoute({
    getParentRoute: () => AuthRoute,
    component: MalAuthCallback,
    path: "malauth/callback",
    pendingMs: 3000,
});

export const SettingsRoute = createRoute({
    getParentRoute: () => AuthIndexRoute,
    path: "settings",
    pendingMs: 3000,
    component: Settings,
});

export const SettingsApplicationRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/",
    pendingMs: 3000,
    component: Application,
});

export const SettingsUserRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/user",
    pendingMs: 3000,
    component: User,
});

export const SettingsApiRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/api",
    pendingMs: 3000,
    component: Api,
});

export const SettingsNotificationsRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/notifications",
    pendingMs: 3000,
    component: Notifications,
});

export const SettingsLogsRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/logs",
    pendingMs: 3000,
    component: Logs,
});

export const SettingsPlexRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/plex",
    pendingMs: 3000,
    component: Plex,
});

export const SettingsMalRoute = createRoute({
    getParentRoute: () => SettingsRoute,
    path: "/mal",
    pendingMs: 3000,
    component: Mal,
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
        <div className="full-height">
            <Outlet/>
            {settings.debug ? (
                <>
                    <TanStackRouterDevtools/>
                    <ReactQueryDevtools initialIsOpen={false}/>
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

const settingsRouteTree = SettingsRoute.addChildren([
    SettingsApplicationRoute,
    SettingsUserRoute,
    SettingsApiRoute,
    SettingsNotificationsRoute,
    SettingsLogsRoute,
    SettingsPlexRoute,
    SettingsMalRoute,
]);
const authenticatedTree = AuthRoute.addChildren([
    AuthIndexRoute.addChildren([DashboardRoute, settingsRouteTree]),
]);
const routeTree = RootRoute.addChildren([
    LoginRoute,
    OnboardRoute,
    MalAuthRoute,
    authenticatedTree,
]);

export const Router = createRouter({
    routeTree,
    context: {
        queryClient,
    },
});
