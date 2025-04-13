import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {useState, useEffect} from "react";
import {Divider, Paper, Text, Title, Stack, PasswordInput, Button, Loader, Flex, Group} from "@mantine/core";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {useForm} from "@mantine/form";
import {MalAuth} from "@app/types/MalAuth";
import {MalQueryOptions} from "@api/queries.ts";
import {MalAuthKeys} from "@api/query_keys.ts";
import {ConfirmDeleteButton} from "@components/alerts/ConfirmDeleteButton";

export const Mal = () => {
    const queryClient = useQueryClient();
    const {data: malauth, isLoading} = useQuery(MalQueryOptions());
    const [testSucess, setTestSuccess] = useState<boolean | null>(null);
    ;

    const isEmptySettings = !malauth || Object.keys(malauth).length === 0;

    useEffect(() => {
        if (!isEmptySettings) {
            APIClient.malauth.test()
                .then(() => setTestSuccess(true))
                .catch(() => setTestSuccess(false));
        } else {
            setTestSuccess(null);
        }
    }, [malauth]);

    useEffect(() => {
        const handleMessage = (event: MessageEvent) => {
            if (event.origin !== window.location.origin) {
                return;
            }
            if (event.data?.type === "mal-auth-success") {
                queryClient.invalidateQueries({queryKey: MalAuthKeys.config()});
            }
        };

        window.addEventListener("message", handleMessage);

        return () => {
            window.removeEventListener("message", handleMessage);
        };
    }, [queryClient]);

    const deleteMutation = useMutation({
        mutationFn: APIClient.malauth.delete,
        onSuccess: () => {
            queryClient.invalidateQueries(MalQueryOptions());
            displayNotification({
                title: "MyAnimeList Authentication Deleted Successfully",
                message: "MyAnimeList updates will no longer be sent",
                type: "success",
            });
        },
        onError: (error) => {
            displayNotification({
                title: "Deletion Failed",
                message: error.message || "Could not delete MyAnimeList authentication",
                type: "error",
            });
        },
    })

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

    const startMutation = useMutation({
        mutationFn: () => APIClient.malauth.start(form.getValues().clientID, form.getValues().clientSecret),
        onSuccess: (data) => {
            const url = data?.url;
            if (url) {
                const width = 600;
                const height = 700;
                const left = window.screenX + (window.innerWidth - width) / 2;
                const top = window.screenY + (window.innerHeight - height) / 2;

                window.open(
                    url,
                    "MyAnimeList.net OAuth", // popup window name
                    `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`
                );
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

    const handleFormSubmit = () => {
        startMutation.mutate();
    };

    return (
        <Paper withBorder={true} p="md" shadow="xl">
            <Stack>
                <Title order={1} mt="md">
                    MyAnimeList
                </Title>
                <Text>
                    Manage the connection to your MyAnimeList account here.
                </Text>
                <Divider/>
                {(isLoading || (!isEmptySettings && testSucess === null)) ? (
                    <Loader/>
                ) : (
                    <>
                        {isEmptySettings ? (
                            <>
                                <Flex
                                    direction={"column"}
                                    gap={"md"}
                                    align={"stretch"}
                                    w={"350px"}
                                    mx={"auto"}
                                >
                                    <Group justify={"center"}>
                                        <Text fw={600} size={"md"}>
                                            Login with MyAnimeList.net
                                        </Text>
                                    </Group>
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
                                        <Group justify={"center"}>
                                            <Button type="submit" mt={"md"}>
                                                LOGIN
                                            </Button>
                                        </Group>
                                    </form>
                                </Flex>
                            </>
                        ) : (
                            <>
                                <Flex align={"center"} justify={"center"}>
                                    <Text size={"xl"} fw={600}>
                                        Authentication Status:
                                    </Text>
                                    <Text c={testSucess ? "green" : "red"} size={"md"} fw={600} ml={"xs"} mt={3}>
                                        {testSucess ? "OK" : "Failed"}
                                    </Text>
                                </Flex>
                                <Group justify="center">
                                    <ConfirmDeleteButton
                                        message={"MyAnimeList.net credentials will be deleted."}
                                        confirmText={"REMOVE ACCESS"}
                                        onConfirm={() => deleteMutation.mutate()}
                                    />
                                </Group>
                            </>
                        )}
                    </>
                )}
            </Stack>
        </Paper>
    );
};