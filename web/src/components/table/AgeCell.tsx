import { Text } from "@mantine/core";
import { formatDistanceToNowStrict } from "date-fns";

interface AgeCellProps {
    timestamp?: string | Date | null;
}

export const AgeCell = ({ timestamp }: AgeCellProps) => {
    if (!timestamp) {
        return (
            <Text size="sm" ta="center">
                -
            </Text>
        );
    }

    const date = typeof timestamp === "string" ? new Date(timestamp) : timestamp;
    return (
        <Text size="sm" ta="center">
            {formatDistanceToNowStrict(date, { addSuffix: false })}
        </Text>
    );
};

