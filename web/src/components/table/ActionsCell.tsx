import { Group, ActionIcon } from "@mantine/core";
import { FaEye } from "react-icons/fa";
import { ConfirmDeleteIcon } from "@components/alerts/ConfirmDeleteIcon";
import { InfoTooltip } from "./InfoTooltip";

interface ActionsCellProps {
    onView: () => void;
    onDelete?: () => void;
    deleteTitle?: string;
    deleteMessage?: string;
    isDeleting?: boolean;
}

export const ActionsCell = ({ 
    onView, 
    onDelete, 
    deleteTitle = "Delete",
    deleteMessage = "Are you sure you want to delete this item?",
    isDeleting = false
}: ActionsCellProps) => {
    return (
        <Group gap="xs" justify="center">
            <InfoTooltip label="View details">
                <ActionIcon variant="outline" onClick={onView}>
                    <FaEye size={16} />
                </ActionIcon>
            </InfoTooltip>
            {onDelete && (
                <ConfirmDeleteIcon
                    onConfirm={onDelete}
                    title={deleteTitle}
                    message={deleteMessage}
                    loading={isDeleting}
                    variant="outline"
                />
            )}
        </Group>
    );
};

