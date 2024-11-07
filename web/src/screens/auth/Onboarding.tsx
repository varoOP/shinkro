import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { useForm } from "@mantine/form";
import {
  Button,
  Paper,
  Image,
  Stack,
  TextInput,
  PasswordInput,
  Group,
  Text,
  Container,
} from "@mantine/core";

import { APIClient } from "@api/APIClient";

import Logo from "@app/logo.svg";

interface InputValues {
  username: string;
  pass: string;
  confirmPass: string;
}

export const Onboarding = () => {
  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      username: "",
      pass: "",
      confirmPass: "",
    },
    validate: {
      username: (value) => (value.length < 1 ? "Input a username" : null),
      pass: (value) => (value.length < 1 ? "Input a password" : null),
      confirmPass: (value, values) =>
        value !== values.pass ? "Passwords do not match" : null,
    },
  });
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: (data: InputValues) =>
      APIClient.auth.onboard(data.username, data.pass),
    onSuccess: () => navigate({ to: "/login" }),
  });

  return (
    <Container>
      <Image src={Logo} fit="contain" h={100} alt="Logo" />
      <Text ta="center" size="xl" fw={700}>
        shinkro
      </Text>
      <Paper
        shadow="md"
        radius="xl"
        withBorder
        p="xl"
        style={{ width: 450, margin: "0 auto" }}
      >
        <Stack align="stretch" justify="center" gap="sm">
          <form
            onSubmit={form.onSubmit((values) => mutation.mutate(values))}
            style={{ width: "100%" }}
          >
            <TextInput
              label="USERNAME"
              placeholder="Your waifu's name works"
              {...form.getInputProps("username")}
            />
            <PasswordInput
              label="PASSWORD"
              placeholder="Cite the deep magic!"
              {...form.getInputProps("pass")}
            />
            <PasswordInput
              label="CONFIRM PASSWORD"
              placeholder="2X Deep Magics"
              {...form.getInputProps("confirmPass")}
            />
            <Group justify="center" mt="md">
              <Button type="submit" fullWidth>
                Create Account
              </Button>
            </Group>
          </form>
        </Stack>
      </Paper>
    </Container>
  );
};
