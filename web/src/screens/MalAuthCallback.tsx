import {useMutation} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient.ts";
import {Flex, Group, Loader, Stack, Text, Button, Paper, Image} from "@mantine/core";
import {SiMyanimelist} from "react-icons/si";
import {FaArrowRightArrowLeft} from "react-icons/fa6";
import {useEffect} from "react";
import Logo from "@app/logo.svg";

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
        if (mutation.isSuccess || mutation.isError) {
            window.opener?.postMessage({type: "mal-auth"}, window.location.origin);
        }
    }, [mutation.isSuccess, mutation.isError]);

    return (
        <Flex
            direction={"column"}
            w={"100%"}
            maw={"600px"}
            miw={"280px"}
            mx={"auto"}
            pt={"10vh"}
            align={"stretch"}
        >
            <Paper
                withBorder
                p="md"
                shadow="xl"
            >
                <Group justify={"center"}>
                    <Image src={Logo} fit="contain" h={80}/>
                    <FaArrowRightArrowLeft size={50}/>
                    <SiMyanimelist size={100} color="#2e51a2"/>
                </Group>
                <Stack align="center">
                    {mutation.isPending ? (
                        <>
                            <Loader size="xl"/>
                            <Text>Authenticating with MyAnimeList...</Text>
                        </>
                    ) : (
                        <>
                            {mutation.isSuccess ? (
                                <Text size="xl" fw={600} c="green">Authentication Successful!</Text>
                            ) : (
                                <Text size="xl" fw={600} c="red">Authentication Failed!</Text>
                            )}
                        </>
                    )}
                    <Text c={"dimmed"}>You may close this window now.</Text>
                    <Button onClick={() => window.close()}>CLOSE WINDOW</Button>
                </Stack>
            </Paper>
        </Flex>
    );
};