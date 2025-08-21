import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {ConfigQueryOptions, LogQueryOptions} from "@api/queries.ts";
import {CenteredEmptyState, SettingsSectionHeader} from "@screens/settings/components.tsx";
import { Stack, Text, Select } from "@mantine/core";
import { useState } from "react";
import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {SettingsKeys} from "@api/query_keys.ts";

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

    // Always set log level to TRACE on mount
    useState(() => {
        if (data && data.log_level !== "TRACE") {
            mutation.mutate({ log_level: "TRACE" });
        }
    });
    return (
        <main>
            <SettingsSectionHeader title={"Logs"} description={"Manage shinkro logs here."}/>
            <LogFiles/>
        </main>
    );
}

// ...existing code...
import { LogViewer } from "./LogViewer";

export const LogFiles = () => {
    const { data: logs } = useQuery(LogQueryOptions());
    const [selected, setSelected] = useState<string | null>(null);
    const isEmpty = !logs || !(logs.length > 0);

    return (
        <main>
            <Text fw={900} size={"xl"} mt={"md"}>
                Log Viewer:
            </Text>
            {isEmpty ? (
                <CenteredEmptyState message={"No Logs Found"} />
            ) : (
                <Stack>
                    <Select
                        label="Log file"
                        data={logs.map((log) => ({ value: log.name, label: log.name }))}
                        value={selected}
                        onChange={setSelected}
                        placeholder="Select a log file"
                        size="xs"
                        mb="sm"
                    />
                    {selected && <LogViewer filename={selected} />}
                </Stack>
            )}
        </main>
    );
}