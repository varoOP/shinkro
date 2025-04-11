export interface PlexConfig {
    host: string;
    port: number;
    tls: boolean;
    tls_skip: boolean;
    token: string;
    client_id: string;
    anime_libs: string[];
    plex_user: string;
    plex_client_enabled: boolean;
}

export interface PlexOAuthStartResponse {
    pin_id: number;
    code: string;
    client_id: string;
    auth_url: string;
}

export interface PlexOAuthPollResponse {
    token: string;
    plex_user: string;
    client_id: string;
    message: string;
}

export interface PlexConnection {
    protocol: string;
    address: string;
    port: number;
    uri: string;
    local: boolean;
    relay: boolean;
    IPv6: boolean;
}

export interface PlexServerResponse {
    Servers: PlexServer[];
}

export interface PlexServer {
    name: string;
    product: string;
    productVersion: string;
    platform: string | null;
    platformVersion: string | null;
    device: string | null;
    clientIdentifier: string;
    createdAt: string; // or Date if you want to parse it
    lastSeenAt: string;
    provides: string;
    ownerId: number | null;
    sourceTitle: string | null;
    publicAddress: string;
    accessToken: string;
    owned: boolean;
    home: boolean;
    synced: boolean;
    relay: boolean;
    presence: boolean;
    httpsRequired: boolean;
    publicAddressMatches: boolean;
    dnsRebindingProtection: boolean;
    natLoopbackSupported: boolean;
    connections: PlexConnection[];
}

export interface PlexLibraryResponse {
    MediaContainer: {
        size: number;
        allowSync: boolean;
        title1: string;
        Directory: PlexLibrary[];
    };
}

export interface PlexLibrary {
    allowSync: boolean;
    art: string;
    composite: string;
    filters: boolean;
    refreshing: boolean;
    thumb: string;
    key: string;
    type: string;
    title: string;
    agent: string;
    scanner: string;
    language: string;
    uuid: string;
    updatedAt: number;
    createdAt: number;
    scannedAt: number;
    content: boolean;
    directory: boolean;
    contentChangedAt: number;
    hidden: number;
    Location: PlexLibraryLocation[];
}

export interface PlexLibraryLocation {
    id: number;
    path: string;
}