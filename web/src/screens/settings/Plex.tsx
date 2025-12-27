import {
    Button,
    Stack,
    Group,
    Table,
} from "@mantine/core";
import {useMutation, useSuspenseQuery, useQueryClient} from "@tanstack/react-query";
import {useEffect, useState, useRef, useLayoutEffect} from "react";
import {useDisclosure} from "@mantine/hooks";

import {APIClient} from "@api/APIClient";
import {PlexSettingsQueryOptions} from "@api/queries";
import {PlexSettingsKeys} from "@api/query_keys";
import {displayNotification} from "@components/notifications";
import {PlexSettings} from "@forms/settings/PlexSettings";
import {ConfirmDeleteButton} from "@components/alerts/ConfirmDeleteButton";
import {InfoTooltip} from "@components/table/InfoTooltip";
import {
    SettingsSectionHeader,
    StatusIndicator,
    CenteredEmptyState,
} from "./components";

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

export const Plex = () => {
    const queryClient = useQueryClient();
    const [isReachable, setIsReachable] = useState<boolean | null>(null);
    const [opened, {open, close}] = useDisclosure(false);
    const {data: settings} = useSuspenseQuery(PlexSettingsQueryOptions());
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
        <main>
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
                    <Stack mt={"md"}>
                        <StatusIndicator
                            label="Connection Status:"
                            status={isReachable}
                            loadStatus={isReachable === null}
                        />
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
                        <Group justify="flex-end">
                            <ConfirmDeleteButton
                                onConfirm={() => deleteMutation.mutate()}
                                message="All your Plex Media Server settings will be deleted."
                            />
                            <Button onClick={open}>EDIT SETTINGS</Button>
                        </Group>
                    </Stack>
                )}
            <PlexSettings
                opened={opened}
                onClose={close}
                defaultValues={settings}
            />
        </main>
    );
};