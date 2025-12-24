import { Modal, Stack, Group, Button, ScrollArea, Code } from "@mantine/core";
import { FaCopy } from "react-icons/fa";
import { displayNotification } from "@components/notifications";

interface ViewDetailsModalProps<T> {
    opened: boolean;
    onClose: () => void;
    title: string;
    data: T | null;
}

export const ViewDetailsModal = <T,>({ opened, onClose, title, data }: ViewDetailsModalProps<T>) => {
    const handleCopy = () => {
        if (data) {
            const json = JSON.stringify(data, null, 2);
            navigator.clipboard.writeText(json);
            displayNotification({
                title: "Copied",
                message: "Data copied to clipboard",
                type: "success",
            });
        }
    };

    return (
        <Modal opened={opened} onClose={onClose} title={title} size="xl">
            {data && (
                <Stack>
                    <Group justify="flex-end">
                        <Button
                            leftSection={<FaCopy size={14} />}
                            variant="outline"
                            onClick={handleCopy}
                        >
                            Copy JSON
                        </Button>
                    </Group>
                    <ScrollArea h={500}>
                        <Code block>{JSON.stringify(data, null, 2)}</Code>
                    </ScrollArea>
                </Stack>
            )}
        </Modal>
    );
};

