import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {ConfigQueryOptions, LogQueryOptions} from "@api/queries.ts";
import {CenteredEmptyState, SettingsSectionHeader} from "@screens/settings/components.tsx";
import { Divider, Stack, Text, Select, Group, ActionIcon } from "@mantine/core";
import { useState } from "react";
import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {SettingsKeys} from "@api/query_keys.ts";
import {FaDownload} from "react-icons/fa";
import {baseUrl} from "@utils";

export const Logs = () => {
    const queryClient = useQueryClient();
    const {data} = useQuery(ConfigQueryOptions());

    const mutation = useMutation({
        mutationFn: (config: ConfigUpdate) => APIClient.config.update(config),
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: SettingsKeys.config()});
            displayNotification(
                {
                    title: "Log Level Updated",
                    type: "success"
                }
            );
        },
        onError: () => {
            queryClient.invalidateQueries({queryKey: SettingsKeys.config()});
        }
    });

    return (
        <main>
            <SettingsSectionHeader title={"Logs"} description={"Manage shinkro logs here."}/>
            
            {/* Log Settings Section */}
            {data && (
                <Stack mt={"md"}>
                    <Text fw={900} size={"xl"}>
                        Log Settings:
                    </Text>
                    <Group>
                        <Text fw={600} size={"md"}>
                            Log Level:
                        </Text>
                        <Select
                            size={"xs"}
                            data={[
                                {value: "DEBUG", label: "DEBUG"},
                                {value: "INFO", label: "INFO"},
                                {value: "ERROR", label: "ERROR"},
                                {value: "TRACE", label: "TRACE"},
                            ]}
                            defaultValue={data.log_level}
                            onChange={(value: LogLevel | string | null) => mutation.mutate({log_level: value ? value as LogLevel : ""})}
                        />
                    </Group>
                </Stack>
            )}
            
            <Divider mt={"md"}/>
            
            {/* Log Files Section */}
            <LogFiles/>
        </main>
    );
}

export const LogFiles = () => {
    const {data: logs} = useQuery(LogQueryOptions());
    const [isDownloading, setIsDownloading] = useState(false);
    const isEmpty = !logs || !(logs.length > 0)

    const handleDownload = async (filename: string) => {
        setIsDownloading(true);
        const response = await fetch(`${baseUrl()}api/fs/logs/${filename}`);
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = filename;
        link.click();
        URL.revokeObjectURL(url);
        setIsDownloading(false);
    }

    return (
        <main>
            <Text fw={900} size={"xl"} mt={"md"}>
                Log Files:
            </Text>
            {isEmpty ? (
                <CenteredEmptyState message={"No Logs Found"}/>
            ) : (
                <Stack>
                    {logs.map((log) => (
                        <Group key={log.name} align={"flex-start"} mt={"md"}>
                            <Stack gap={0}>
                                <Group>
                                    <Text fw={600}>
                                        {log.name}
                                    </Text>
                                    <ActionIcon
                                        loading={isDownloading}
                                        onClick={() => handleDownload(log.name)}
                                    >
                                        <FaDownload size={12}/>
                                    </ActionIcon>
                                </Group>
                                <Text size={"xs"} c={"dimmed"}>
                                    Size: {log.size_human}
                                </Text>
                                <Text size={"xs"} c={"dimmed"}>
                                    Last Modified: {new Date(log.modified_at).toLocaleString()}
                                </Text>
                            </Stack>
                        </Group>
                    ))}
                </Stack>
            )}
        </main>
    );
}