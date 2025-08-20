import {Tabs, Divider, Paper, Container} from "@mantine/core";
import {useNavigate, Outlet, useRouterState} from "@tanstack/react-router";
import {SiMyanimelist, SiPlex} from "react-icons/si";
import {FaCog, FaUserCog, FaKey, FaBell, FaMap} from "react-icons/fa";
import {BsStack} from "react-icons/bs";

const tabsList = [
    {value: "application", label: "Application", icon: <FaCog/>, path: "/settings"},
    {value: "user", label: "User", icon: <FaUserCog/>, path: "/settings/user"},
    {value: "api", label: "API Keys", icon: <FaKey/>, path: "/settings/api"},
    {value: "mapping", label: "Mapping", icon: <FaMap/>, path: "/settings/mapping"},
    {value: "notifications", label: "Notifications", icon: <FaBell/>, path: "/settings/notifications"},
    {value: "logs", label: "Logs", icon: <BsStack/>, path: "/settings/logs"},
    {value: "plex", label: "Plex", icon: <SiPlex/>, path: "/settings/plex"},
    {value: "mal", label: "MyAnimeList", icon: <SiMyanimelist/>, path: "/settings/mal"},
];

export const Settings = () => {
    const navigate = useNavigate();
    const pathname = useRouterState().location.pathname;
    const activeTab =
        tabsList.find((tab) => pathname.endsWith(tab.path))?.value || "application";

    return (
        <Container size={1200} px="md" component="main">
            <Paper mt="md" withBorder p={"md"} h={"100%"} mih={"80vh"}>
                <Tabs
                    value={activeTab}
                    onChange={(value) => {
                        const tab = tabsList.find((t) => t.value === value);
                        if (tab) void navigate({to: tab.path});
                    }}
                    variant="pills"
                    radius="sm"
                >
                    <Tabs.List justify="space-between" grow>
                        {tabsList.map((tab) => (
                            <Tabs.Tab key={tab.value} value={tab.value} leftSection={tab.icon}>
                                {tab.label}
                            </Tabs.Tab>
                        ))}
                    </Tabs.List>
                    <Divider mt="xs"/>
                    <Outlet/>
                </Tabs>
            </Paper>
        </Container>
    );
};
