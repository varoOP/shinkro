interface NavItem {
  name: string;
  path: string;
  exact?: boolean;
}

export const NAV_ROUTES: Array<NavItem> = [
  { name: "Dashboard", path: "/", exact: true },
  { name: "Settings", path: "/settings" },
];
