import { Tooltip, Text } from "@mantine/core";
import { ReactNode } from "react";

interface InfoTooltipProps {
    label: string;
    children: ReactNode;
}

export const InfoTooltip = ({ label, children }: InfoTooltipProps) => {
    const tooltipContent = (
            <Text size="xs" fw={600}>
                {label}
            </Text> 
    );

    return (
        <Tooltip
            label={tooltipContent}
            position="top"
            maw={320}
            multiline={true}
            withArrow={true}
            arrowSize={8}
            arrowPosition="center"
        >
            {children}
        </Tooltip>
    );
};

