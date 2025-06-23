import {
    Button, Stack, Group, Text, PasswordInput,
    ActionIcon, CopyButton, useMantineColorScheme, Divider
} from "@mantine/core";
import {CenteredEmptyState, SettingsSectionHeader} from "@screens/settings/components.tsx";
import {useQuery, useQueryClient, useMutation} from "@tanstack/react-query";
import {ApikeysQueryOptions} from "@api/queries.ts";
import {displayNotification} from "@components/notifications";
import {useDisclosure} from "@mantine/hooks";
import {APIClient} from "@api/APIClient.ts";
import {ApiAddKey} from "@forms/settings/ApiAddKey.tsx";
import {FaCopy} from "react-icons/fa";
import {ConfirmDeleteIcon} from "@components/alerts/ConfirmDeleteIcon";

export const Api = () => {
    const [opened, {open, close}] = useDisclosure(false);
    const {data: keys} = useQuery(ApikeysQueryOptions());
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
        <main>
            <SettingsSectionHeader
                title="API Keys"
                description="Manage your shinkro API keys here."
            />
            {keys?.length ? (
                <Stack mt={"md"} mb={"md"}>
                    <KeyList keys={keys} onDelete={mutation.mutate} isDark={isDark}/>
                </Stack>
            ) : (
                <CenteredEmptyState message="No API Keys Found"/>
            )}
            <Group justify={"center"}>
                <Button onClick={open} size="compact-xs">
                    ADD NEW
                </Button>
            </Group>
            <ApiAddKey opened={opened} onClose={close}/>
        </main>
    );
};

interface KeyListProps {
    keys: { name: string; key: string }[];
    onDelete: (key: string) => void;
    isDark: boolean;
}

const KeyList = ({keys, onDelete, isDark}: KeyListProps) => (
    <div>
        <Group justify="flex-start" grow gap="xl">
            <div style={{minWidth: "180px", maxWidth: "640px"}}>
                <Text fw={700}>Name</Text>
            </div>
            <Text fw={700}>API Key</Text>
        </Group>
        <Divider mt={"xs"}/>

        {keys.map((key) => (
            <KeyRow key={key.key} name={key.name} value={key.key} onDelete={onDelete} isDark={isDark}/>
        ))}
    </div>
);

interface KeyRowProps {
    name: string;
    value: string;
    onDelete: (key: string) => void;
    isDark: boolean;
}

const KeyRow = ({name, value, onDelete, isDark}: KeyRowProps) => (
    <Stack>
        <Group gap="xl" justify="flex-start" grow mt="md">
            <Stack justify="center" align="stretch" miw="150px" maw="500px">
                <Text fw={900} truncate c={isDark ? "plex" : "mal"}>
                    {name}
                </Text>
            </Stack>

            <Stack justify="flex-start">
                <Group grow>
                    <PasswordInput
                        value={value}
                        readOnly
                        variant="filled"
                        leftSection={
                            <div style={{pointerEvents: "all"}}>
                                <CopyButton value={value}>
                                    {({copied, copy}) => (
                                        <ActionIcon color={copied ? 'teal' : ''} onClick={copy}>
                                            <FaCopy/>
                                        </ActionIcon>
                                    )}
                                </CopyButton>
                            </div>
                        }
                    />
                    <ConfirmDeleteIcon
                        onConfirm={() => onDelete(value)}
                        title="Delete API Key"
                        message={`Are you sure you want to delete the API key "${name}"?`}
                        confirmText="DELETE"
                        variant="outline"
                    />
                </Group>
            </Stack>
        </Group>
        <Divider/>
    </Stack>
);