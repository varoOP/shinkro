import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import {
    useReactTable,
    getCoreRowModel,
    flexRender,
    type ColumnDef,
    type ColumnFilter,
    type PaginationState,
} from "@tanstack/react-table";
import {
    Container,
    Paper,
    Stack,
    Title,
    Table,
    Badge,
    Group,
    Text,
    Button,
    Select,
    TextInput,
    Modal,
    Code,
    ScrollArea,
    ActionIcon,
    Tooltip,
} from "@mantine/core";
import { formatDistanceToNowStrict } from "date-fns";
import { FaEye, FaCheckCircle, FaTimesCircle, FaCopy, FaChevronLeft, FaChevronRight, FaAngleDoubleLeft, FaAngleDoubleRight } from "react-icons/fa";
import { useDisclosure } from "@mantine/hooks";
import { AuthContext } from "@utils/Context";
import { Navigate } from "@tanstack/react-router";
import { plexPayloadsQueryOptions } from "@api/queries";
import type { PlexPayloadListItem } from "@app/types/Plex";
import { displayNotification } from "@components/notifications";
import { formatEventName } from "@utils";

export const PlexPayloads = () => {
    const isLoggedIn = AuthContext.useSelector((s) => s.isLoggedIn);
    if (!isLoggedIn) {
        return <Navigate to="/login" />;
    }

    const [pagination, setPagination] = useState<PaginationState>({
        pageIndex: 0,
        pageSize: 10,
    });

    const [columnFilters, setColumnFilters] = useState<ColumnFilter[]>([]);
    const [selectedPayload, setSelectedPayload] = useState<PlexPayloadListItem | null>(null);
    const [opened, { open, close }] = useDisclosure(false);

    const { isLoading, error, data } = useQuery(
        plexPayloadsQueryOptions(pagination.pageIndex, pagination.pageSize, columnFilters)
    );

    const columns = useMemo<ColumnDef<PlexPayloadListItem>[]>(
        () => [
            {
                header: "Age",
                accessorKey: "plex.timestamp",
                cell: ({ row }) => {
                    const timestamp = row.original.plex.timestamp;
                    return (
                        <Text size="sm" ta="center">
                            {timestamp ? formatDistanceToNowStrict(new Date(timestamp), { addSuffix: false }) : "-"}
                        </Text>
                    );
                },
            },
            {
                header: "Plex Payload",
                accessorKey: "payload",
                id: "payload",
                cell: ({ row }) => {
                    const plex = row.original.plex;
                    const title = plex.Metadata?.grandparentTitle || plex.Metadata?.title || "Unknown";
                    const library = plex.Metadata?.librarySectionTitle || "";
                    const event = plex.event || "";

                    let season = "";
                    let episode = "";
                    const isMovie = plex.Metadata?.type === "movie";
                    if (plex.Metadata && !isMovie) {
                        if (plex.Metadata.parentIndex !== undefined && plex.Metadata.parentIndex > 0) {
                            season = `S${plex.Metadata.parentIndex}`;
                        }
                        if (plex.Metadata.index !== undefined && plex.Metadata.index > 0) {
                            episode = `E${plex.Metadata.index}`;
                        }
                    }

                    return (
                        <Stack gap={4} style={{ maxWidth: "500px" }}>
                            <Tooltip label={title} disabled={title.length <= 50}>
                                <Text 
                                    fw={600} 
                                    size="sm" 
                                    style={{ 
                                        overflow: "hidden",
                                        textOverflow: "ellipsis",
                                        whiteSpace: "nowrap"
                                    }}
                                >
                                    {title}
                                </Text>
                            </Tooltip>
                            <Group gap="xs">
                                {library && (
                                    <Badge size="xs" variant="transparent" color="gray">
                                        {library}
                                    </Badge>
                                )}
                                {event && (
                                    <Badge size="xs" variant="transparent" color="plex">
                                        {formatEventName(event)}
                                    </Badge>
                                )}
                                {season && episode && (
                                    <Badge size="xs" variant="transparent">
                                        {season} {episode}
                                    </Badge>
                                )}
                            </Group>
                        </Stack>
                    );
                },
                meta: {
                    filterVariant: "search",
                },
            },
            {
                header: "View Payload",
                accessorKey: "view",
                cell: ({ row }) => {
                    return (
                        <Tooltip label="View full payload">
                            <ActionIcon
                                variant="light"
                                onClick={() => {
                                    setSelectedPayload(row.original);
                                    open();
                                }}
                            >
                                <FaEye size={16} />
                            </ActionIcon>
                        </Tooltip>
                    );
                },
            },
            {
                header: "Processing Status",
                accessorKey: "status",
                id: "status",
                cell: ({ row }) => {
                    // Read from consolidated fields
                    const plex = row.original.plex;
                    const plexSuccess = plex?.success;
                    const errorType = plex?.errorType;
                    const errorMsg = plex?.errorMsg;
                    
                    if (plexSuccess === undefined || plexSuccess === null) {
                        return (
                            <Tooltip label="Failed">
                                <FaTimesCircle size={20} color="red" />
                            </Tooltip>
                        );
                    }
                    if (plexSuccess === true) {
                        return (
                            <Tooltip label="Success">
                                <FaCheckCircle size={20} color="green" />
                            </Tooltip>
                        );
                    }
                    // Build tooltip with errorType and errorMessage
                    const tooltipLabel = errorType && errorMsg
                        ? `${errorType}: ${errorMsg}`
                        : errorMsg || errorType || "Failed";
                    return (
                        <Tooltip label={tooltipLabel}>
                            <FaTimesCircle size={20} color="red" />
                        </Tooltip>
                    );
                },
            },
            {
                header: "Source Type",
                accessorKey: "plex.source",
                id: "source",
                cell: ({ row }) => {
                    const source = row.original.plex.source || "";
                    return (
                        <Badge variant="transparent" color={source === "Plex Webhook" ? "plex" : "blue"}>
                            {source || "Unknown"}
                        </Badge>
                    );
                },
            },
        ],
        []
    );

    const tableInstance = useReactTable({
        columns,
        data: data?.data || [],
        getCoreRowModel: getCoreRowModel(),
        manualFiltering: true,
        manualPagination: true,
        rowCount: data?.count || 0,
        state: {
            columnFilters,
            pagination,
        },
        onPaginationChange: setPagination,
        onColumnFiltersChange: setColumnFilters,
    });

    const handleCopyPayload = () => {
        if (selectedPayload) {
            const payloadJson = JSON.stringify(selectedPayload, null, 2);
            navigator.clipboard.writeText(payloadJson);
            displayNotification({
                title: "Copied",
                message: "Payload copied to clipboard",
                type: "success",
            });
        }
    };

    return (
        <Container size={1200} px="md" component="main">
            <Stack gap="md" mt="md">
            <Title order={2}>Plex Payloads</Title>
                {/* Filters */}
                <Paper mt="md" withBorder p="md">
                <Stack gap="md">
                <Group>
                    <TextInput
                        placeholder="Search by title"
                        value={(columnFilters.find((f) => f.id === "payload")?.value as string) || ""}
                        onChange={(e) => {
                            const filter = columnFilters.filter((f) => f.id !== "payload");
                            if (e.target.value) {
                                filter.push({ id: "payload", value: e.target.value });
                            }
                            setColumnFilters(filter);
                        }}
                        style={{ flex: 1 }}
                    />
                    <Select
                        placeholder="Event"
                        data={[
                            { value: "media.scrobble", label: "Scrobble" },
                            { value: "media.rate", label: "Rate" },
                        ]}
                        value={(columnFilters.find((f) => f.id === "event")?.value as string) || ""}
                        onChange={(value) => {
                            const filter = columnFilters.filter((f) => f.id !== "event");
                            if (value) {
                                filter.push({ id: "event", value });
                            }
                            setColumnFilters(filter);
                        }}
                        clearable
                    />
                    <Select
                        placeholder="Source"
                        data={[
                            { value: "Plex Webhook", label: "Plex Webhook" },
                            { value: "Tautulli", label: "Tautulli" },
                        ]}
                        value={(columnFilters.find((f) => f.id === "source")?.value as string) || ""}
                        onChange={(value) => {
                            const filter = columnFilters.filter((f) => f.id !== "source");
                            if (value) {
                                filter.push({ id: "source", value });
                            }
                            setColumnFilters(filter);
                        }}
                        clearable
                    />
                    <Select
                        placeholder="Status"
                        data={[
                            { value: "success", label: "Success" },
                            { value: "failed", label: "Failed" },
                        ]}
                        value={(columnFilters.find((f) => f.id === "status")?.value as string) || ""}
                        onChange={(value) => {
                            const filter = columnFilters.filter((f) => f.id !== "status");
                            if (value) {
                                filter.push({ id: "status", value });
                            }
                            setColumnFilters(filter);
                        }}
                        clearable
                    />
                </Group>

                {/* Table */}
                {error ? (
                    <Text c="red">Error loading payloads</Text>
                ) : (
                    <Table.ScrollContainer minWidth={800}>
                        <Table highlightOnHover>
                            <Table.Thead>
                                {tableInstance.getHeaderGroups().map((headerGroup) => (
                                    <Table.Tr key={headerGroup.id}>
                                        {headerGroup.headers.map((header) => {
                                            const shouldCenter = header.column.id !== "payload";
                                            return (
                                                <Table.Th key={header.id} style={{ textAlign: shouldCenter ? "center" : "left" }}>
                                                    {header.isPlaceholder
                                                        ? null
                                                        : flexRender(header.column.columnDef.header, header.getContext())}
                                                </Table.Th>
                                            );
                                        })}
                                    </Table.Tr>
                                ))}
                            </Table.Thead>
                            <Table.Tbody>
                                {isLoading ? (
                                    <Table.Tr>
                                        <Table.Td colSpan={columns.length} style={{ textAlign: "center" }}>
                                            <Text>Loading...</Text>
                                        </Table.Td>
                                    </Table.Tr>
                                ) : tableInstance.getRowModel().rows.length === 0 ? (
                                    <Table.Tr>
                                        <Table.Td colSpan={columns.length} style={{ textAlign: "center" }}>
                                            <Text c="dimmed">No payloads found</Text>
                                        </Table.Td>
                                    </Table.Tr>
                                ) : (
                                    tableInstance.getRowModel().rows.map((row) => (
                                        <Table.Tr key={row.id}>
                                            {row.getVisibleCells().map((cell) => {
                                                const shouldCenter = cell.column.id !== "payload";
                                                const isSourceColumn = cell.column.id === "source";
                                                return (
                                                    <Table.Td 
                                                        key={cell.id} 
                                                        style={{ 
                                                            textAlign: shouldCenter ? "center" : "left",
                                                            whiteSpace: isSourceColumn ? "nowrap" : "normal"
                                                        }}
                                                    >
                                                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                                    </Table.Td>
                                                );
                                            })}
                                        </Table.Tr>
                                    ))
                                )}
                            </Table.Tbody>
                        </Table>
                    </Table.ScrollContainer>
                )}

                {/* Pagination */}
                <Group justify="space-between">
                    <Group>
                        <Text size="sm">
                            Page <strong>{pagination.pageIndex + 1}</strong> of{" "}
                            <strong>{Math.ceil((data?.count || 0) / pagination.pageSize)}</strong>
                        </Text>
                        <Select
                            value={pagination.pageSize.toString()}
                            onChange={(value) =>
                                setPagination({ ...pagination, pageSize: parseInt(value || "20"), pageIndex: 0 })
                            }
                            data={[
                                { value: "5", label: "5 entries" },
                                { value: "10", label: "10 entries" },
                                { value: "20", label: "20 entries" },
                                { value: "50", label: "50 entries" },
                            ]}
                            style={{ width: 150 }}
                        />
                    </Group>
                    <Group>
                        <Tooltip label="First">
                            <ActionIcon
                                variant="outline"
                                onClick={() => tableInstance.setPageIndex(0)}
                                disabled={!tableInstance.getCanPreviousPage()}
                            >
                                <FaAngleDoubleLeft size={16} />
                            </ActionIcon>
                        </Tooltip>
                        <Tooltip label="Previous">
                            <ActionIcon
                                variant="outline"
                                onClick={() => tableInstance.previousPage()}
                                disabled={!tableInstance.getCanPreviousPage()}
                            >
                                <FaChevronLeft size={16} />
                            </ActionIcon>
                        </Tooltip>
                        <Tooltip label="Next">
                            <ActionIcon
                                variant="outline"
                                onClick={() => tableInstance.nextPage()}
                                disabled={!tableInstance.getCanNextPage()}
                            >
                                <FaChevronRight size={16} />
                            </ActionIcon>
                        </Tooltip>
                        <Tooltip label="Last">
                            <ActionIcon
                                variant="outline"
                                onClick={() => tableInstance.setPageIndex(tableInstance.getPageCount() - 1)}
                                disabled={!tableInstance.getCanNextPage()}
                            >
                                <FaAngleDoubleRight size={16} />
                            </ActionIcon>
                        </Tooltip>
                    </Group>
                </Group>
            </Stack>
            </Paper>
            </Stack>

            {/* Payload View Modal */}
            <Modal
                opened={opened}
                onClose={close}
                title="Plex Payload"
                size="xl"
            >
                {selectedPayload && (
                    <Stack>
                        <Group justify="flex-end">
                            <Button
                                leftSection={<FaCopy size={14} />}
                                variant="outline"
                                onClick={handleCopyPayload}
                            >
                                Copy JSON
                            </Button>
                        </Group>
                        <ScrollArea h={500}>
                            <Code block>{JSON.stringify(selectedPayload, null, 2)}</Code>
                        </ScrollArea>
                    </Stack>
                )}
            </Modal>
        </Container>
    );
};

