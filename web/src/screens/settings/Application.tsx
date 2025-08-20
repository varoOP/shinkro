import {
    Stack,
    Table,
    Text,
} from "@mantine/core";
import {useQuery} from "@tanstack/react-query";
import {ConfigQueryOptions} from "@api/queries.ts";
import {SettingsSectionHeader} from "@screens/settings/components.tsx";

export const Application = () => {
    const {data: config} = useQuery(ConfigQueryOptions());

    const rows = config ? [
        {label: 'Binary', value: config.application || 'Not set'},
        {label: 'Config Directory', value: config.config_dir || 'Not set'},
        {label: 'Host', value: config.host || 'Not set'},
        {label: 'Port', value: config.port?.toString() || 'Not set'},
        {label: 'Base URL', value: config.base_url || 'Not set'},
        {label: 'Log Level', value: config.log_level || 'Not set'},
        {label: 'Log Path', value: config.log_path || 'Not set'},
        {label: 'Log Max Size (MB)', value: config.log_max_size?.toString() || 'Not set'},
        {label: 'Log Max Backups', value: config.log_max_backups?.toString() || 'Not set'},
        {label: 'Check for Updates', value: config.check_for_updates ? 'Enabled' : 'Disabled'},
        {label: 'Version', value: config.version || 'Not set'},
        {label: 'Commit', value: config.commit || 'Not set'},
        {label: 'Build Date', value: config.date || 'Not set'},
    ] : [];

    return (
        <main>
            <SettingsSectionHeader
                title="Application Settings"
                description="To change settings, edit config.toml found in the config directory and restart the application."
            />
            
            {config ? (
                <Stack mt="md">
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
                </Stack>
            ) : (
                <Stack mt="md">
                    <Text c="dimmed">Loading application settings...</Text>
                </Stack>
            )}
        </main>
    );
};