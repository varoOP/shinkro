import {useForm} from "@mantine/form";
import {Button, Group, Modal, PasswordInput, Text, Code} from "@mantine/core";
import {MalAuth} from "@app/types/MalAuth";
import {useMutation} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {useState} from "react";
import {CopyTextToClipboard} from "@utils/index";

interface Props {
    opened: boolean;
    onClose: () => void;
    loading: boolean;
    setLoading: (loading: boolean) => void;
}

export const MalForm = ({opened, onClose, loading, setLoading}: Props) => {
    const [copied, setCopied] = useState(false);
    const [copyError, setCopyError] = useState(false);
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

    const handleCopy = async () => {
        try {
            await CopyTextToClipboard(appRedirectURL);
            setCopied(true);
            setCopyError(false);
            setTimeout(() => setCopied(false), 2000);
        } catch (error) {
            console.error('Copy failed:', error);
            setCopyError(true);
            displayNotification({
                title: "Copy Failed",
                message: "Please manually copy the URL from the text above. Clipboard access may be restricted over HTTP.",
                type: "info",
            });
        }
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
                    <Button 
                        color={copied ? 'teal' : copyError ? 'red' : 'mal'} 
                        onClick={handleCopy}
                    >
                        {copied ? 'COPIED URL' : copyError ? 'COPY FAILED' : 'COPY APP REDIRECT URL'}
                    </Button>
                    <Button type="submit" mt={"md"} loading={loading}>
                        LOGIN
                    </Button>
                </Group>
            </form>
        </Modal>
    )
}