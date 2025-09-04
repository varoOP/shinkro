import {
    Stack,
    Group,
    Button,
    PasswordInput,
    Text,
    Paper,
    Table,
    Modal,
    TextInput,
    Switch,
    ActionIcon,
} from "@mantine/core";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {SettingsSectionHeader} from "@screens/settings/components.tsx";
import {useForm} from "@mantine/form";
import {AuthContext} from "@utils/Context";
import {useState} from "react";
import {FaTrash} from "react-icons/fa";

export const User = () => {
    const [auth] = AuthContext.use();
    const [createModalOpened, setCreateModalOpened] = useState(false);
    const queryClient = useQueryClient();

    const form = useForm({
        initialValues: {
            currentPassword: "",
            newPassword: "",
            confirmPassword: "",
        },
        validate: {
            currentPassword: (value) => (value ? null : "Current password is required"),
            newPassword: (value) => {
                if (!value) return "New password is required";
                if (value.length < 8) return "Password must be at least 8 characters";
                return null;
            },
            confirmPassword: (value, values) => {
                if (!value) return "Please confirm your new password";
                if (value !== values.newPassword) return "Passwords do not match";
                return null;
            },
        },
    });

    const createForm = useForm({
        initialValues: {
            username: "",
            password: "",
            admin: false,
        },
        validate: {
            username: (value) => (value ? null : "Username is required"),
            password: (value) => {
                if (!value) return "Password is required";
                if (value.length < 8) return "Password must be at least 8 characters";
                return null;
            },
        },
    });

    // Query for getting all users (only for admins)
    const {data: users} = useQuery({
        queryKey: ["users"],
        queryFn: () => APIClient.auth.getUsers(),
        enabled: auth.admin,
    });

    const updateMutation = useMutation({
        mutationFn: (values: typeof form.values) => 
            APIClient.auth.updateUser({
                username_current: auth.username,
                password_current: values.currentPassword,
                password_new: values.newPassword,
            }),
        onSuccess: () => {
            displayNotification({
                title: "Password Updated",
                message: "Your password has been successfully updated",
                type: "success",
            });
            form.reset();
        },
        onError: (error) => {
            displayNotification({
                title: "Password Update Failed",
                message: error.message || "Could not update password. Please check your current password.",
                type: "error",
            });
        },
    });

    const createMutation = useMutation({
        mutationFn: (values: typeof createForm.values) => 
            APIClient.auth.createUser(values),
        onSuccess: () => {
            displayNotification({
                title: "User Created",
                message: "User has been successfully created",
                type: "success",
            });
            createForm.reset();
            setCreateModalOpened(false);
            queryClient.invalidateQueries({queryKey: ["users"]});
        },
        onError: (error) => {
            displayNotification({
                title: "User Creation Failed",
                message: error.message || "Could not create user.",
                type: "error",
            });
        },
    });

    const deleteMutation = useMutation({
        mutationFn: (username: string) => APIClient.auth.deleteUser(username),
        onSuccess: () => {
            displayNotification({
                title: "User Deleted",
                message: "User has been successfully deleted",
                type: "success",
            });
            queryClient.invalidateQueries({queryKey: ["users"]});
        },
        onError: (error) => {
            displayNotification({
                title: "User Deletion Failed",
                message: error.message || "Could not delete user.",
                type: "error",
            });
        },
    });

    const handleSubmit = (values: typeof form.values) => {
        updateMutation.mutate(values);
    };

    const handleCreateUser = (values: typeof createForm.values) => {
        createMutation.mutate(values);
    };

    const handleDeleteUser = (username: string) => {
        if (confirm(`Are you sure you want to delete user "${username}"?`)) {
            deleteMutation.mutate(username);
        }
    };

    return (
        <main>
            <SettingsSectionHeader
                title="User Settings"
                description="Change your account password and manage users."
            />
            <Stack mt="md" maw={800} mx="auto">
                <Text size="lg">
                    Current User: <strong>{auth.username}</strong> {auth.admin && <Text span c="blue">(Admin)</Text>}
                </Text>
                
                {/* Password Change Section */}
                <Paper withBorder mt="sm" p="xl">
                    <Text size="md" fw={500} mb="md">Change Password</Text>
                    <form onSubmit={form.onSubmit(handleSubmit)}>
                        <Stack gap="md">
                            <PasswordInput
                                label="Current Password"
                                placeholder="Enter your current password"
                                {...form.getInputProps("currentPassword")}
                            />
                            
                            <PasswordInput
                                label="New Password"
                                placeholder="Enter your new password"
                                {...form.getInputProps("newPassword")}
                            />
                            
                            <PasswordInput
                                label="Confirm New Password"
                                placeholder="Confirm your new password"
                                {...form.getInputProps("confirmPassword")}
                            />
                            
                            <Group justify="center" mt="md">
                                <Button 
                                    type="submit" 
                                    loading={updateMutation.isPending}
                                >
                                    Change Password
                                </Button>
                            </Group>
                        </Stack>
                    </form>
                </Paper>

                {/* Admin Section - User Management */}
                {auth.admin && (
                    <Paper withBorder mt="md" p="xl">
                        <Group justify="space-between" mb="md">
                            <Text size="md" fw={500}>User Management</Text>
                            <Button onClick={() => setCreateModalOpened(true)}>
                                Create User
                            </Button>
                        </Group>
                        
                        {users && users.length > 0 ? (
                            <Table>
                                <Table.Thead>
                                    <Table.Tr>
                                        <Table.Th>Username</Table.Th>
                                        <Table.Th>Admin</Table.Th>
                                        <Table.Th>Actions</Table.Th>
                                    </Table.Tr>
                                </Table.Thead>
                                <Table.Tbody>
                                    {users.map((user) => (
                                        <Table.Tr key={user.id}>
                                            <Table.Td>{user.username}</Table.Td>
                                            <Table.Td>{user.admin ? "Yes" : "No"}</Table.Td>
                                            <Table.Td>
                                                <Group gap="xs">
                                                    <ActionIcon
                                                        color="red"
                                                        onClick={() => handleDeleteUser(user.username)}
                                                        disabled={user.username === auth.username}
                                                    >
                                                        <FaTrash />
                                                    </ActionIcon>
                                                </Group>
                                            </Table.Td>
                                        </Table.Tr>
                                    ))}
                                </Table.Tbody>
                            </Table>
                        ) : (
                            <Text c="dimmed">No users found.</Text>
                        )}
                    </Paper>
                )}

                {/* Create User Modal */}
                <Modal
                    opened={createModalOpened}
                    onClose={() => setCreateModalOpened(false)}
                    title="Create New User"
                >
                    <form onSubmit={createForm.onSubmit(handleCreateUser)}>
                        <Stack gap="md">
                            <TextInput
                                label="Username"
                                placeholder="Enter username"
                                {...createForm.getInputProps("username")}
                            />
                            
                            <PasswordInput
                                label="Password"
                                placeholder="Enter password"
                                {...createForm.getInputProps("password")}
                            />
                            
                            <Switch
                                label="Admin"
                                description="Give this user admin privileges"
                                {...createForm.getInputProps("admin", { type: "checkbox" })}
                            />
                            
                            <Group justify="flex-end" mt="md">
                                <Button 
                                    variant="subtle" 
                                    onClick={() => setCreateModalOpened(false)}
                                >
                                    Cancel
                                </Button>
                                <Button 
                                    type="submit" 
                                    loading={createMutation.isPending}
                                >
                                    Create User
                                </Button>
                            </Group>
                        </Stack>
                    </form>
                </Modal>
            </Stack>
        </main>
    );
};