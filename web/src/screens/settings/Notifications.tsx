import {Button, Stack, Paper, Text, Group, ActionIcon, Badge, Switch} from "@mantine/core";
import {CenteredEmptyState, SettingsSectionHeader} from "@screens/settings/components.tsx";
import {NotificationsQueryOptions} from "@api/queries.ts";
import {useQuery, useMutation, useQueryClient} from "@tanstack/react-query";
import {AddNotification} from "@forms/settings/AddNotification";
import {useState} from "react";
import {APIClient} from "@api/APIClient.ts";
import {NotificationKeys} from "@api/query_keys.ts";
import {displayNotification} from "@components/notifications";
import { FaEdit } from "react-icons/fa";
import {ConfirmDeleteIcon} from "@components/alerts/ConfirmDeleteIcon";

export const Notifications = () => {
    const [modalOpened, setModalOpened] = useState(false);
    const [editingNotification, setEditingNotification] = useState<ServiceNotification | undefined>();
    const queryClient = useQueryClient();
    const {data: notifications} = useQuery(NotificationsQueryOptions());

    const deleteMutation = useMutation({
        mutationFn: (id: number) => APIClient.notifications.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: NotificationKeys.lists()});
            displayNotification({
                title: "Notification Deleted",
                message: "The notification was successfully deleted",
                type: "success",
            });
        },
    });

    const updateMutation = useMutation({
        mutationFn: (notification: ServiceNotification) => APIClient.notifications.update(notification),
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: NotificationKeys.lists()});
            displayNotification({
                title: "Notification Updated",
                message: "The notification settings were updated successfully",
                type: "success",
            });
        },
    });

    const handleEdit = (notification: ServiceNotification) => {
        setEditingNotification(notification);
        setModalOpened(true);
    };

    const handleToggleEnabled = (notification: ServiceNotification) => {
        updateMutation.mutate({
            ...notification,
            enabled: !notification.enabled
        });
    };

    const handleCloseModal = () => {
        setModalOpened(false);
        setEditingNotification(undefined);
    };

    return (
        <main>
            <SettingsSectionHeader 
                title={"Notifications"} 
                description={"Manage your notification settings here."}
            />
            
            {notifications && notifications.length > 0 ? (
                <Stack mt="md">
                    {notifications.map((notification) => (
                        <Paper key={notification.id} p="md" withBorder>
                            <Group justify="space-between" align="flex-start">
                                <Stack gap="xs">
                                    <Group>
                                        <Text fw={500} size="lg">{notification.name}</Text>
                                        <Switch
                                            checked={notification.enabled}
                                            onChange={() => handleToggleEnabled(notification)}
                                            size="md"
                                        />
                                    </Group>
                                    <Text size="sm" c="dimmed">
                                        Type: {notification.type}
                                    </Text>
                                    <Group gap="xs">
                                        {notification.events.map((event) => (
                                            <Badge key={event} variant="light">
                                                {event}
                                            </Badge>
                                        ))}
                                    </Group>
                                </Stack>
                                <Group>
                                    <ActionIcon 
                                        variant="outline" 
                                        color="blue" 
                                        onClick={() => handleEdit(notification)}
                                    >
                                        <FaEdit size={18} />
                                    </ActionIcon>
                                    <ConfirmDeleteIcon
                                        onConfirm={() => deleteMutation.mutate(notification.id)}
                                        title="Delete Notification"
                                        message={`Are you sure you want to delete the notification "${notification.name}"?`}
                                        confirmText="DELETE"
                                    />
                                </Group>
                            </Group>
                        </Paper>
                    ))}
                </Stack>
            ) : (
                <CenteredEmptyState
                    message={"No Notification Settings Found"}
                />
            )}

            <Group justify="center" mt="md">
                <Button onClick={() => setModalOpened(true)}>
                    Add New Notification
                </Button>
            </Group>

            <AddNotification 
                opened={modalOpened}
                onClose={handleCloseModal}
                defaultValues={editingNotification}
            />
        </main>
    );
}