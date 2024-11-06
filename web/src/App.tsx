import { useEffect } from "react";
import { RouterProvider } from "@tanstack/react-router";
import { QueryClientProvider } from "@tanstack/react-query";
import { Notifications } from "@mantine/notifications";
import { Router } from "@app/routes";
import { routerBasePath } from "@utils";
import { queryClient } from "@api/QueryClient";
import { SettingsContext } from "@utils/Context";

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof Router;
  }
}

export function App() {
  const [, setSettings] = SettingsContext.use();

  useEffect(() => {
    const themeMediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const handleThemeChange = (e: MediaQueryListEvent) => {
      setSettings((prevState) => ({ ...prevState, darkTheme: e.matches }));
    };

    themeMediaQuery.addEventListener("change", handleThemeChange);
    return () =>
      themeMediaQuery.removeEventListener("change", handleThemeChange);
  }, [setSettings]);

  return (
    <QueryClientProvider client={queryClient}>
      <Notifications position="top-right" />
      <RouterProvider basepath={routerBasePath()} router={Router} />
    </QueryClientProvider>
  );
}
