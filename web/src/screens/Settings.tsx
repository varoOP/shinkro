import {Paper, Tabs, rem, Divider} from "@mantine/core";
import {SiPlex, SiMyanimelist} from "react-icons/si";
import {Plex} from "@screens/settings/Plex";


export const Settings = () => {
    const iconStyle = {width: rem(30), height: rem(30)};

    return (
        <main style={{height: '100%', width: '100%'}}>
            <Paper style={{
                width: '100%',
                flexGrow: 1,
                padding: '8px'
            }}>
                <Tabs defaultValue="plex" variant="pills" radius="sm">
                    <Tabs.List grow justify="center">
                        <Tabs.Tab value="application">
                            Application
                        </Tabs.Tab>
                        <Tabs.Tab value="user">
                            User
                        </Tabs.Tab>
                        <Tabs.Tab value="plex" leftSection={<SiPlex style={iconStyle}/>}>
                        </Tabs.Tab>
                        <Tabs.Tab value="mal" leftSection={<SiMyanimelist style={iconStyle}/>}>
                        </Tabs.Tab>
                        <Tabs.Tab value="logs">
                            Logs
                        </Tabs.Tab>
                    </Tabs.List>
                    <Divider size="md" mt="xs"/>
                    <Tabs.Panel mt="xs" value="plex">
                        <Plex/>
                    </Tabs.Panel>

                    <Tabs.Panel value="mal">
                        Load Mal Auth Component here
                    </Tabs.Panel>
                </Tabs>

            </Paper>
        </main>
    );
};
