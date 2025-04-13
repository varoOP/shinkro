import {useMutation} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient.ts";
import {Center, Loader, Stack, Text, Button} from "@mantine/core";
import {useEffect} from "react";

export const MalAuthCallback = () => {
    const mutation = useMutation({
        mutationFn: () => {
            const params = new URLSearchParams(window.location.search);
            const code = params.get("code");
            const state = params.get("state");

            if (!code || !state) {
                throw new Error("Missing code or state in URL");
            }

            return APIClient.malauth.callback(code, state);
        },
    });

    useEffect(() => {
        mutation.mutate();
    }, []);

    useEffect(() => {
        if (mutation.isSuccess) {
            window.opener?.postMessage({type: "mal-auth-success"}, window.location.origin);
            const timer = setTimeout(() => {
                window.close();
            }, 100);

            return () => clearTimeout(timer);
        }
    }, [mutation.isSuccess]);

    if (mutation.isPending) {
        return (
            <Center style={{height: "100vh"}}>
                <Stack align="center">
                    <Loader size="xl"/>
                    <Text>Authenticating with MyAnimeList...</Text>
                </Stack>
            </Center>
        );
    }

    if (mutation.isSuccess) {
        return (
            <Center style={{height: "100vh"}}>
                <Stack align="center">
                    <Text size="xl" fw={600} c="green">Authentication successful!</Text>
                    <Text>Closing Window in a few seconds..</Text>
                </Stack>
            </Center>
        );
    }

    if (mutation.isError) {
        return (
            <Center style={{height: "100vh"}}>
                <Stack align="center">
                    <Text size="xl" fw={600} c="red">Authentication failed!</Text>
                    <Text>Please close this window and try again.</Text>
                    <Button onClick={() => window.close()}>Close Window</Button>
                </Stack>
            </Center>
        );
    }

    return null;
};