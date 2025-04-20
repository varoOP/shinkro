import {
    Button,
    Stack,
    Group,
    Table,
} from "@mantine/core";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {useEffect, useState} from "react";
import {useDisclosure} from "@mantine/hooks";

import {APIClient} from "@api/APIClient";
import {PlexSettingsQueryOptions} from "@api/queries";
import {PlexSettingsKeys} from "@api/query_keys";
import {displayNotification} from "@components/notifications";
import {PlexSettings} from "@forms/settings/PlexSettings";
import {ConfirmDeleteButton} from "@components/alerts/ConfirmDeleteButton";
import {
    SettingsSectionHeader,
    StatusIndicator,
    CenteredEmptyState,
} from "./components";

export const Plex = () => {
    const queryClient = useQueryClient();
    const [isReachable, setIsReachable] = useState<boolean | null>(null);
    const [opened, {open, close}] = useDisclosure(false);
    const {data: settings} = useQuery(PlexSettingsQueryOptions());
    const isEmptySettings = !settings || Object.keys(settings).length === 0;

    useEffect(() => {
        if (!isEmptySettings) {
            APIClient.plex
                .test(settings)
                .then(() => setIsReachable(true))
                .catch(() => setIsReachable(false));
        } else {
            setIsReachable(null);
        }
    }, [settings, isEmptySettings]);

    const deleteMutation = useMutation({
        mutationFn: APIClient.plex.delete,
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: PlexSettingsKeys.config()});
            displayNotification({
                title: "Success",
                message: "Plex settings deleted successfully.",
                type: "success",
            });
        },
        onError: (error) => {
            displayNotification({
                title: "Delete failed",
                message: error.message || "Could not delete Plex settings",
                type: "error",
            });
        },
    });

    const rows = settings && !isEmptySettings
        ? [
            {label: 'Plex User', value: settings.plex_user},
            {label: 'Anime Libraries', value: settings.anime_libs.join(', ')},
            {label: 'Host', value: settings.host},
            {label: 'Port', value: settings.port.toString()},
            {label: 'TLS', value: settings.tls ? 'Enabled' : 'Disabled'},
            {label: 'TLS Skip Verification', value: settings.tls_skip ? 'Enabled' : 'Disabled'},
        ]
        : [];

    return (
        <>
            <Stack>
                <SettingsSectionHeader
                    title="Plex Media Server"
                    description="Manage the connection to your Plex Media Server here."
                />
                {isEmptySettings ? (
                    <CenteredEmptyState
                        message="Plex Setup Not Found"
                        button={<Button onClick={open}>SETUP PLEX</Button>}
                    />
                ) : (
                    <>
                        <StatusIndicator
                            label="Connection Status:"
                            status={isReachable}
                            loadStatus={isReachable === null}
                        />
                        <Table variant="vertical" verticalSpacing="sm" withTableBorder>
                            <Table.Tbody>
                                {rows.map((row, index) => (
                                    <Table.Tr key={index}>
                                        <Table.Th w={180}>{row.label}</Table.Th>
                                        <Table.Td>{row.value}</Table.Td>
                                    </Table.Tr>
                                ))}
                            </Table.Tbody>
                        </Table>
                        <Group justify="flex-end">
                            <ConfirmDeleteButton
                                onConfirm={() => deleteMutation.mutate()}
                                message="All your Plex Media Server settings will be deleted."
                            />
                            <Button onClick={open}>EDIT SETTINGS</Button>
                        </Group>
                    </>
                )}
            </Stack>
            <PlexSettings
                opened={opened}
                onClose={close}
                defaultValues={settings}
            />
        </>
    );
};