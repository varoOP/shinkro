import React, {useEffect, useState} from "react";
import {useMutation, useQueryErrorResetBoundary} from "@tanstack/react-query";
import {useNavigate, useRouter, useSearch} from "@tanstack/react-router";
import {FaQuestion} from "react-icons/fa";
import {APIClient} from "@api/APIClient";
import {displayNotification} from "@components/notifications";
import {
    PasswordInput,
    TextInput,
    Button,
    Title,
    Stack,
    Paper,
    Tooltip,
    Image,
    Group,
    ActionIcon
} from "@mantine/core";
import {useForm} from "@mantine/form";
import {LoginRoute} from "@app/routes";
import Logo from "@app/logo.svg";
import {AuthContext} from "@utils/Context";
import classes from "./Auth.module.css";

type LoginFormFields = {
    username: string;
    password: string;
};

export const Login = () => {
    const [auth, setAuth] = AuthContext.use();

    const queryErrorResetBoundary = useQueryErrorResetBoundary();

    const router = useRouter();
    const search = useSearch({from: LoginRoute.id});
    const navigate = useNavigate();

    const form = useForm<LoginFormFields>({
        initialValues: {username: "", password: ""},
        validate: {
            username: (value) => (value ? null : "Username is required"),
            password: (value) => (value ? null : "Password is required"),
        },
    });

    useEffect(() => {
        queryErrorResetBoundary.reset();
        AuthContext.reset();
    }, [queryErrorResetBoundary]);

    const loginMutation = useMutation({
        mutationFn: (data: LoginFormFields) =>
            APIClient.auth.login(data.username, data.password),
        onSuccess: (_, variables: LoginFormFields) => {
            queryErrorResetBoundary.reset();
            setAuth({
                isLoggedIn: true,
                username: variables.username,
            });
            router.invalidate();
        },
        onError: (error) => {
            displayNotification({
                title: "Login Error",
                message: error.message || "An error occurred!",
                type: "error",
            });
        },
    });

    const onSubmit = (data: LoginFormFields) => loginMutation.mutate(data);

    const [tooltipOpened, setTooltipOpened] = useState(false);

    const commandText =
        "shinkro --config=${HOME}/.config/shinkro change-password ${USER} --password <new_password>";

    const handleCopyCommand = () => {
        navigator.clipboard.writeText(commandText);

        // Show the tooltip and notification
        setTooltipOpened(true);
        displayNotification({
            title: "Password Reset",
            message: "Command copied to clipboard",
            type: "info",
        });

        // Hide the tooltip after a delay
        setTimeout(() => setTooltipOpened(false), 5000);
    };

    React.useLayoutEffect(() => {
        if (auth.isLoggedIn && search.redirect) {
            void navigate({to: search.redirect});
        } else if (auth.isLoggedIn) {
            void navigate({to: "/"});
        }
    }, [auth.isLoggedIn, search.redirect]);

    return (
        <div className={classes.outerContainer}>
            <div className={classes.innerContainer}>
                <Image src={Logo} fit="contain" h={100}/>
                <Title order={2} ta="center">
                    shinkro
                </Title>
                <Paper withBorder={true} shadow={"xl"} mt={"md"} p={"xl"}>
                    <form onSubmit={form.onSubmit(onSubmit)}>
                        <Stack
                        >
                            <TextInput
                                {...form.getInputProps("username")}
                                label="Username"
                                placeholder="Enter username"
                            />
                            <PasswordInput
                                {...form.getInputProps("password")}
                                label="Password"
                                placeholder="Enter password"
                                mt="sm"
                            />
                            <Group justify={"center"} gap="xs">
                                <Button type="submit" mt="sm">
                                    LOGIN
                                </Button>
                                <Tooltip
                                    label={`Set new password using terminal command: ${commandText}`}
                                    opened={tooltipOpened}
                                    position="bottom"
                                    arrowOffset={50}
                                    arrowSize={5}
                                    withArrow
                                >
                                    <ActionIcon onClick={handleCopyCommand} mt="sm" radius={"xl"} size={"md"}
                                                variant={"outline"}>
                                        <FaQuestion
                                            size={15}
                                        />
                                    </ActionIcon>
                                </Tooltip>
                            </Group>
                        </Stack>
                    </form>
                </Paper>
            </div>
        </div>
    );
};
