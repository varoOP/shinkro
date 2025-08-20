import {useInfiniteQuery, useQuery} from "@tanstack/react-query";
import React from "react";
import {Container, Stack, Title} from "@mantine/core";
import {
    plexCountsQueryOptions,
    animeUpdateCountQueryOptions,
    recentAnimeUpdatesQueryOptions
} from "@api/queries";
import {APIClient} from "@api/APIClient";
import {PlexKeys} from "@api/query_keys";
import {AuthContext} from "@utils/Context";
import {Navigate} from "@tanstack/react-router";
import {StatisticsSection, RecentlyUpdatedAnimeCarousel, RecentTimeline} from "@components/dashboard";
import type { RecentAnimeItem } from "@app/types/Anime";
import type { PlexHistoryItem, PlexHistoryResponse } from "@app/types/Plex";

export const Dashboard = () => {
    const isLoggedIn = AuthContext.useSelector((s) => s.isLoggedIn);
    if (!isLoggedIn) {
        return <Navigate to="/login"/>;
    }

    const {data: plexCounts, isLoading: plexLoading} = useQuery(plexCountsQueryOptions());
    const {data: animeUpdateCount, isLoading: animeUpdateLoading} = useQuery(animeUpdateCountQueryOptions());
    const {data: recentAnime, isLoading: recentLoading} = useQuery(recentAnimeUpdatesQueryOptions(8));

    const pageSize = 5;
    const historyParams = React.useMemo(() => ({ limit: pageSize }), [pageSize]);
    const timelineQuery = useInfiniteQuery({
        queryKey: PlexKeys.history("timeline", historyParams),
        queryFn: ({ pageParam }): Promise<PlexHistoryResponse> => APIClient.plex.history({ type: "timeline", limit: pageSize, cursor: pageParam as string | undefined }),
        initialPageParam: undefined as string | undefined,
        getNextPageParam: (lastPage: PlexHistoryResponse) => lastPage?.pagination?.next || undefined,
    });

    const timelineItems: PlexHistoryItem[] = ((timelineQuery.data?.pages as PlexHistoryResponse[] | undefined) || [])
        .flatMap((p) => p?.data || []);

    return (
        <Container size={1200} px="md" component="main">
            <Stack gap="md">
                <Title order={2}>
                    Statistics
                </Title>
                <StatisticsSection
                    plexCounts={plexCounts}
                    animeUpdateCount={animeUpdateCount}
                    plexLoading={plexLoading}
                    animeUpdateLoading={animeUpdateLoading}
                />
            </Stack>
            <Stack gap="md" mt="xl">
                <Title order={2}>Recently Updated Anime</Title>
                <RecentlyUpdatedAnimeCarousel
                    items={recentAnime as RecentAnimeItem[] | undefined}
                    loading={recentLoading}
                />
            </Stack>

            {/* Timeline Section */}
            <Stack gap="md" mt="xl">
                <Title order={2}>Recent Activity Timeline</Title>
                <RecentTimeline 
                    timelineQuery={timelineQuery}
                    timelineItems={timelineItems}
                />
            </Stack>
        </Container>
    );
};
