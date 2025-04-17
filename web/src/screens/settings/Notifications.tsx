import {Stack} from "@mantine/core";
import {SettingsSectionHeader} from "@screens/settings/components.tsx";

export const Notifications = () => {
    return (
        <Stack>
            <SettingsSectionHeader title={"Notifications"} description={"Manage your notification settings here."}/>
        </Stack>
    );
}