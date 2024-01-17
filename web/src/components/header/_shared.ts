interface NavItem {
    name: string;
    path: string;
  }
  
  export interface RightNavProps {
    logoutMutation: () => void;
  }
  
  export const NAV_ROUTES: Array<NavItem> = [
    { name: "Anime Updates", path: "/" },
    { name: "Plex Payloads", path: "/plexPayloads" },
    { name: "Settings", path: "/settings" },
    { name: "Logs", path: "/logs" }
  ];