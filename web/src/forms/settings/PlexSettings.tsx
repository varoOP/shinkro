import {
    Modal,
    TextInput,
    Switch,
    Button,
    Stack,
    Group,
    MultiSelect,
    Select,
    PasswordInput,
} from "@mantine/core";
import {TfiReload} from "react-icons/tfi";
import {useForm} from "@mantine/form";
import {useEffect, useState} from "react";
import {
    PlexConfig,
    PlexServer,
    PlexLibrary,
    PlexServerResponse,
} from "@app/types/Plex";
import {APIClient} from "@api/APIClient";
import {displayNotification} from "@components/notifications";

interface Props {
    opened: boolean;
    onClose: () => void;
    onSubmit: (values: PlexConfig) => void;
    defaultValues?: Partial<PlexConfig>;
}

export const PlexSettings = ({
                                 opened,
                                 onClose,
                                 onSubmit,
                                 defaultValues,
                             }: Props) => {
    const form = useForm<PlexConfig>({
        initialValues: {
            host: "",
            port: 0,
            tls: false,
            tls_skip: false,
            token: "",
            plex_user: "",
            anime_libs: [],
            plex_client_enabled: true,
            client_id: "",
            ...defaultValues,
        },
        validate: {
            host: (v: string) => (v ? null : "Host is required"),
            port: (v: number) => (v ? null : "Port is required"),
            token: (v: string) => (v ? null : "Token is required"),
            plex_user: (v: string) => (v ? null : "Username is required"),
            anime_libs: (v: string[]) =>
                v.length === 0 ? "At least one anime library is required" : null,
            client_id: (v: string) => (v ? null : "Client ID is required"),
        },
    });

    const [polling, setPolling] = useState(false);
    const [oauthInfo, setOauthInfo] = useState<{
        pin_id: number;
        client_id: string;
        code: string;
    } | null>(null);

    // Servers state and selected server (using Select)
    const [servers, setServers] = useState<PlexServerResponse | null>(null);
    const [authenticated, setAuthenticated] = useState(false);
    const [selectedServer, setSelectedServer] = useState<PlexServer | null>(null);
    const [loadingServers, setLoadingServers] = useState(false);

    // Libraries state for anime libraries
    const [libraries, setLibraries] = useState<PlexLibrary[]>([]);
    const [loadingLibraries, setLoadingLibraries] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);

    useEffect(() => {
        if (defaultValues) {
            form.setValues(defaultValues);
        }
    }, [defaultValues]);

    const handlePlexLogin = async () => {
        try {
            await APIClient.plex.testToken()
            setAuthenticated(true);
            displayNotification(
                {
                    title: "Plex Already Authenticated",
                    message: "Loading Servers..",
                    type: "info",
                }
            )
            return
        } catch (err) {
            setAuthenticated(false);
            const error = err as Error;
            console.error(error);
        }
        try {
            setPolling(true);
            const result = await APIClient.plex.startOAuth();
            setOauthInfo({
                pin_id: result.pin_id,
                client_id: result.client_id,
                code: result.code,
            });
            window.open(result.auth_url, "_blank");
        } catch (err) {
            const error = err as Error;
            setPolling(false);
            displayNotification({
                title: "OAuth Error",
                message: error.message || "Failed to start Plex login",
                type: "error",
            });
        }
    };

    useEffect(() => {
        if (!polling || !oauthInfo) return;
        let timeoutReached = false;
        const timeout = setTimeout(() => {
            timeoutReached = true;
            setPolling(false);
            displayNotification({
                title: "Timeout",
                message: "Authentication with Plex timed out after 60 seconds",
                type: "error",
            });
        }, 60000);

        const interval = setInterval(async () => {
            if (timeoutReached) return;
            try {
                const result = await APIClient.plex.pollOAuth(
                    oauthInfo.pin_id,
                    oauthInfo.client_id,
                    oauthInfo.code
                );
                if (result.message === "waiting for auth") {
                    return;
                }
                setAuthenticated(true)
                form.setValues((prev) => ({
                    ...prev,
                    token: result.token,
                    plex_user: result.plex_user,
                    client_id: result.client_id,
                }));
                displayNotification({
                    title: "Plex Login Successful",
                    message: `Logged in as ${result.plex_user}`,
                    type: "success",
                });
                clearInterval(interval);
                clearTimeout(timeout);
                setPolling(false);
                setOauthInfo(null);
            } catch (err) {
                setAuthenticated(false);
                const error = err as Error;
                clearInterval(interval);
                clearTimeout(timeout);
                setPolling(false);
                displayNotification({
                    title: "Plex Login Failed",
                    message: error.message || "Polling error",
                    type: "error",
                });
            }
        }, 1000);

        return () => {
            clearInterval(interval);
            clearTimeout(timeout);
        };
    }, [polling, oauthInfo, form]);

    useEffect(() => {
        if (authenticated) {
            void loadServers();
        }
    }, [authenticated]);

    useEffect(() => {
        if (servers && servers.Servers.length > 0) {
            setSelectedServer(servers.Servers[0]);
        }
    }, [servers]);

    useEffect(() => {
        if (selectedServer) {
            form.setFieldValue("host", selectedServer.connections[0].address);
            form.setFieldValue("port", selectedServer.connections[0].port);
            form.setFieldValue("tls", selectedServer.connections[0].protocol === "https");

            void loadLibrariesForServer();
        }
    }, [selectedServer]);

    const loadServers = async () => {
        if (!form.getValues().token || !form.getValues().client_id) {
            displayNotification({
                title: "Missing Credentials",
                message: "Please authenticate with Plex first.",
                type: "error",
            });
            return;
        }
        try {
            setLoadingServers(true);
            const response = await APIClient.plex.servers(form.getValues());
            const serverList = (response?.Servers ?? []).filter(
                (server: PlexServer) => server.owned && server.provides === "server"
            );
            if (!serverList.length) {
                displayNotification({
                    title: "No Servers Found",
                    message: "No Plex servers found for your account.",
                    type: "error",
                });
                setServers(null);
                return;
            }
            setServers({Servers: serverList});
            displayNotification({
                title: "Servers Loaded",
                message: "Plex servers loaded successfully.",
                type: "success",
            });
        } catch (err) {
            const error = err as Error;
            displayNotification({
                title: "Failed to Load Servers",
                message: error.message || "Could not load Plex servers.",
                type: "error",
            });
            setServers(null);
        } finally {
            setLoadingServers(false);
        }
    };

    const handleServerSelection = (value: string | null) => {
        if (!value) {
            setSelectedServer(null);
            return;
        }
        const server = servers?.Servers.find((s) => s.clientIdentifier === value);
        setSelectedServer(server || null);
    };


    const loadLibrariesForServer = async () => {
        if (!selectedServer?.connections || selectedServer?.connections.length === 0) {
            displayNotification({
                title: "Invalid Server Data",
                message: "Selected server does not have any connections.",
                type: "error",
            });
            return;
        }

        try {
            setLoadingLibraries(true);
            const response = await APIClient.plex.libraries(form.getValues());
            console.log("Libraries response:", response);
            setLibraries(response.MediaContainer.Directory || []);
            displayNotification({
                title: "Libraries Loaded",
                message: "Anime libraries loaded successfully.",
                type: "success",
            });
        } catch (err) {
            const error = err as Error;
            displayNotification({
                title: "Failed to Load Libraries, Choose different Host and Port",
                message: error.message || "Could not load libraries.",
                type: "error",
            });
        } finally {
            setLoadingLibraries(false);
        }
    };


    // Prepare options for the MultiSelect (deduplicated)
    const libraryOptions = Array.from(
        new Set(libraries.map((lib) => lib.title))
    ).map((title) => ({value: title, label: title}));

    const testPlex = async () => {
        try {
            setTestingConnection(true)
            await APIClient.plex.test(form.getValues());
            displayNotification({
                title: "Test Successful",
                message: "Plex connection test was successful",
                type: "success",
            });
        } catch (err) {
            const error = err as Error;
            displayNotification({
                title: "Test Failed",
                message: error.message || "Plex connection test failed.",
                type: "error",
            });
        } finally {
            setTestingConnection(false);
        }
    }

    const handleFormSubmit = (values: PlexConfig) => {
        // Prepare settings for submission.
        const settings: PlexConfig = {
            host: values.host,
            port: values.port,
            token: "", // token is managed via OAuth
            client_id: values.client_id,
            tls: values.tls,
            tls_skip: values.tls_skip,
            anime_libs: values.anime_libs,
            plex_user: values.plex_user,
            plex_client_enabled: values.plex_client_enabled,
        };
        onSubmit(settings);
    };

    const {error} = form.getInputProps("port", {type: "input"});
    return (
        <Modal opened={opened} onClose={onClose} title="Plex Settings">
            <form onSubmit={form.onSubmit(handleFormSubmit)}>
                <Stack align="stretch">
                    <Button loading={polling} onClick={handlePlexLogin} variant="light">
                        Authenticate with Plex
                    </Button>
                    <Group align="flex-end" justify="flex-start">
                        <Select style={{flex: 1}}
                                label="Select Plex Server"
                                placeholder="Load servers after authenticating"
                                data={servers?.Servers.map((server) => ({
                                    value: server.clientIdentifier,
                                    label: server.name,
                                }))}
                                value={selectedServer ? selectedServer.clientIdentifier : ""}
                                onChange={handleServerSelection}
                                disabled={!servers || polling || loadingServers}
                                searchable
                                clearable
                        />
                        <Button onClick={loadServers} disabled={loadingServers || polling} loading={loadingServers}>
                            <TfiReload/>
                        </Button>
                    </Group>
                    <Select
                        label="Host"
                        placeholder="Select host after selecting server"
                        data={
                            selectedServer
                                ? selectedServer.connections.map((conn) => conn.address)
                                : [form.getValues().host]
                        }
                        value={form.getValues().host}
                        disabled={!selectedServer || loadingServers}
                        {...form.getInputProps("host", {type: "input"})}
                    />
                    <Select
                        label="Port"
                        placeholder="Select Port after selecting Host"
                        data={
                            selectedServer
                                ? selectedServer.connections.map((conn) => String(conn.port))
                                : [form.getValues().port === 0 ? "" : String(form.getValues().port)]
                        }
                        value={String(form.getValues().port)}
                        onChange={(value: string | null) => {
                            form.setFieldValue("port", value ? parseInt(value) : 0);
                        }}
                        disabled={!selectedServer || loadingServers}
                        error={error}
                    />
                    <Group align="flex-end" justify="flex-start">
                        <Switch
                            label="TLS (HTTPS)"
                            {...form.getInputProps("tls", {type: "checkbox"})}
                        />
                        <Switch
                            label="Skip TLS Verification"
                            {...form.getInputProps("tls_skip", {type: "checkbox"})}
                        />
                    </Group>
                    <Group align="flex-end">
                        <MultiSelect style={{flex: 1}}
                                     label="Select Anime Libraries"
                                     placeholder="Pick Anime Libraries"
                                     data={libraryOptions}
                                     value={form.getValues().anime_libs}
                                     disabled={loadingLibraries || libraryOptions.length === 0}
                                     searchable
                                     clearable
                                     hidePickedOptions
                                     {...form.getInputProps("anime_libs", {type: "input"})}
                        />
                        <Button onClick={loadLibrariesForServer} loading={loadingLibraries}
                                disabled={loadingLibraries || polling || loadingServers || !selectedServer}>
                            <TfiReload/>
                        </Button>
                    </Group>
                    <TextInput
                        label="Username"
                        placeholder="Enter Plex username"
                        {...form.getInputProps("plex_user")}
                    />
                    <TextInput
                        label="Client ID"
                        placeholder="To generate client ID, authenticate with Plex"
                        disabled
                        {...form.getInputProps("client_id")}
                    />
                    <PasswordInput
                        label="Token"
                        placeholder="To get token, authenticate with Plex"
                        disabled
                        {...form.getInputProps("token")}
                    />
                    <Switch
                        label="Enable Plex Client"
                        {...form.getInputProps("plex_client_enabled", {type: "checkbox"})}
                    />

                    <Group mt="sm" justify="flex-end">
                        <Button variant="default" onClick={onClose} disabled={polling}>
                            Cancel
                        </Button>
                        <Button onClick={testPlex} loading={testingConnection}>
                            Test
                        </Button>
                        <Button type="submit" disabled={polling}>
                            Save
                        </Button>
                    </Group>
                </Stack>
            </form>
        </Modal>
    );
};