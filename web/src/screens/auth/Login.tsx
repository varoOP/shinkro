import React, { useEffect, useState } from "react";
import { useMutation, useQueryErrorResetBoundary } from "@tanstack/react-query";
import { useRouter, useSearch } from "@tanstack/react-router";

import { APIClient } from "@api/APIClient";
import { showNotification } from "@mantine/notifications";
import {
  PasswordInput,
  TextInput,
  Button,
  Text,
  Stack,
  Paper,
  Tooltip,
  Container,
  Image,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { LoginRoute } from "@app/routes";

import Logo from "@app/logo.svg";
import { AuthContext } from "@utils/Context";

type LoginFormFields = {
  username: string;
  password: string;
};

export const Login = () => {
  const [auth, setAuth] = AuthContext.use();

  const queryErrorResetBoundary = useQueryErrorResetBoundary();

  const router = useRouter();
  const search = useSearch({ from: LoginRoute.id });

  const form = useForm<LoginFormFields>({
    initialValues: { username: "", password: "" },
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
    onError: (error: any) => {
      showNotification({
        title: "Error",
        message: error.message || "An error occurred!",
        color: "red",
      });
    },
  });

  const onSubmit = (data: LoginFormFields) => loginMutation.mutate(data);

  const [tooltipOpened, setTooltipOpened] = useState(false);

  const handleCopyCommand = () => {
    const commandText =
      "shinkro --config=/home/username/.config/shinkro change-password $USERNAME";
    navigator.clipboard.writeText(commandText);

    // Show the tooltip and notification
    setTooltipOpened(true);
    showNotification({
      title: "Password Reset",
      message: "Command copied to clipboard",
    });

    // Hide the tooltip after a delay
    setTimeout(() => setTooltipOpened(false), 5000);
  };

  React.useLayoutEffect(() => {
    if (auth.isLoggedIn && search.redirect) {
      router.history.push(search.redirect);
    } else if (auth.isLoggedIn) {
      router.history.push("/");
    }
  }, [auth.isLoggedIn, search.redirect]);

  return (
    <Container>
      <Image src={Logo} fit="contain" h={100} />
      <Text ta="center" size="xl" c="dark" fw={700}>
        shinkro
      </Text>
      <Paper shadow="md" radius="xl" withBorder p="xl">
        <Stack
          bg="var(--mantine-color-body)"
          align="strech"
          justify="center"
          gap="sm"
        >
          <form onSubmit={form.onSubmit(onSubmit)}>
            <TextInput label="USERNAME" {...form.getInputProps("username")} />
            <PasswordInput
              label="PASSWORD"
              {...form.getInputProps("password")}
            />
            <Button fullWidth type="submit" mt="sm">
              Login
            </Button>
            <Tooltip
              label="Reset password using terminal command: shinkro --config=/home/username/.config/shinkro change-password $USERNAME"
              opened={tooltipOpened}
              position="bottom"
              arrowOffset={50}
              arrowSize={5}
              withArrow
            >
              <Button fullWidth onClick={handleCopyCommand} mt="sm">
                Forgot Password
              </Button>
            </Tooltip>
          </form>
        </Stack>
      </Paper>
    </Container>
  );
};
