import {Alert, Button, Group, Modal} from "@mantine/core";
import {HiExclamationTriangle} from "react-icons/hi2";
import {useDisclosure} from "@mantine/hooks";

interface ConfirmDeleteButtonProps {
    onConfirm: () => void;
    title?: string;
    message?: string;
    confirmText?: string;
    cancelText?: string;
    loading?: boolean;
}

export const ConfirmDeleteButton = ({
                                        onConfirm,
                                        title = "Confirm Deletion",
                                        message = "",
                                        confirmText = "DELETE",
                                        cancelText = "CANCEL",
                                        loading = false,
                                    }: ConfirmDeleteButtonProps) => {
    const [opened, {open, close}] = useDisclosure(false);

    return (
        <>
            <Modal opened={opened} onClose={close} title={title}>
                <Alert
                    icon={<HiExclamationTriangle size="1.5rem"/>}
                    title="Warning"
                    color="red"
                    variant="outline"
                    mb="md"
                >
                    This action cannot be undone! Are you sure you want to delete? {message}
                </Alert>

                <Group justify="flex-end" mt="md">
                    <Button onClick={close}>
                        {cancelText}
                    </Button>
                    <Button
                        color="red"
                        loading={loading}
                        onClick={() => {
                            onConfirm();
                            close();
                        }}
                    >
                        {confirmText}
                    </Button>
                </Group>
            </Modal>

            <Button color="red" variant={"filled"} onClick={open}>
                {confirmText}
            </Button>
        </>
    );
};