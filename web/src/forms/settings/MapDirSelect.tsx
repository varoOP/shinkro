import {
    Button, Divider,
    Group,
    Loader,
    Modal,
    Stack,
    Text,
} from "@mantine/core";
import {useState, useEffect} from "react";
import {useMutation} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient";
import {FileSystem} from "@app/types/FileSystem";
import {getFileExtension, getParentPath} from "@utils";
import {FaArrowLeft} from "react-icons/fa";

interface MapDirSelectProps {
    opened: boolean;
    onClose: () => void;
    onSelect: (path: string) => void;
}

export function MapDirSelect({opened, onClose, onSelect}: MapDirSelectProps) {
    const [currentPath, setCurrentPath] = useState<string | undefined>();
    const [entries, setEntries] = useState<FileSystem[]>([]);

    const {
        mutate: fetchDir,
        isPending: loading,
        error,
    } = useMutation({
        mutationFn: (path?: string) => APIClient.fs.listDirs(path ?? ""),
        onSuccess: (data) => {
            setEntries(data);
            if (data.length > 0) {
                const parent = getParentPath(data[0].path);
                setCurrentPath(parent);
            }
        },
    });

    // Reset to ConfigPath on open
    useEffect(() => {
        if (opened) {
            setCurrentPath(undefined);
            fetchDir(undefined);
        }
    }, [opened]);

    const goUp = () => {
        if (!currentPath || currentPath === "/") return;
        const parent = getParentPath(currentPath);
        if (parent !== currentPath) {
            setCurrentPath(parent);
            fetchDir(parent);
        }
    };

    const navigate = (dirPath: string) => {
        setCurrentPath(dirPath);
        fetchDir(dirPath);
    };

    return (
        <Modal opened={opened} onClose={onClose} title={<Text fw={700}>Select a .yml or .yaml File</Text>} size="lg">
            <Stack>
                <Text size="sm" fw={600}>
                    Current Path: <code>{currentPath || "[ConfigPath]"}</code>
                </Text>
                <Divider/>

                <Group justify="flex-start">
                    <Button onClick={goUp} variant="light" size="xs">
                        <FaArrowLeft/>
                    </Button>
                </Group>

                {loading ? (
                    <Loader/>
                ) : error ? (
                    <Text c="red">Failed to load directory.</Text>
                ) : entries.length > 0 ? (
                    entries.map((entry) => (
                        <Group key={entry.path} justify="space-between">
                            <Text>{entry.name}</Text>
                            {entry.is_dir ? (
                                <Button
                                    onClick={() => navigate(entry.path)}
                                    size="xs"
                                    w={80}
                                >
                                    Open
                                </Button>
                            ) : (
                                <Button
                                    onClick={() => {
                                        onSelect(entry.path);
                                        onClose();
                                    }}
                                    size="xs"
                                    color="blue"
                                    w={80}
                                    disabled={getFileExtension(entry.path) !== "yml" && getFileExtension(entry.path) !== "yaml"}
                                >
                                    Select
                                </Button>
                            )}
                        </Group>
                    ))
                ) : (
                    <Text>No files or directories found.</Text>
                )}
            </Stack>
        </Modal>
    );
}