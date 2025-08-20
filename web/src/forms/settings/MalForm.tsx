import {useForm} from "@mantine/form";
import {Button, CopyButton, Group, Modal, PasswordInput, Text, Code} from "@mantine/core";
import {MalAuth} from "@app/types/MalAuth";
import {useMutation} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";

interface Props {
    opened: boolean;
    onClose: () => void;
    loading: boolean;
    setLoading: (loading: boolean) => void;
}

export const MalForm = ({opened, onClose, loading, setLoading}: Props) => {
    const appRedirectURL = `${window.location.origin}${window.APP.baseUrl}malauth/callback`

    const form = useForm<MalAuth>({
        initialValues: {
            clientID: "",
            clientSecret: "",
        },
        validate: {
            clientID: (value: string) => (value ? null : "Client ID is required"),
            clientSecret: (value: string) => (value ? null : "Client Secret is required"),
        },
    });

    const mutation = useMutation({
        mutationFn: APIClient.malauth.start,
        onSuccess: (data) => {
            const url = data?.url;
            if (url) {
                setLoading(true);
                window.open(url, "_blank");
            }
        },
        onError: (error) => {
            displayNotification({
                title: "MyAnimeList Authentication Failed",
                message: error.message || "Could not start MyAnimeList authentication",
                type: "error",
            });
        },
    })

    const handleFormSubmit = (mal: MalAuth) => {
        mutation.mutate(mal);
        form.reset();
    };

    return (
        <Modal opened={opened} onClose={onClose} title={"Login to MyAnimeList"}>
            <form onSubmit={form.onSubmit(handleFormSubmit)}>
                <PasswordInput
                    label="Client ID"
                    placeholder="Enter Client ID"
                    {...form.getInputProps("clientID")}
                />
                <PasswordInput
                    label="Client Secret"
                    placeholder="Enter Client Secret"
                    {...form.getInputProps("clientSecret")}
                    mt={"md"}
                />
                <Group mt={"md"}>
                    <Text size="sm">App Redirect URL</Text>
                    <Code c="dimmed">{appRedirectURL}</Code>
                </Group>
                <Group justify={"center"} align={"flex-end"}>
                    <CopyButton value={appRedirectURL}>
                        {({copied, copy}) => (
                                <Button color={copied ? 'teal' : 'mal'} onClick={copy}>
                                    {copied ? 'COPIED URL' : 'COPY APP REDIRECT URL'}
                                </Button>
                        )}
                    </CopyButton>
                    <Button type="submit" mt={"md"} loading={loading}>
                        LOGIN
                    </Button>
                </Group>
            </form>
        </Modal>
    )
}