import {Title, Text, Divider, Group, Flex, Stack} from "@mantine/core";
import {ReactNode} from "react";

export const SettingsSectionHeader = ({
                                          title,
                                          description,
                                          right,
                                      }: {
    title: string;
    description: string;
    right?: ReactNode;
}) => (
    <div>
        <Group align={"flex-end"} gap={"xl"}>
            <Stack>
                <Title order={1} mt="md">{title}</Title>
                <Text>{description}</Text>
            </Stack>
            {right && right}
        </Group>
        <Divider mt={"md"}/>
    </div>
);

export const StatusIndicator = ({
                                    label,
                                    status,
                                }: {
    label: string;
    status: boolean | null;
}) => (
    <Flex align="center" justify="flex-start">
        <Text size="xl" fw={600}>{label}</Text>
        <Text
            c={status === null ? "gray" : status ? "green" : "red"}
            size="md"
            fw={600}
            ml="xs"
            mt={3}
        >
            {status === null ? "Unknown" : status ? "OK" : "Failed"}
        </Text>
    </Flex>
);

export const CenteredEmptyState = ({
                                       message,
                                       button,
                                   }: {
    message: string;
    button?: ReactNode;
}) => (
    <>
        <Group justify="center">
            <Text size="md" fw={600}>{message}</Text>
        </Group>
        <Group justify="center">{button}</Group>
    </>
);

export const BoundedFormContainer = ({children}: { children: ReactNode }) => (
    <Flex
        direction="column"
        gap="md"
        align="stretch"
        w="100%"
        maw="600px"
        miw="280px"
        mx="auto"
    >
        {children}
    </Flex>
);