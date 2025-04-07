import {queryOptions} from "@tanstack/react-query";
import {APIClient} from "@api/APIClient";
import {ApiKeys, NotificationKeys, SettingsKeys} from "@api/query_keys";

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
        queryKey: SettingsKeys.plex(),
        queryFn: () => APIClient.plex.getSettings(),
        enabled: enabled,
        retry: false,
        refetchOnWindowFocus: false,
    });

export const UpdatesQueryOptions = (enabled: boolean) =>
    queryOptions({
        queryKey: SettingsKeys.updates(),
        queryFn: () => APIClient.updates.getLatestRelease(),
        retry: false,
        refetchOnWindowFocus: false,
        enabled: enabled,
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
