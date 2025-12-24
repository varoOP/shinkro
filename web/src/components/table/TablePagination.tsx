import { Group, Text, Select, ActionIcon } from "@mantine/core";
import { FaChevronLeft, FaChevronRight, FaAngleDoubleLeft, FaAngleDoubleRight } from "react-icons/fa";
import type { Table } from "@tanstack/react-table";
import { InfoTooltip } from "./InfoTooltip";

interface TablePaginationProps<T> {
    table: Table<T>;
    totalCount: number;
    pageSizeOptions?: Array<{ value: string; label: string }>;
}

export const TablePagination = <T,>({ 
    table, 
    totalCount,
    pageSizeOptions = [
        { value: "5", label: "5 entries" },
        { value: "10", label: "10 entries" },
        { value: "20", label: "20 entries" },
        { value: "50", label: "50 entries" },
    ]
}: TablePaginationProps<T>) => {
    const pagination = table.getState().pagination;
    const totalPages = Math.ceil(totalCount / pagination.pageSize);

    return (
        <Group justify="space-between">
            <Group>
                <Text size="sm">
                    Page <strong>{pagination.pageIndex + 1}</strong> of <strong>{totalPages}</strong>
                </Text>
                <Select
                    value={pagination.pageSize.toString()}
                    onChange={(value) =>
                        table.setPagination({ ...pagination, pageSize: parseInt(value || "20"), pageIndex: 0 })
                    }
                    data={pageSizeOptions}
                    style={{ width: 150 }}
                />
            </Group>
            <Group>
                <InfoTooltip label="First">
                    <ActionIcon
                        variant="outline"
                        onClick={() => table.setPageIndex(0)}
                        disabled={!table.getCanPreviousPage()}
                    >
                        <FaAngleDoubleLeft size={16} />
                    </ActionIcon>
                </InfoTooltip>
                <InfoTooltip label="Previous">
                    <ActionIcon
                        variant="outline"
                        onClick={() => table.previousPage()}
                        disabled={!table.getCanPreviousPage()}
                    >
                        <FaChevronLeft size={16} />
                    </ActionIcon>
                </InfoTooltip>
                <InfoTooltip label="Next">
                    <ActionIcon
                        variant="outline"
                        onClick={() => table.nextPage()}
                        disabled={!table.getCanNextPage()}
                    >
                        <FaChevronRight size={16} />
                    </ActionIcon>
                </InfoTooltip>
                <InfoTooltip label="Last">
                    <ActionIcon
                        variant="outline"
                        onClick={() => table.setPageIndex(table.getPageCount() - 1)}
                        disabled={!table.getCanNextPage()}
                    >
                        <FaAngleDoubleRight size={16} />
                    </ActionIcon>
                </InfoTooltip>
            </Group>
        </Group>
    );
};

