import {queryOptions} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient";
import {
    ApiKeys,
    NotificationKeys,
    SettingsKeys,
    PlexSettingsKeys,
    MalAuthKeys,
    MappingKeys,
    LogKeys,
    PlexKeys,
} from "@api/query_keys";

export const ConfigQueryOptions = (enabled: boolean = true) =>
    queryOptions({
        queryKey: SettingsKeys.config(),
        queryFn: () => APIClient.config.get(),
        retry: false,
        refetchOnWindowFocus: false,
        enabled: enabled,
    });

export const PlexSettingsQueryOptions = (enabled: boolean = true) =>
    queryOptions({
        queryKey: PlexSettingsKeys.config(),
        queryFn: () => APIClient.plex.getSettings(),
        enabled: enabled,
        retry: false,
        refetchOnWindowFocus: false,
    });

export const MalQueryOptions = (enabled: boolean = true) =>
    queryOptions({
        queryKey: MalAuthKeys.config(),
        queryFn: () => APIClient.malauth.get(),
        enabled: enabled,
        retry: false,
        refetchOnWindowFocus: false,
    });

export const MappingQueryOptions = (enabled: boolean = true) =>
    queryOptions({
        queryKey: MappingKeys.lists(),
        queryFn: () => APIClient.mapping.get(),
        enabled: enabled,
        retry: false,
        refetchOnWindowFocus: false,
    });

export const LogQueryOptions = (enabled: boolean = true) =>
    queryOptions({
        queryKey: LogKeys.lists(),
        queryFn: () => APIClient.fs.listLogs(),
        enabled: enabled,
        retry: false,
        refetchOnWindowFocus: true,
    });

export const LogContentQueryOptions = (enabled: boolean = true) =>
    queryOptions({
        queryKey: LogKeys.content(),
        queryFn: async () => {
            const response = await fetch(`${window.location.origin}/api/fs/logs/shinkro.log`);
            if (!response.ok) throw new Error("Failed to fetch log");
            return response.text();
        },
        enabled: enabled,
        retry: false,
        refetchOnWindowFocus: false,
    });

export const NotificationsQueryOptions = () =>
    queryOptions({
        queryKey: NotificationKeys.lists(),
        queryFn: () => APIClient.notifications.getAll(),
    });

export const ApikeysQueryOptions = () =>
    queryOptions({
        queryKey: ApiKeys.lists(),
        queryFn: () => APIClient.apikeys.getAll(),
        refetchOnWindowFocus: false,
    });

export const plexHistoryQueryOptions = (opts: {
    type?: "timeline" | "table";
    limit?: number;
    cursor?: string;
    offset?: number;
    search?: string;
    status?: string;
    event?: string;
    from?: string;
    to?: string;
} = { type: "timeline", limit: 10 }) =>
    queryOptions({
        queryKey: PlexKeys.history(opts.type ?? "timeline", opts),
        queryFn: () => APIClient.plex.history(opts),
    });

export const plexCountsQueryOptions = () =>
    queryOptions({
        queryKey: ["plexCounts"],
        queryFn: () => APIClient.plex.getCounts(),
    });

export const animeUpdateCountQueryOptions = () =>
    queryOptions({
        queryKey: ["animeUpdateCount"],
        queryFn: () => APIClient.animeupdate.getCount(),
    });

export const recentAnimeUpdatesQueryOptions = (limit: number = 5) =>
    queryOptions({
        queryKey: ["recentAnimeUpdates", limit],
        queryFn: () => APIClient.animeupdate.getRecent(limit),
    });

export const latestReleaseQueryOptions = () =>
    queryOptions({
        queryKey: ["latestRelease"],
        queryFn: () => APIClient.updates.getLatestRelease(),
        staleTime: 1000 * 60 * 60, // 1 hour
    });
