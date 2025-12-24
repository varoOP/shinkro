import { Tooltip, Stack, Text, Divider } from "@mantine/core";
import { ReactNode } from "react";

interface ErrorTooltipProps {
    errorType?: string;
    errorMessage?: string;
    children: ReactNode;
}

export const ErrorTooltip = ({ errorType, errorMessage, children }: ErrorTooltipProps) => {
    // Build the tooltip content
    const hasError = errorType || errorMessage;
    
    if (!hasError) {
        // If no error info, just return children without tooltip
        return <>{children}</>;
    }

    // Format error type: replace underscores with spaces
    const formattedErrorType = errorType ? errorType.replace(/_/g, " ") : undefined;

    const tooltipContent = (
        <div>
            <Stack gap="xs">
                {formattedErrorType && (
                    <Text size="sm" fw={600} c="red.6">
                        {formattedErrorType}
                    </Text>
                )}
                {errorType && errorMessage && <Divider size="xs" />}
                {errorMessage && (
                    <Text 
                        size="xs" 
                        fw={600}
                        style={{
                            wordBreak: "normal",
                            overflowWrap: "break-word",
                            whiteSpace: "pre-wrap"
                        }}
                    >
                        {errorMessage}
                    </Text>
                )}
            </Stack>
        </div>
    );

    return (
        <Tooltip
            label={tooltipContent}
            position="top"
            multiline={true}
            withArrow={true}
            maw={320}
            arrowSize={8}
            arrowPosition="center"
        >
            {children}
        </Tooltip>
    );
};

