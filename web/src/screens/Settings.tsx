import { Container, Paper, Tabs } from "@mantine/core";
import { Outlet, useRouterState, useNavigate } from "@tanstack/react-router";
import { FaCog, FaUserCog, FaKey, FaBell, FaMap } from "react-icons/fa";
import { BsStack } from "react-icons/bs";

export const Settings = () => {
    const pathname = useRouterState().location.pathname;
    const navigate = useNavigate();
    
    // Determine active tab based on current path
    const getActiveTab = () => {
        if (pathname === '/settings') return 'application';
        if (pathname === '/settings/user') return 'user';
        if (pathname === '/settings/api') return 'api';
        if (pathname === '/settings/mapping') return 'mapping';
        if (pathname === '/settings/notifications') return 'notifications';
        if (pathname === '/settings/logs') return 'logs';
        return 'application';
    };

    const handleTabChange = (value: string | null) => {
        if (!value) return;
        
        const routeMap: Record<string, string> = {
            'application': '/settings',
            'user': '/settings/user',
            'api': '/settings/api',
            'mapping': '/settings/mapping',
            'notifications': '/settings/notifications',
            'logs': '/settings/logs',
        };
        
        const route = routeMap[value];
        if (route) {
            navigate({ to: route });
        }
    };

    return (
        <Container size={1200} px="md" component="main">
            <Paper mt="md" withBorder p={"md"} h={"100%"} mih={"80vh"}>
                {pathname === '/settings' || (pathname.startsWith('/settings/') && !pathname.includes('/plex') && !pathname.includes('/mal')) ? (
                    <Tabs value={getActiveTab()} onChange={handleTabChange} variant="pills">
                        <Tabs.List grow>
                            <Tabs.Tab value="application" leftSection={<FaCog size={14} />}>
                                Application
                            </Tabs.Tab>
                            <Tabs.Tab value="user" leftSection={<FaUserCog size={14} />}>
                                User
                            </Tabs.Tab>
                            <Tabs.Tab value="api" leftSection={<FaKey size={14} />}>
                                API Keys
                            </Tabs.Tab>
                            <Tabs.Tab value="mapping" leftSection={<FaMap size={14} />}>
                                Mapping
                            </Tabs.Tab>
                            <Tabs.Tab value="notifications" leftSection={<FaBell size={14} />}>
                                Notifications
                            </Tabs.Tab>
                            <Tabs.Tab value="logs" leftSection={<BsStack size={14} />}>
                                Logs
                            </Tabs.Tab>
                        </Tabs.List>
                    </Tabs>
                ) : null}
                <Outlet />
            </Paper>
        </Container>
    );
};
