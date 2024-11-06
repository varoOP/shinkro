import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { useForm } from "@mantine/form";
import {
  Button,
  Center,
  Image,
  Stack,
  TextInput,
  PasswordInput,
  Group,
  Title,
} from "@mantine/core";

import { APIClient } from "@api/APIClient";

import Logo from "@app/logo.svg";

interface InputValues {
  username: string;
  password: {
    pass: string;
    confirmPass: string;
  };
}

export const Onboarding = () => {
  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      username: "",
      password: {
        pass: "",
        confirmPass: "",
      },
    },
    validate: {
      username: (value) => (value.length < 1 ? "Input a username" : null),
      password: {
        pass: (value) => (value.length < 1 ? "Input a password" : null),
        confirmPass: (value, values) =>
          value !== values.password.pass ? "Passwords do not match" : null,
      },
    },
  });
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: (data: InputValues) =>
      APIClient.auth.onboard(data.username, data.password.pass),
    onSuccess: () => navigate({ to: "/login" }),
  });

  const handleSubmit = (values: InputValues) => {
    mutation.mutate(values);
  };

  return (
    <Center style={{ height: "100vh" }}>
      <Stack align="center">
        <Image src={Logo} alt="Logo" width={80} height={80} />
        <Title order={2}>shinkro</Title>
        <form onSubmit={form.onSubmit(handleSubmit)} style={{ width: "100%" }}>
          <TextInput
            label="Username"
            placeholder="Enter username"
            {...form.getInputProps("username")}
          />
          <PasswordInput
            label="Password"
            placeholder="Enter password"
            {...form.getInputProps("pass")}
          />
          <PasswordInput
            label="Confirm Password"
            placeholder="Confirm password"
            {...form.getInputProps("confirmPass")}
          />
          <Group justify="center" mt="md">
            <Button type="submit">Create Account</Button>
          </Group>
        </form>
      </Stack>
    </Center>
  );
};
