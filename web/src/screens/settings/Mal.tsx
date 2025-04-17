import {APIClient} from "@api/APIClient.ts";
import {displayNotification} from "@components/notifications";
import {useState, useEffect} from "react";
import {
    Stack,
    Button,
    Loader,
    Group,
} from "@mantine/core";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {MalQueryOptions} from "@api/queries.ts";
import {MalAuthKeys} from "@api/query_keys.ts";
import {ConfirmDeleteButton} from "@components/alerts/ConfirmDeleteButton";
import {useDisclosure} from "@mantine/hooks";
import {MalForm} from "@forms/settings/MalForm.tsx";
import {CenteredEmptyState, SettingsSectionHeader, StatusIndicator} from "@screens/settings/components.tsx";

export const Mal = () => {
    const queryClient = useQueryClient();
    const {data: malauth, isLoading} = useQuery(MalQueryOptions());
    const [loading, setLoading] = useState(false);

    const [opened, {open, close}] = useDisclosure(false);
    const [testSucess, setTestSuccess] = useState<boolean | null>(null);

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
            if (event.data?.type === "mal-auth") {
                queryClient.invalidateQueries({queryKey: MalAuthKeys.config()});
                setLoading(false);
                close();
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

    return (
        <>
            <Stack>
                <SettingsSectionHeader
                    title={"MyAnimeList"}
                    description={"Manage the connection to your MyAnimeList account here."}
                />
                {(isLoading || (!isEmptySettings && testSucess === null)) ? (
                    <Loader mt={"md"} mx={"auto"}/>
                ) : (
                    <>
                        {isEmptySettings ? (
                            <CenteredEmptyState
                                message={"No MyAnimeList Credentials Found"}
                                button={
                                    <Button onClick={open}>
                                        START AUTHENTICATION
                                    </Button>
                                }
                            />
                        ) : (
                            <>
                                <StatusIndicator label={"Authentication Status:"} status={testSucess}/>
                                <Group justify="flex-start">
                                    <ConfirmDeleteButton
                                        message={"MyAnimeList.net credentials will be deleted."}
                                        confirmText={"REMOVE ACCESS"}
                                        onConfirm={() => deleteMutation.mutate()}
                                    />
                                    <Button onClick={open}>RE - AUTHENTICATE</Button>
                                </Group>
                            </>
                        )}
                    </>
                )}
            </Stack>
            <MalForm
                opened={opened}
                onClose={close}
                loading={loading}
                setLoading={setLoading}
            />
        </>
    );
};