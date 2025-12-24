import { Badge } from "@mantine/core";
import { FaCheckCircle, FaTimesCircle, FaHourglassHalf } from "react-icons/fa";
import { ErrorTooltip } from "./ErrorTooltip";
import { InfoTooltip } from "./InfoTooltip";

export type StatusType = "SUCCESS" | "FAILED" | "PENDING" | boolean | null | undefined;

interface StatusBadgeProps {
    status: StatusType;
    errorType?: string;
    errorMessage?: string;
    variant?: "icon" | "badge";
}

export const StatusBadge = ({ status, errorType, errorMessage, variant = "icon" }: StatusBadgeProps) => {
    if (variant === "icon") {
        // Icon variant (like Plex payloads)
        if (status === "SUCCESS" || status === true) {
            return (
                <InfoTooltip label="Success">
                    <FaCheckCircle size={20} color="green" />
                </InfoTooltip>
            );
        }

        if (status === "FAILED" || status === false || status === null || status === undefined) {
            return (
                <ErrorTooltip errorType={errorType} errorMessage={errorMessage}>
                    <FaTimesCircle size={20} color="red" />
                </ErrorTooltip>
            );
        }

        // PENDING
        return (
            <InfoTooltip label="Pending">
                <FaHourglassHalf size={20} color="orange" />
            </InfoTooltip>
        );
    }

    // Badge variant (for Anime Updates)
    if (status === "SUCCESS") {
        return <Badge color="green" variant="outline">SUCCESS</Badge>;
    }

    if (status === "FAILED") {
        return (
            <ErrorTooltip errorType={errorType} errorMessage={errorMessage}>
                <Badge color="red" variant="outline">FAILED</Badge>
            </ErrorTooltip>
        );
    }

    if (status === "PENDING") {
        return <Badge color="yellow" variant="outline">PENDING</Badge>;
    }

    // Fallback for null/undefined
    return <Badge color="red" variant="outline">FAILED</Badge>;
};

