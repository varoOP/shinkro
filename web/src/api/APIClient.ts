import {baseUrl, sseBaseUrl} from "@utils";
import {GithubRelease} from "@app/types/Update";
import {
    PlexConfig,
    PlexOAuthPollResponse,
    PlexOAuthStartResponse,
    PlexLibraryResponse, PlexServerResponse
} from "@app/types/Plex";
import {AuthContext} from "@utils/Context";
import {MalAuthOpts} from "@app/types/MalAuth";

type RequestBody = BodyInit | object | Record<string, unknown> | null;
type Primitive = string | number | boolean | symbol | undefined;

interface HttpConfig {
    /**
     * One of "GET", "POST", "PUT", "PATCH", "DELETE", etc.
     * See https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
     */
    method?: string;
    /**
     * JSON body for this request. Once this is set to an object,
     * then `Content-Type` for this request is set to `application/json`
     * automatically.
     */
    body?: RequestBody;
    /**
     * Helper to work with a query string/search param of a URL.
     * E.g. ?a=1&b=2&c=3
     *
     * Using this interface will automatically convert
     * the object values into RFC-3986-compliant strings.
     *
     * Keys will *NOT* be sanitized, and any whitespace and
     * invalid characters will remain.
     *
     * The only supported value types are:
     * numbers, booleans, strings and flat 1-D arrays.
     *
     * Objects as values are not supported.
     *
     * The supported values are serialized as follows:
     *  - undefined values are ignored
     *  - empty strings are ignored
     *  - empty strings inside arrays are ignored
     *  - empty arrays are ignored
     *  - arrays append each time with the key and for each child
     *    e.g. `{ arr: [1, 2, 3] }` will yield `?arr=1&arr=2&arr=3`
     *  - array items with an undefined value (or which serialize to an empty string) are ignored,
     *    e.g. `{ arr: [1, undefined, undefined] }` will yield `?arr=1`
     *    (NaN, +Inf, -Inf, etc. will remain since they are valid serializations)
     */
    queryString?: Record<string, Primitive | Primitive[]>;
}

/**
 * Encodes a string into a RFC-3986-compliant string.
 *
 * By default, encodeURIComponent will not encode
 * any of the following characters: !'()*
 *
 * So a simple regex replace is done which will replace
 * these characters with their hex-value representation.
 *
 * @param str Input string (dictionary value).
 * @returns A RFC-3986-compliant string variation of the input string.
 * @note See https://stackoverflow.com/a/62969380
 */
function encodeRFC3986URIComponent(str: string): string {
    return encodeURIComponent(str).replace(
        /[!'()*]/g,
        (c) => `%${c.charCodeAt(0).toString(16).toUpperCase()}`
    );
}

/**
 * Makes a request on the network and returns a promise.
 *
 * This function serves as both a request builder and a response interceptor.
 *
 * @param endpoint The endpoint path relative to the backend instance.
 * @param config A dictionary which specifies what information this network
 * request must relay during transport. See @ref HttpClient.
 * @returns A promise for the *sent* network request which must *  be await'ed or .then()-chained before it can be used.
 *
 * If the status code returned by the server is in the [200, 300) range, then this is considered a success.
 *    - This function resolves with an empty dictionary object, i.e. {}, if the status code is 204 No data
 *    - The parsed JSON body is returned by this method if the server returns `Content-Type: application/json`.
 *    - In all other scenarios, the raw Response object from window.fetch() is returned,
 *      which must be handled manually by awaiting on one of its methods.
 *
 * The following is done if the status code that the server returns is NOT successful,
 * that is, if it falls outside of the [200, 300] range:
 *  - A unique Error object is returned if the user is logged in and the status code is 403 Forbidden.
 *    This Error object *should* be consumed by the @tanstack/query code, which indirectly calls HttpClient.
 *    The current user is then prompted to log in again after being logged out.
 *  - The `ErrorPage` screen appears in all other scenarios.
 */
export async function HttpClient<T = unknown>(
    endpoint: string,
    config: HttpConfig = {}
): Promise<T> {
    const init: RequestInit = {
        method: config.method,
        headers: {Accept: "*/*", "x-requested-with": "XMLHttpRequest"},
        credentials: "include",
    };

    if (config.body) {
        init.body = JSON.stringify(config.body);

        if (typeof config.body === "object") {
            init.headers = {
                ...init.headers,
                "Content-Type": "application/json",
            };
        }
    }

    if (config.queryString) {
        const params: string[] = [];

        for (const [key, value] of Object.entries(config.queryString)) {
            const serializedKey = encodeRFC3986URIComponent(key);

            if (typeof value === "undefined") {
                // Skip case when the value is undefined.
                // The solution in this case is to use the request body instead with JSON
                continue;
            } else if (Array.isArray(value)) {
                // Append (don't set) each array member as a query parameter
                // e.g. ?a=1&a=2&a=3
                value.forEach((child) => {
                    // Skip undefined member values
                    const v = typeof child !== "undefined" ? String(child) : "";
                    if (v.length) {
                        params.push(`${serializedKey}=${encodeRFC3986URIComponent(v)}`);
                    }
                });
            } else {
                // This is a primitive value, just add as string
                // e.g. ?a=1
                const v = String(value);
                if (v.length) {
                    params.push(`${serializedKey}=${encodeRFC3986URIComponent(v)}`);
                }
            }
        }

        if (params.length) {
            endpoint += `?${params.join("&")}`;
        }
    }

    const response = await window.fetch(`${baseUrl()}${endpoint}`, init);
    const isJson = response.headers
        .get("Content-Type")
        ?.includes("application/json");

    if (response.status >= 200 && response.status < 300) {
        // We received a successful response
        if (response.status === 204) {
            // 204 contains no data, but indicates success
            return Promise.resolve<T>({} as T);
        }

        // If Content-Type is application/json, then parse response as JSON
        // otherwise, just resolve the Response object returned by window.fetch
        // and the consumer can call await response.text() if needed.
        if (isJson) {
            return Promise.resolve<T>((await response.json()) as T);
        } else {
            return Promise.resolve<T>(response as T);
        }
    } else {
        // This is not a successful response.
        // It is most likely an error.
        switch (response.status) {
            case 403: {
                if (AuthContext.get().isLoggedIn) {
                    return Promise.reject(new Error("Cookie expired or invalid."));
                }
                break;
            }
            case 500: {
                const health = await window.fetch(`${baseUrl()}api/healthz/liveness`);
                if (!health.ok) {
                    return Promise.reject(
                        new Error(`[500] Offline (Internal server error): "${endpoint}"`)
                    );
                }
                break;
            }
            case 503: {
                // Show an error toast to notify the user what occurred
                return Promise.reject(
                    new Error(`[503] Service unavailable: "${endpoint}"`)
                );
            }
            default:
                break;
        }

        let reason = response.statusText;
        if (isJson) {
            const json = await response.json();
            if (Object.hasOwn(json, "message")) {
                reason = json.message as string;
            }
        }

        if (reason.length) {
            reason = ` (${reason})`;
        }

        const defaultError = new Error(
            `HTTP request to '${endpoint}' failed with code ${response.status}${reason}`
        );
        return Promise.reject(defaultError);
    }
}

const appClient = {
    Get: <T>(endpoint: string, config: HttpConfig = {}) =>
        HttpClient<T>(endpoint, {
            ...config,
            method: "GET",
        }),
    Post: <T = void>(endpoint: string, config: HttpConfig = {}) =>
        HttpClient<T>(endpoint, {
            ...config,
            method: "POST",
        }),
    Put: <T = void>(endpoint: string, config: HttpConfig = {}) =>
        HttpClient<T>(endpoint, {
            ...config,
            method: "PUT",
        }),
    Patch: (endpoint: string, config: HttpConfig = {}) =>
        HttpClient<void>(endpoint, {
            ...config,
            method: "PATCH",
        }),
    Delete: (endpoint: string, config: HttpConfig = {}) =>
        HttpClient<void>(endpoint, {
            ...config,
            method: "DELETE",
        }),
};

export const APIClient = {
    auth: {
        login: (username: string, password: string) =>
            appClient.Post("api/auth/login", {
                body: {username, password},
            }),
        logout: () => appClient.Post("api/auth/logout"),
        validate: () => appClient.Get<void>("api/auth/validate"),
        onboard: (username: string, password: string) =>
            appClient.Post("api/auth/onboard", {
                body: {username, password},
            }),
        canOnboard: () => appClient.Get("api/auth/onboard"),
        updateUser: (req: UserUpdate) =>
            appClient.Patch(`api/auth/user/${req.username_current}`, {body: req}),
    },
    apikeys: {
        getAll: () => appClient.Get<APIKey[]>("api/keys"),
        create: (key: APIKey) =>
            appClient.Post("api/keys", {
                body: key,
            }),
        delete: (key: string) => appClient.Delete(`api/keys/${key}`),
    },
    config: {
        get: () => appClient.Get<Config>("api/config"),
        update: (config: ConfigUpdate) =>
            appClient.Patch("api/config", {
                body: config,
            }),
    },
    plex: {
        getSettings: () =>
            appClient.Get<PlexConfig>("api/plex/settings"),

        updateSettings: (config: PlexConfig) =>
            appClient.Put("api/plex/settings", {
                body: config,
            }),

        delete: () =>
            appClient.Delete("api/plex/settings"),

        test: (config: PlexConfig) =>
            appClient.Post("api/plex/settings/test", {
                body: config,
            }),

        testToken: () =>
            appClient.Get("api/plex/settings/testToken"),

        startOAuth: () =>
            appClient.Post<PlexOAuthStartResponse>("api/plex/settings/oauth"),

        pollOAuth: (pin_id: number, client_id: string, code: string) =>
            appClient.Get<PlexOAuthPollResponse>("api/plex/settings/oauth", {
                queryString: {
                    pin_id,
                    client_id,
                    code,
                },
            }),

        servers: (config: PlexConfig) =>
            appClient.Post<PlexServerResponse>("api/plex/settings/servers", {
                body: config,
            }),

        libraries: (config: PlexConfig) =>
            appClient.Post<PlexLibraryResponse>("api/plex/settings/libraries", {
                body: config,
            }),
    },

    malauth: {
        callback: (data: MalAuthOpts) =>
            appClient.Post("api/malauth/callback", {
                body: data,
            }),


        getOpts: () =>
            appClient.Get<MalAuthOpts>("api/malauth/opts"),


        storeOpts: (data: MalAuthOpts) =>
            appClient.Post("api/malauth/opts", {
                body: data,
            }),

        test: () =>
            appClient.Get("api/malauth/test"),

    },

    logs: {
        files: () => appClient.Get<LogFileResponse>("api/logs/files"),
        getFile: (file: string) => appClient.Get(`api/logs/files/${file}`),
    },
    events: {
        logs: () =>
            new EventSource(`${sseBaseUrl()}api/events?stream=logs`, {
                withCredentials: true,
            }),
    },
    notifications: {
        getAll: () => appClient.Get<ServiceNotification[]>("api/notification"),
        create: (notification: ServiceNotification) =>
            appClient.Post("api/notification", {
                body: notification,
            }),
        update: (notification: ServiceNotification) =>
            appClient.Put(`api/notification/${notification.id}`, {
                body: notification,
            }),
        delete: (id: number) => appClient.Delete(`api/notification/${id}`),
        test: (notification: ServiceNotification) =>
            appClient.Post("api/notification/test", {
                body: notification,
            }),
    },
    updates: {
        check: () => appClient.Get("api/updates/check"),
        getLatestRelease: () => appClient.Get<GithubRelease>("api/updates/latest"),
    },
};
