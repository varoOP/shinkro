import { Badge, Anchor, Button, Text } from "@mantine/core";
import { FaExternalLinkAlt } from "react-icons/fa";
import { getSourceLink } from "@utils/sourceLinks";

interface SourceBadgeProps {
    sourceDB?: string;
    sourceID?: number;
    variant?: "transparent" | "outline" | "button";
}

export const SourceBadge = ({ sourceDB, sourceID, variant = "transparent" }: SourceBadgeProps) => {
    if (!sourceDB || !sourceID) {
        if (variant === "button") {
            return (
                <Text size="sm" c="dimmed">
                    -
                </Text>
            );
        }
        return (
            <Badge variant={variant} color="gray">
                Unknown
            </Badge>
        );
    }

    const link = getSourceLink(sourceDB, sourceID);
    const displayText = `${sourceDB.toUpperCase()}: ${sourceID}`;

    if (variant === "button") {
        if (link) {
            return (
                <Button
                    component="a"
                    href={link}
                    target="_blank"
                    rel="noopener noreferrer"
                    variant="subtle"
                    size="xs"
                    rightSection={<FaExternalLinkAlt size={12} />}
                >
                    {displayText}
                </Button>
            );
        }
        return (
            <Text size="sm" c="dimmed">
                {displayText}
            </Text>
        );
    }

    if (link) {
        return (
            <Anchor href={link} target="_blank" rel="noopener noreferrer" underline="never">
                <Badge variant={variant} color={sourceDB.toUpperCase() === "TVDB" ? "blue" : "gray"}>
                    {displayText}
                </Badge>
            </Anchor>
        );
    }

    return (
        <Badge variant={variant} color="gray">
            {displayText}
        </Badge>
    );
};

