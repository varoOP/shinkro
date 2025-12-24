import {Suspense} from "react";
import {useSuspenseQuery} from "@tanstack/react-query";
import {Container, Stack, Title, Group, Select, Loader, Center} from "@mantine/core";
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

function StatisticsContent() {
    const {data: plexCounts} = useSuspenseQuery(plexCountsQueryOptions());
    const {data: animeUpdateCount} = useSuspenseQuery(animeUpdateCountQueryOptions());

    return (
        <StatisticsSection
            plexCounts={plexCounts}
            animeUpdateCount={animeUpdateCount}
        />
    );
}

function RecentAnimeContent() {
    const {data: recentAnime} = useSuspenseQuery(recentAnimeUpdatesQueryOptions(8));

    return (
        <RecentlyUpdatedAnimeCarousel
            items={recentAnime as RecentAnimeItem[] | undefined}
        />
    );
}

function TimelineContent() {
    const [settings] = SettingsContext.use();
    const limit = settings.timelineLimit;
    const {data: timelineData} = useSuspenseQuery(plexHistoryQueryOptions({ limit }));

    return (
        <RecentTimeline 
            timelineItems={timelineData || []}
        />
    );
}

export const Dashboard = () => {
    const isLoggedIn = AuthContext.useSelector((s) => s.isLoggedIn);
    if (!isLoggedIn) {
        return <Navigate to="/login"/>;
    }

    const [settings, setSettings] = SettingsContext.use();
    const limit = settings.timelineLimit;

    return (
        <Container size={1200} px="md" component="main">
            <Stack gap="md" p="md">
                <Title order={2}>
                    Statistics
                </Title>
                <Suspense fallback={
                    <Center py="xl">
                        <Loader size="lg" />
                    </Center>
                }>
                    <StatisticsContent />
                </Suspense>
            </Stack>
            <Stack gap="md" mt="xl" p="md">
                <Title order={2}>Recently Updated Anime</Title>
                <Suspense fallback={
                    <Center py="xl">
                        <Loader size="lg" />
                    </Center>
                }>
                    <RecentAnimeContent />
                </Suspense>
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
                <Suspense fallback={
                    <Center py="xl">
                        <Loader size="lg" />
                    </Center>
                }>
                    <TimelineContent />
                </Suspense>
            </Stack>
        </Container>
    );
};
