import {ActionIcon, AppShell, Code, Flex, Group, Image, Menu, NavLink, rem, Title,} from "@mantine/core";
import Logo from "@app/logo.svg";
import {displayNotification} from "@components/notifications";
import {MdDarkMode, MdLightMode} from "react-icons/md";
import {FaDiscord, FaGithub, FaUser} from "react-icons/fa";
import {GrHelpBook} from "react-icons/gr";
import {BiLogOut} from "react-icons/bi";
import {useMutation, useQuery} from "@tanstack/react-query";
import {Link, Outlet, useRouter} from "@tanstack/react-router";
import {APIClient} from "@api/APIClient";
import {ConfigQueryOptions} from "@api/queries";
import {AuthContext, useThemeToggle} from "@utils/Context";
import {ExternalLink} from "@components/ExternalLink";
import {NAV_ROUTES} from "./navigation";
import classes from "./Layout.module.css";

export const Layout = () => {
    const router = useRouter();

    const {isError: isConfigError, error: configError, data: config} = useQuery(ConfigQueryOptions(true));
    if (isConfigError) {
        console.log(configError);
    }

    const {colorScheme, toggleTheme} = useThemeToggle();
    // const settings = SettingsContext.use();

    // Handle logout mutation
    const logoutMutation = useMutation({
        mutationFn: APIClient.auth.logout,
        onSuccess: () => {
            displayNotification({
                title: "Logged out",
                message: "You have been logged out. Goodbye!",
                type: "success",
            });
            AuthContext.reset();
            router.history.push("/login");
        },
        onError: (err) => {
            console.error("logout error", err);
        },
    });

    return (
        <AppShell
            padding="md"
            header={{height: 60}}
            navbar={{
                width: 250,
                breakpoint: "sm",
            }}
            className={classes.appshell}
        >
            <AppShell.Header className={classes.header}>
                <Group h="100%" px="md" align="center">
                    <Image src={Logo} height={60}/>
                    <Flex align="flex-end" gap="xs">
                        <Title order={3}>shinkro</Title>
                        <Code fw={700} className={classes.code}>
                            {config?.version}
                        </Code>
                    </Flex>
                    <Menu
                        shadow="md"
                        width={200}
                        position="bottom-start"
                        offset={8}
                        withArrow
                        arrowPosition="center"
                        transitionProps={{transition: "skew-up", duration: 150}}
                    >
                        <Menu.Target>
                            <ActionIcon
                                variant="outline"
                                radius="md"
                                size="lg"
                                style={{marginLeft: "auto"}}
                            >
                                <FaUser style={{width: "80%", height: "80%"}}/>
                            </ActionIcon>
                        </Menu.Target>

                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => logoutMutation.mutate()}
                                leftSection={
                                    <BiLogOut style={{width: rem(20), height: rem(20)}}/>
                                }
                            >
                                Logout
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>

                    <ExternalLink href="https://discord.gg/ZkYdfNgbAT">
                        <ActionIcon variant="filled" color="#7289da" radius="md" size="lg">
                            <FaDiscord style={{width: "80%", height: "80%"}}/>
                        </ActionIcon>
                    </ExternalLink>

                    <ExternalLink href="https://github.com/varoOP/shinkro">
                        <ActionIcon variant="default" radius="md" size="lg">
                            <FaGithub style={{width: "80%", height: "80%"}}/>
                        </ActionIcon>
                    </ExternalLink>

                    <ExternalLink href="https://docs.shinkro.com">
                        <ActionIcon variant="default" radius="md" size="lg">
                            <GrHelpBook style={{width: "80%", height: "80%"}}/>
                        </ActionIcon>
                    </ExternalLink>

                    <ActionIcon
                        variant="outline"
                        radius="xl"
                        size="lg"
                        onClick={toggleTheme}
                        aria-label="Theme Switch"
                    >
                        {colorScheme === "dark" ? (
                            <MdLightMode style={{width: "80%", height: "80%"}}/>
                        ) : (
                            <MdDarkMode style={{width: "80%", height: "80%"}}/>
                        )}
                    </ActionIcon>
                </Group>
            </AppShell.Header>

            <AppShell.Navbar className={classes.navbar}>
                {NAV_ROUTES.map((item, itemIdx) => (
                    <Link
                        key={item.name + itemIdx}
                        to={item.path}
                        params={{}}
                        style={{textDecoration: "none", color: "inherit"}}
                    >
                        {({isActive}) => {
                            return (
                                <>
                                    <NavLink
                                        component="button"
                                        label={item.name}
                                        active={isActive}
                                        variant="light"
                                        color="blue"
                                    />
                                </>);
                        }}
                    </Link>
                ))}
            </AppShell.Navbar>
            <AppShell.Main>
                <Flex
                    h="80vh"
                    w="70vw"
                    justify="center"
                    align="center"
                    direction="column"
                >
                    <Outlet/>
                </Flex>
            </AppShell.Main>
        </AppShell>
    );
};
