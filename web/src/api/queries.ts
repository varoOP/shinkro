import {queryOptions} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient";
import {
    ApiKeys,
    NotificationKeys,
    SettingsKeys,
    PlexSettingsKeys,
    MalAuthKeys,
    MappingKeys,
    LogKeys
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

// export const UpdatesQueryOptions = (enabled: boolean) =>
//     queryOptions({
//         queryKey: SettingsKeys.updates(),
//         queryFn: () => APIClient.updates.getLatestRelease(),
//         retry: false,
//         refetchOnWindowFocus: false,
//         enabled: enabled,
//     });

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
