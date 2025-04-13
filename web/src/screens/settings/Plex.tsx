import {
    Text,
    Button,
    Stack,
    Title,
    Center,
    Loader,
    Group,
    Paper,
    Table,
    Flex,
    Divider,
} from "@mantine/core";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {useEffect, useState} from "react";
import {useDisclosure} from "@mantine/hooks";

import {PlexConfig} from "@app/types/Plex";
import {APIClient} from "@api/APIClient";
import {PlexSettingsQueryOptions} from "@api/queries";
import {PlexSettingsKeys} from "@api/query_keys";
import {displayNotification} from "@components/notifications";
import {PlexSettings} from "@forms/settings/PlexSettings";
import {ConfirmDeleteButton} from "@components/alerts/ConfirmDeleteButton.tsx";

export const Plex = () => {
    const queryClient = useQueryClient();
    const [isReachable, setIsReachable] = useState<boolean | null>(null);
    const [opened, {open, close}] = useDisclosure(false);
    const {data: settings, isLoading} = useQuery(PlexSettingsQueryOptions());
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
            queryClient.resetQueries({queryKey: PlexSettingsKeys.config()});
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

    const updateMutation = useMutation({
        mutationFn: APIClient.plex.updateSettings,
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: PlexSettingsKeys.config()});
            close(); // Close modal after saving
            displayNotification({
                title: "Success",
                message: "Plex settings updated successfully",
                type: "success",
            });
        },
        onError: (error) => {
            displayNotification({
                title: "Update failed",
                message: error.message || "Could not update Plex settings",
                type: "error",
            });
        },
    });

    const handleFormSubmit = (values: PlexConfig) => {
        updateMutation.mutate(values);
    };

    const rows =
        settings && !isEmptySettings
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
        <Paper withBorder={true} p="md" shadow="xl">
            <Stack justify="center">
                <Title order={1} mt="md">
                    Plex Media Server
                </Title>
                <Text>
                    Manage the connection to your Plex Media Server here.
                </Text>
                <Divider/>
                {(isLoading || (!isEmptySettings && isReachable === null)) ? (
                    <>
                        <Center mt="md">
                            <Loader/>
                        </Center>
                    </>
                ) : (
                    <>
                        {isEmptySettings ? (
                            <>
                                <Group justify={"center"}>
                                    <Text size={"md"} fw={600}>Plex Setup Not Found</Text>
                                </Group>
                                <Group justify={"center"}>
                                    <Button onClick={open}>SETUP PLEX</Button>
                                </Group>
                            </>
                        ) : (
                            <>
                                <Flex align={"center"}>
                                    <Text size={"xl"} fw={600}>
                                        Connection Status:
                                    </Text>
                                    <Text c={isReachable ? "green" : "red"} size={"md"} fw={600} ml="xs" mr="xs"
                                          mt={3}>
                                        {isReachable ? "OK" : "Not Reachable"}
                                    </Text>
                                </Flex>
                                {settings && (
                                    <Table variant="vertical" verticalSpacing={"sm"} withTableBorder>
                                        <Table.Tbody>
                                            {rows.map((row, index) => (
                                                <Table.Tr key={index}>
                                                    <Table.Th w={180}>{row.label}</Table.Th>
                                                    <Table.Td>{row.value}</Table.Td>
                                                </Table.Tr>
                                            ))}
                                        </Table.Tbody>
                                    </Table>
                                )}
                                <Group mt="sm" justify={"flex-end"}>
                                    <ConfirmDeleteButton
                                        onConfirm={() => deleteMutation.mutate()}
                                        message={"All your Plex Media Server settings will be deleted."}
                                    />
                                    <Button variant={"outline"} onClick={open}>EDIT SETTINGS</Button>
                                </Group>
                            </>
                        )}
                    </>
                )}
            </Stack>
            <PlexSettings
                opened={opened}
                onClose={close}
                onSubmit={handleFormSubmit}
                defaultValues={settings}
            />
        </Paper>
    );
};