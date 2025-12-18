export interface RecentAnimeItem {
    animeStatus: 'watching' | 'completed' | 'on_hold' | 'dropped' | 'plan_to_watch' | string;
    finishDate: string;
    lastUpdated: string;
    malId: number;
    pictureUrl: string;
    rating: number;
    rewatchNum: number;
    startDate: string;
    title: string;
    totalEpisodeNum: number;
    watchedNum: number;
}

export interface TimelineAnimeUpdate {
    malid: number;
    listStatus?: {
        num_episodes_watched?: number;
        score?: number;
        status?: string;
    };
    listDetails?: {
        totalEpisodeNum?: number;
        title?: string;
        pictureUrl?: string;
    };
    // Status fields (consolidated from anime_update_status table)
    status?: AnimeUpdateStatusType;
    errorType?: AnimeUpdateErrorType;
    errorMessage?: string;
    sourceDB?: string;
    sourceID?: number;
    seasonNum?: number;
    episodeNum?: number;
}
