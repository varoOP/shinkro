import React, { useEffect, useState } from "react";
import { useMutation, useQueryErrorResetBoundary } from "@tanstack/react-query";
import { useRouter, useSearch } from "@tanstack/react-router";

import { APIClient } from "@api/APIClient";
import { displayNotification } from "@components/notifications";
import {
  PasswordInput,
  TextInput,
  Button,
  Title,
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
    "shinkro --config=${HOME}/.config/shinkro change-password ${USER}";

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
      router.history.push(search.redirect);
    } else if (auth.isLoggedIn) {
      router.history.push("/");
    }
  }, [auth.isLoggedIn, search.redirect]);

  return (
    <Container>
      <Image src={Logo} fit="contain" h={100} />
      <Title order={2} ta="center">
        shinkro
      </Title>
      <Paper>
        <Stack
          bg="var(--mantine-color-body)"
          align="strech"
          justify="center"
          gap="sm"
        >
          <form onSubmit={form.onSubmit(onSubmit)}>
            <TextInput
              {...form.getInputProps("username")}
              placeholder="USERNAME"
            />
            <PasswordInput
              {...form.getInputProps("password")}
              placeholder="PASSWORD"
              mt="sm"
            />
            <Button fullWidth type="submit" mt="sm">
              Login
            </Button>
            <Tooltip
              label={`Reset password using terminal command: ${commandText}`}
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
