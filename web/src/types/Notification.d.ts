type NotificationType = "DISCORD" | "GOTIFY";
type NotificationEvent = "SUCCESS" | "ERROR" | "APP_UPDATE_AVAILABLE" | "PLEX_PROCESSING_ERROR" | "ANIME_UPDATE_ERROR";

interface ServiceNotification {
    id: number;
    name: string;
    enabled: boolean;
    type: NotificationType;
    events: NotificationEvent[];
    webhook?: string;
    token?: string;
    api_key?: string;
    channel?: string;
    priority?: number;
    topic?: string;
    host?: string;
    username?: string;
}
