import {Tabs, Divider, Paper} from "@mantine/core";
import {useParams, useNavigate} from "@tanstack/react-router";
import {Application} from "@screens/settings/Application";
import {User} from "@screens/settings/User";
import {Plex} from "@screens/settings/Plex";
import {Mal} from "@screens/settings/Mal";
import {Logs} from "@screens/settings/Logs";
import {Api} from "@screens/settings/Api";
import {Notifications} from "@screens/settings/Notifications";
import {SiMyanimelist, SiPlex} from "react-icons/si";
import {FaCog, FaUserCog, FaKey, FaBell} from "react-icons/fa";
import {BsStack} from "react-icons/bs";

const tabsList = [
    {value: "application", label: "Application", component: Application, icon: <FaCog/>},
    {value: "user", label: "User", component: User, icon: <FaUserCog/>},
    {value: "api", label: "API Keys", component: Api, icon: <FaKey/>},
    {value: "notifications", label: "Notifications", component: Notifications, icon: <FaBell/>},
    {value: "logs", label: "Logs", component: Logs, icon: <BsStack/>},
    {value: "plex", label: "Plex Media Server", component: Plex, icon: <SiPlex size={25}/>},
    {value: "mal", label: "MyAnimeList", component: Mal, icon: <SiMyanimelist size={25}/>},
];

export const Settings = () => {
    const params = useParams({strict: false});
    const navigate = useNavigate();
    const activeTab = params.activeTab ?? "application";
    const isValidTab = tabsList.some(tab => tab.value === activeTab);
    const currentTab = isValidTab ? activeTab : "application";

    return (
        <div>
            <Paper mt="md" withBorder p={"md"} h={"100%"} mih={"80vh"}>
                <Tabs
                    value={currentTab}
                    onChange={(value) => {
                        if (value === "application" || !value) {
                            navigate({to: "/settings", replace: true});
                        } else {
                            navigate({to: "/settings/$activeTab", params: {activeTab: value}});
                        }
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

                    <Divider size="md" mt="xs"/>
                    {tabsList.map((tab) => (
                        <Tabs.Panel key={tab.value} value={tab.value} mt="xs">
                            <tab.component/>
                        </Tabs.Panel>
                    ))}
                </Tabs>
            </Paper>
        </div>
    );
};
