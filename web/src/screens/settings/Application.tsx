import {
    Stack,
    Table,
} from "@mantine/core";
import {useSuspenseQuery} from "@tanstack/react-query";
import {useRef, useState, useEffect, useLayoutEffect} from "react";
import {ConfigQueryOptions} from "@api/queries.ts";
import {SettingsSectionHeader} from "@screens/settings/components.tsx";
import {InfoTooltip} from "@components/table/InfoTooltip";

const TruncatedText = ({ value }: { value: string }) => {
    const textRef = useRef<HTMLDivElement>(null);
    const [isTruncated, setIsTruncated] = useState(false);

    const checkTruncation = (element: HTMLDivElement | null) => {
        if (element) {
            // Use requestAnimationFrame to ensure layout is complete
            requestAnimationFrame(() => {
                const isOverflowing = element.scrollWidth > element.clientWidth;
                setIsTruncated(isOverflowing);
            });
        }
    };

    const setRef = (element: HTMLDivElement | null) => {
        textRef.current = element;
        checkTruncation(element);
    };

    useLayoutEffect(() => {
        checkTruncation(textRef.current);
    }, [value]);

    useEffect(() => {
        const handleResize = () => {
            checkTruncation(textRef.current);
        };
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    const textElement = (
        <div
            ref={setRef}
            style={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                display: 'block',
                width: '100%',
                minWidth: 0,
            }}
        >
            {value}
        </div>
    );

    if (isTruncated) {
        return <InfoTooltip label={value}>{textElement}</InfoTooltip>;
    }

    return textElement;
};

export const Application = () => {
    const {data: config} = useSuspenseQuery(ConfigQueryOptions());

    const rows = [
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
    ];

    return (
        <main>
            <SettingsSectionHeader
                title="Application"
                description="To change settings, edit config.toml found in the config directory and restart the application."
            />
            
            <Stack mt="md">
                <div style={{ overflowX: 'auto', width: '100%' }}>
                    <Table variant="vertical" verticalSpacing="sm" withTableBorder style={{ minWidth: '100%' }}>
                        <Table.Tbody>
                            {rows.map((row, index) => (
                                <Table.Tr key={index}>
                                    <Table.Th w={180} style={{ whiteSpace: 'nowrap' }}>{row.label}</Table.Th>
                                    <Table.Td style={{ maxWidth: 0, width: '100%' }}>
                                        <TruncatedText value={row.value} />
                                    </Table.Td>
                                </Table.Tr>
                            ))}
                        </Table.Tbody>
                    </Table>
                </div>
            </Stack>
        </main>
    );
};