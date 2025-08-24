interface NavItem {
  name: string;
  path: string;
  exact?: boolean;
}

export const NAV_ROUTES: Array<NavItem> = [
  { name: "Dashboard", path: "/", exact: true },
  { name: "Logs", path: "/logs" },
  { name: "Settings", path: "/settings" },
];
