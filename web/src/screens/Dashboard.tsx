import {useQuery} from "@tanstack/react-query";
import {useState} from "react";
import {Container, Stack, Title, Group, Select} from "@mantine/core";
import {
    plexCountsQueryOptions,
    animeUpdateCountQueryOptions,
    recentAnimeUpdatesQueryOptions,
    plexHistoryQueryOptions
} from "@api/queries";
import {AuthContext} from "@utils/Context";
import {Navigate} from "@tanstack/react-router";
import {StatisticsSection, RecentlyUpdatedAnimeCarousel, RecentTimeline} from "@components/dashboard";
import type { RecentAnimeItem } from "@app/types/Anime";

export const Dashboard = () => {
    const isLoggedIn = AuthContext.useSelector((s) => s.isLoggedIn);
    if (!isLoggedIn) {
        return <Navigate to="/login"/>;
    }

    const {data: plexCounts, isLoading: plexLoading} = useQuery(plexCountsQueryOptions());
    const {data: animeUpdateCount, isLoading: animeUpdateLoading} = useQuery(animeUpdateCountQueryOptions());
    const {data: recentAnime, isLoading: recentLoading} = useQuery(recentAnimeUpdatesQueryOptions(8));

    const [limit, setLimit] = useState<number>(5);
    const {data: timelineData, isLoading: timelineLoading} = useQuery(plexHistoryQueryOptions({ limit }));

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
                <Group justify="space-between" align="center">
                    <Title order={2}>Recent Activity Timeline</Title>
                    <Select
                        value={limit.toString()}
                        onChange={(value) => setLimit(value ? parseInt(value, 10) : 5)}
                        data={[
                            { value: "5", label: "5" },
                            { value: "10", label: "10" },
                            { value: "25", label: "25" },
                            { value: "50", label: "50" },
                            { value: "100", label: "100" },
                        ]}
                        style={{ width: 100 }}
                    />
                </Group>
                <RecentTimeline 
                    timelineItems={timelineData || []}
                    isLoading={timelineLoading}
                />
            </Stack>
        </Container>
    );
};
