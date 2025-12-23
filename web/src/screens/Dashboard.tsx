import {useQuery} from "@tanstack/react-query";
import {Container, Stack, Title, Group, Select} from "@mantine/core";
import {
    plexCountsQueryOptions,
    animeUpdateCountQueryOptions,
    recentAnimeUpdatesQueryOptions,
    plexHistoryQueryOptions
} from "@api/queries";
import {AuthContext, SettingsContext} from "@utils/Context";
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

    const [settings, setSettings] = SettingsContext.use();
    const limit = settings.timelineLimit;
    
    const {data: timelineData, isLoading: timelineLoading} = useQuery(plexHistoryQueryOptions({ limit }));

    return (
        <Container size={1200} px="md" component="main">
            <Stack gap="md" p="md">
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
            <Stack gap="md" mt="xl" p="md">
                <Title order={2}>Recently Updated Anime</Title>
                <RecentlyUpdatedAnimeCarousel
                    items={recentAnime as RecentAnimeItem[] | undefined}
                    loading={recentLoading}
                />
            </Stack>

            {/* Timeline Section */}
            <Stack gap="md" mt="xl" p="md">
                <Group justify="space-between" align="center">
                    <Title order={2}>Recent Activity Timeline</Title>
                    <Select
                        value={limit.toString()}
                        onChange={(value) => {
                            const newLimit = value ? parseInt(value, 10) : 5;
                            if ([5, 10, 25, 50].includes(newLimit)) {
                                setSettings((prevState) => ({
                                    ...prevState,
                                    timelineLimit: newLimit,
                                }));
                            }
                        }}
                        data={[
                            { value: "5", label: "5 items" },
                            { value: "10", label: "10 items" },
                            { value: "25", label: "25 items" },
                            { value: "50", label: "50 items" },
                        ]}
                        style={{ width: 115 }}
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
