import { MdSpaceDashboard, MdSettings, MdHistory } from "react-icons/md";
import { BsStack } from "react-icons/bs";
import { IconType } from "react-icons";

interface NavItem {
  name: string;
  path: string;
  exact?: boolean;
  icon: IconType;
}

export const NAV_ROUTES: Array<NavItem> = [
  { name: "Dashboard", path: "/", exact: true, icon: MdSpaceDashboard },
  { name: "Plex Payloads", path: "/plex-payloads", icon: MdHistory },
  { name: "Logs", path: "/logs", icon: BsStack },
  { name: "Settings", path: "/settings", icon: MdSettings },
];
