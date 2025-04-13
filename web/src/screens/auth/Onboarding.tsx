import {useMutation} from "@tanstack/react-query";
import {useNavigate} from "@tanstack/react-router";
import {useForm} from "@mantine/form";
import classes from "./Auth.module.css";
import {
    Button,
    Paper,
    Image,
    Stack,
    Group,
    TextInput,
    PasswordInput,
    Title,
} from "@mantine/core";

import {APIClient} from "@api/APIClient";

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
        onSuccess: () => navigate({to: "/login"}),
    });

    return (
        <div className={classes.outerContainer}>
            <div className={classes.innerContainer}>
                <Image src={Logo} fit="contain" h={100} alt="Logo"/>
                <Title ta="center" order={2}>
                    shinkro
                </Title>
                <Paper withBorder={true} shadow={"xl"} mt={"md"} p={"xl"}>
                    <form
                        onSubmit={form.onSubmit((values) => mutation.mutate(values))}
                    >
                        <Stack>
                            <TextInput mt="sm"
                                       placeholder="USERNAME"
                                       {...form.getInputProps("username")}
                            />
                            <PasswordInput mt="sm"
                                           placeholder="PASSWORD"
                                           {...form.getInputProps("pass")}
                            />
                            <PasswordInput mt="sm"
                                           placeholder="CONFIRM PASSWORD"
                                           {...form.getInputProps("confirmPass")}
                            />
                            <Group justify={"center"}>
                                <Button type="submit">
                                    CREATE ACCOUNT
                                </Button>
                            </Group>
                        </Stack>
                    </form>
                </Paper>
            </div>
        </div>
    );
};
