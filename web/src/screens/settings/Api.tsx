import {
    Button, Stack, Group, Text, PasswordInput,
    ActionIcon, CopyButton, useMantineColorScheme, Divider
} from "@mantine/core";
import {CenteredEmptyState, SettingsSectionHeader} from "@screens/settings/components.tsx";
import {useSuspenseQuery, useQueryClient, useMutation} from "@tanstack/react-query";
import {ApikeysQueryOptions} from "@api/queries.ts";
import {displayNotification} from "@components/notifications";
import {useDisclosure} from "@mantine/hooks";
import {APIClient} from "@api/APIClient.ts";
import {ApiAddKey} from "@forms/settings/ApiAddKey.tsx";
import {FaTrash, FaCopy} from "react-icons/fa";

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
                    title="API Keys"
                    description="Manage your shinkro API keys here."
                />
                {keys?.length ? (
                    <KeyList keys={keys} onDelete={mutation.mutate} isDark={isDark}/>
                ) : (
                    <CenteredEmptyState message="No API Keys Found"/>
                )}
                <Group justify={"center"}>
                    <Button onClick={open} size="compact-xs">
                        ADD NEW
                    </Button>
                </Group>
            </Stack>
            <ApiAddKey opened={opened} onClose={close}/>
        </>
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
                <Text fw={900} c={isDark ? "plex" : "mal"}>
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
                    <ActionIcon
                        maw={"20px"}
                        onClick={() => onDelete(value)}
                        color="red"
                        variant="outline"
                    >
                        <FaTrash/>
                    </ActionIcon>
                </Group>
            </Stack>
        </Group>
        <Divider/>
    </Stack>
);