import { Tooltip, Paper, Text } from "@mantine/core";
import { ReactNode } from "react";

interface InfoTooltipProps {
    label: string;
    children: ReactNode;
}

export const InfoTooltip = ({ label, children }: InfoTooltipProps) => {
    const tooltipContent = (
        <Paper p="sm" style={{ width: "auto", maxWidth: 300 }}>
            <Text size="xs" c="dimmed" fw={600}>
                {label}
            </Text>
        </Paper>
    );

    return (
        <Tooltip
            label={tooltipContent}
            position="top"
            withArrow
            styles={{
                tooltip: {
                    padding: 0,
                    backgroundColor: "transparent",
                }
            }}
        >
            {children}
        </Tooltip>
    );
};

