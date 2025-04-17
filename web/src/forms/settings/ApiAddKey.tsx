import {Button, Modal, TextInput, Group} from "@mantine/core";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useForm} from "@mantine/form";
import {APIClient} from "@api/APIClient.ts";
import {ApiKeys} from "@api/query_keys.ts";
import {displayNotification} from "@components/notifications";

interface AddKeyProps {
    opened: boolean;
    onClose: () => void;
}

export const ApiAddKey = ({opened, onClose}: AddKeyProps) => {
    const queryClient = useQueryClient();
    const form = useForm(
        {
            initialValues: {
                name: "",
                scopes: []
            },
            validate: {
                name: (value: string) => (value ? null : "Required"),
            },
        }
    );

    const mutation = useMutation({
        mutationFn: (key: APIKey) => APIClient.apikeys.create(key),
        onSuccess: (_, key) => {
            queryClient.invalidateQueries({queryKey: ApiKeys.lists()});
            displayNotification(
                {
                    title: "API Key Created",
                    message: `API key ${key.name} was created`,
                    type: "success",
                }
            );
            onClose();
        },
    });

    const handleFormSubmit = (values: unknown) => {
        mutation.mutate(values as APIKey);
        form.reset();
    }

    return (
        <Modal opened={opened} onClose={onClose} title={"Add API Key"}>
            <form onSubmit={form.onSubmit(handleFormSubmit)}>
                <TextInput
                    label="Name"
                    placeholder="Enter API key name"
                    {...form.getInputProps("name")}
                />
                <Group justify={"flex-end"} mt={"md"}>
                    <Button variant="default" onClick={onClose}>
                        CANCEL
                    </Button>
                    <Button type={"submit"} justify={"flex-end"}>
                        CREATE
                    </Button>
                </Group>
            </form>
        </Modal>
    )
}