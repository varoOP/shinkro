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
    AnimeUpdateKeys
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

export const recentPlexPayloadsQueryOptions = (limit: number = 20) =>
    queryOptions({
        queryKey: PlexKeys.recentPayloads(limit),
        queryFn: () => APIClient.plex.getRecent(limit),
    });

export const animeUpdateByPlexIdQueryOptions = (plexId: number) =>
    queryOptions({
        queryKey: AnimeUpdateKeys.byPlexId(plexId),
        queryFn: () => APIClient.animeupdate.getByPlexId(plexId),
        enabled: !!plexId,
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
