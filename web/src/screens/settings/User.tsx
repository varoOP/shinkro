import {
    Stack,
    Group,
    Button,
    PasswordInput,
    Text,
    Paper,
} from "@mantine/core";
import {useMutation} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {SettingsSectionHeader} from "@screens/settings/components.tsx";
import {useForm} from "@mantine/form";
import {AuthContext} from "@utils/Context";

export const User = () => {
    const [auth] = AuthContext.use();

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

    const mutation = useMutation({
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

    const handleSubmit = (values: typeof form.values) => {
        mutation.mutate(values);
    };

    return (
        <main>
            <SettingsSectionHeader
                title="User Settings"
                description="Change your account password here."
            />
            <Stack mt="md" maw={800} mx="auto">
                <Text size="lg">
                    Current User: <strong>{auth.username}</strong>
                </Text>
                <Paper withBorder mt="sm" p="xl">
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
                                loading={mutation.isPending}
                            >
                                Change Password
                            </Button>
                        </Group>
                    </Stack>
                </form>
                </Paper>
            </Stack>
        </main>
    );
};