export const SettingsKeys = {
    all: ["settings"] as const,
    updates: () => [...SettingsKeys.all, "updates"] as const,
    config: () => [...SettingsKeys.all, "config"] as const,
    lists: () => [...SettingsKeys.all, "list"] as const,
    plex: () => [...SettingsKeys.all, "plex"] as const,
};

export const ApiKeys = {
    all: ["api_keys"] as const,
    lists: () => [...ApiKeys.all, "list"] as const,
    details: () => [...ApiKeys.all, "detail"] as const,
    detail: (id: string) => [...ApiKeys.details(), id] as const,
};

export const NotificationKeys = {
    all: ["notifications"] as const,
    lists: () => [...NotificationKeys.all, "list"] as const,
    details: () => [...NotificationKeys.all, "detail"] as const,
    detail: (id: number) => [...NotificationKeys.details(), id] as const,
};
