import {Title, Text, Divider, Group, Flex, Stack, Loader, Tooltip, ActionIcon} from "@mantine/core";
import {ReactNode} from "react";
import {FaExternalLinkAlt} from "react-icons/fa";
import {ExternalLink} from "@components/ExternalLink.tsx";

export const SettingsSectionHeader = ({
                                          title,
                                          description,
                                          note,
                                          link,
                                      }: {
    title: string;
    description: string;
    note?: ReactNode;
    link?: string;
}) => (
    <div>
        <Group align={"flex-end"} gap={"xl"}>
            <Stack>
                <Group align={"flex-end"}>
                    <Title order={1} mt="md">{title}</Title>
                    {note && (
                        <ExternalLink href={link ?? ""}>
                            <Tooltip
                                label={note}
                                position="right"
                                withArrow
                                color={"blue.9"}
                                multiline={true}
                                maw={520}
                            >

                                <ActionIcon size="sm" color={"blue.9"}>
                                    <FaExternalLinkAlt/>
                                </ActionIcon>
                            </Tooltip>
                        </ExternalLink>
                    )}
                </Group>
                <Text>{description}</Text>
            </Stack>
        </Group>
        <Divider mt={"md"} size={"md"}/>
    </div>
);

export const StatusIndicator = ({
                                    label,
                                    status,
                                    loadStatus,
                                }: {
    label: string;
    status: boolean | null;
    loadStatus: boolean | null;
}) => (
    <Flex align="center" justify="flex-start">
        <Text size="xl" fw={600}>{label}</Text>
        {loadStatus ? (
            <Loader
                size={"sm"}
                ml={"xs"}
            />
        ) : (
            <Text
                c={status === null ? "gray" : status ? "green" : "red"}
                size="md"
                fw={600}
                ml="xs"
                mt={3}
            >
                {status === null ? "Unknown" : status ? "OK" : "Failed"}
            </Text>
        )}
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