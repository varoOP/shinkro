import {Button, Stack, Group, Text, PasswordInput, ActionIcon, CopyButton} from "@mantine/core";
import {CenteredEmptyState, SettingsSectionHeader} from "@screens/settings/components.tsx";
import {useSuspenseQuery, useQueryClient, useMutation} from "@tanstack/react-query";
import {ApikeysQueryOptions} from "@api/queries.ts";
import {displayNotification} from "@components/notifications";
import {useDisclosure} from "@mantine/hooks";
import {APIClient} from "@api/APIClient.ts";
import {ApiAddKey} from "@forms/settings/ApiAddKey.tsx";
import {FaTrash, FaCopy} from "react-icons/fa";
import {useMantineColorScheme} from '@mantine/core';


export const Api = () => {
    const [opened, {open, close}] = useDisclosure(false);
    const {data: keys} = useSuspenseQuery(ApikeysQueryOptions());
    const queryClient = useQueryClient();

    const {colorScheme} = useMantineColorScheme();
    const isDark = colorScheme === 'dark';

    const mutation = useMutation({
        mutationFn: (key: string) => APIClient.apikeys.delete(key),
        onSuccess: () => {
            queryClient.invalidateQueries(ApikeysQueryOptions());
            displayNotification({
                title: "API Key Deleted",
                type: "success",
            });
        },
    });

    return (
        <>
            <Stack>
                <SettingsSectionHeader
                    title={"API Keys"}
                    description={"Manage your shinkro API keys here."}
                    right={
                        <Button onClick={open} size={"compact-xs"}>
                            ADD NEW
                        </Button>
                    }
                />
                {keys && keys.length > 0 ? (
                    <div>
                        <Group justify="flex-start" grow gap={"xl"}>
                            <div style={{minWidth: "180px", maxWidth: "282px"}}>
                                <Text fw={700}>
                                    Name
                                </Text>
                            </div>
                            <Text fw={700}>
                                API Key
                            </Text>

                        </Group>
                        {keys.map((key) => (
                            <Group key={key.key} gap={"xl"} justify="flex-start" grow mt={"md"}>
                                <Stack justify="center" align={"stretch"} miw={"150px"} maw={"280px"}>
                                    <Text
                                        fw={900}
                                        c={isDark ? "plex" : "mal"}
                                    >
                                        {key.name}
                                    </Text>
                                </Stack>
                                <Stack justify="flex-start">
                                    <Group grow>
                                        <PasswordInput
                                            value={key.key}
                                            readOnly
                                            variant={"filled"}
                                            leftSection={
                                                <div style={{pointerEvents: "all"}}>
                                                    <CopyButton value={key.key}>
                                                        {({copied, copy}) => (
                                                            <ActionIcon color={copied ? 'teal' : 'plex'} onClick={copy}>
                                                                <FaCopy/>
                                                            </ActionIcon>
                                                        )}
                                                    </CopyButton>
                                                </div>
                                            }
                                        />
                                        <Group>
                                            <ActionIcon
                                                onClick={() => mutation.mutate(key.key)} color={"red"}
                                                variant="outline"
                                            >
                                                <FaTrash/>
                                            </ActionIcon>
                                        </Group>
                                    </Group>
                                </Stack>
                            </Group>
                        ))}
                    </div>
                ) : (
                    <CenteredEmptyState
                        message={"No API Keys Found"}
                    />
                )}
            </Stack>
            <ApiAddKey opened={opened} onClose={close}/>
        </>
    );
}