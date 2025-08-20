import {useInfiniteQuery, useQuery} from "@tanstack/react-query";
import React from "react";
import {
    Card,
    Group,
    Stack,
    Flex,
    Text,
    Title,
    Badge,
    Image,
    Anchor,
    Timeline,
    TimelineItem,
    Container,
    Box,
    Code,
    Button
} from "@mantine/core";
import {Carousel} from '@mantine/carousel';
import {MdMovie, MdStar} from "react-icons/md";
import {SiMyanimelist} from "react-icons/si";
import {formatDistanceToNow, parseISO} from "date-fns";
import {
    plexCountsQueryOptions,
    animeUpdateCountQueryOptions,
    recentAnimeUpdatesQueryOptions
} from "@api/queries";
import {APIClient} from "@api/APIClient";
import {PlexKeys} from "@api/query_keys";
import {AuthContext} from "@utils/Context";
import {Navigate} from "@tanstack/react-router";

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
        queryFn: ({ pageParam }) => APIClient.plex.history({ type: "timeline", limit: pageSize, cursor: pageParam as string | undefined }),
        initialPageParam: undefined as string | undefined,
        getNextPageParam: (lastPage: any) => lastPage?.pagination?.next || undefined,
    });

    const timelineItems = (timelineQuery.data?.pages || []).flatMap((p: any) => p?.data || []);

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
                    items={recentAnime as AnimeItem[] | undefined}
                    loading={recentLoading}
                />
            </Stack>

            {/* Timeline Section */}
            <Stack gap="md" mt="xl">
                <Title order={2}>Recent Activity Timeline</Title>
                {timelineQuery.isLoading ? (
                    <Text>Loading timeline...</Text>
                ) : timelineItems && timelineItems.length > 0 ? (
                    <>
                        <Timeline active={-1} bulletSize={24} lineWidth={2}>
                            {timelineItems.map((item: any) => {
                                const plex = item.plex;
                                const status = item.status;
                                const animeUpdate = item.animeUpdate;
                                return (
                                    <TimelineItem key={plex.id}>
                                        <Card shadow="xs" padding="sm" radius="md" withBorder>
                                            <Group>
                                                <Text fw={700}
                                                      >{plex?.Metadata?.grandparentTitle || plex?.Metadata?.title || 'Unknown Title'}</Text>
                                                <Text size="xs" c="dimmed">
                                                    {plex.timestamp ? formatDistanceToNow(new Date(plex.timestamp), {addSuffix: true}) : "-"}
                                                </Text>
                                            </Group>
                                            <Group gap="sm" mb={4} mt="xs">
                                            {plex?.Metadata?.librarySectionTitle && (
                                                    <Badge color="gray"
                                                           variant="light">{plex.Metadata.librarySectionTitle}</Badge>
                                                )}
                                                <Badge color="plex">
                                                    {plex.event}
                                                </Badge>
                                                {status?.success ? (
                                                    <Badge color="green" variant="filled">Successful</Badge>
                                                ) : (
                                                    <Badge color="red" variant="filled">Failed</Badge>
                                                )}
                                            </Group>
                                            {animeUpdate ? (
                                                <Stack gap={2} mt={4}>
                                                    <Text size="sm" fw={700}>MyAnimeList Update Details</Text>
                                                    <Text
                                                        size="sm">Progress: {animeUpdate?.listStatus?.num_episodes_watched}/{animeUpdate?.listDetails?.totalEpisodeNum || '?'}</Text>
                                                    {animeUpdate?.listStatus?.score > 0 && (
                                                        <Text size="sm">Score: {animeUpdate.listStatus.score}</Text>
                                                    )}
                                                    {animeUpdate?.listStatus?.status && (
                                                        <Text size="sm">Status: {formatStatus(animeUpdate.listStatus.status)}</Text>
                                                    )}
                                                    <Anchor href={`https://myanimelist.net/anime/${animeUpdate.malid}`}
                                                            target="_blank" underline="hover" c="mal" fw={700}>
                                                        View on MAL
                                                    </Anchor>
                                                </Stack>
                                            ) : (
                                                status?.errorMsg ? (
                                                    <Stack gap={4} mt={4}>
                                                        <Flex gap="xs" justify="flex-start" align="flex-start" direction="row">
                                                            <Text size="sm" fw={700}>Error:</Text>
                                                            <Code>{status.errorMsg}</Code>
                                                        </Flex>
                                                    </Stack>
                                                ) : (
                                                    <Text size="sm" c="dimmed">No MyAnimeList update for this event.</Text>
                                                )
                                            )}
                                        </Card>
                                    </TimelineItem>
                                );
                            })}
                        </Timeline>
                        {timelineQuery.hasNextPage && (
                            <Group justify="center" mt="sm">
                                <Button onClick={() => timelineQuery.fetchNextPage()} loading={timelineQuery.isFetchingNextPage} variant="light">
                                    {timelineQuery.isFetchingNextPage ? "Loading..." : "Load more"}
                                </Button>
                            </Group>
                        )}
                    </>
                ) : (
                    <Text>No recent activity found.</Text>
                )}
            </Stack>
        </Container>
    );
};


type StatisticsProps = {
    plexCounts?: { countScrobble?: number; countRate?: number };
    animeUpdateCount?: { count?: number };
    plexLoading: boolean;
    animeUpdateLoading: boolean;
};

function StatisticsSection({
                               plexCounts,
                               animeUpdateCount,
                               plexLoading,
                               animeUpdateLoading,
                           }: StatisticsProps) {
    return (
        <Group gap="xl" justify="start">
            <Card shadow="sm" padding="lg" radius="md" withBorder style={{minWidth: 350}}>
                <Group gap="sm" mb="xs" justify="center">
                    <MdMovie size={24}/>
                    <Text size="lg" fw={500}>
                        Plex Scrobbles
                    </Text>
                </Group>

                {plexLoading ? (
                    <Text size="xl" fw={700} c="dimmed" ta="center">
                        Loading...
                    </Text>
                ) : (
                    <Text size="xl" fw={700} ta="center">
                        {plexCounts?.countScrobble?.toLocaleString() || 0}
                    </Text>
                )}

                <Text size="sm" c="dimmed" mt="xs" ta="center">
                    Total Scrobble events
                </Text>
            </Card>

            {/* Plex Ratings */}
            <Card shadow="sm" padding="lg" radius="md" withBorder style={{minWidth: 350}}>
                <Group gap="sm" mb="xs" justify="center">
                    <MdStar size={24}/>
                    <Text size="lg" fw={500}>
                        Plex Ratings
                    </Text>
                </Group>
                {plexLoading ? (
                    <Text size="xl" fw={700} c="dimmed" ta="center">
                        Loading...
                    </Text>
                ) : (
                    <Text size="xl" fw={700} ta="center">
                        {plexCounts?.countRate?.toLocaleString() || 0}
                    </Text>
                )}
                <Text size="sm" c="dimmed" mt="xs" ta="center">
                    Total Rate Events
                </Text>
            </Card>

            {/* Anime Updates */}
            <Card shadow="sm" padding="lg" radius="md" withBorder style={{minWidth: 350}}>
                <Group gap="sm" mb="xs" justify="center">
                    <SiMyanimelist size={30}/>
                    <Text size="lg" fw={500}>
                        Anime Updates
                    </Text>
                </Group>
                {animeUpdateLoading ? (
                    <Text size="xl" fw={700} c="dimmed" ta="center">
                        Loading...
                    </Text>
                ) : (
                    <Text size="xl" fw={700} ta="center">
                        {animeUpdateCount?.count?.toLocaleString() || 0}
                    </Text>
                )}
                <Text size="sm" c="dimmed" mt="xs" ta="center">
                    Successful MAL Updates
                </Text>
            </Card>
        </Group>
    );
}


type AnimeItem = {
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
};

const statusColor = (s: string) =>
    s === 'watching' ? 'green'
        : s === 'completed' ? 'blue'
            : s === 'on_hold' ? 'yellow'
                : s === 'dropped' ? 'red'
                    : s === 'plan_to_watch' ? 'plex'
                        : 'gray';

const safeDate = (d?: string) => d && d.trim().length > 0 ? d : 'Not set';

// Format MAL status: replace underscores with space and capitalize each word
const formatStatus = (s?: string) => {
    if (!s) return "";
    return s
        .split('_')
        .map((w) => (w.length ? w[0].toUpperCase() + w.slice(1) : w))
        .join(' ');
};

function RecentlyUpdatedAnimeCarousel({
                                          items,
                                          loading,
                                      }: {
    items?: AnimeItem[];
    loading: boolean;
}) {
    if (loading) return <Text>Loading...</Text>;

    if (!items || items.length === 0) {
        return <Text>No Updates yet! Start watching some anime and sync your progress to MyAnimeList.</Text>;
    }

    return (
        <Carousel
            style={{width: "100%"}}
            slideSize={{base: "100%", sm: "50%", md: "33.333%", lg: "25%"}}
            slideGap="md"
            controlSize={50}
            withIndicators
            draggable={true}
            height={560}
            emblaOptions={
                {
                    loop: true,
                    dragFree: true,
                    align: 'start',
                }
            }
        >
            {items.map((anime) => (
                <Carousel.Slide key={anime.malId}>
                    {/* Keep each card from stretching too wide */}
                    <Box mx="auto" w="100%" style={{maxWidth: 300}}>
                        <Card shadow="sm" padding="md" radius="md" withBorder>
                            <Anchor href={`https://myanimelist.net/anime/${anime.malId}`} target="_blank"
                                    underline="always">
                                <Image
                                    src={anime.pictureUrl}
                                    alt={anime.title}
                                    height={360}
                                    fit="cover"
                                    radius="sm"
                                    mb="sm"
                                    style={{width: "100%", objectFit: "cover"}}
                                />
                            </Anchor>

                            <Group gap="xs" justify="center" mb={4}>
                                <Badge color={statusColor(anime.animeStatus)} variant="light">
                                    {anime.animeStatus.replace(/_/g, ' ')}
                                </Badge>
                                {anime.rewatchNum > 0 && (
                                    <Badge color="mal" variant="light">{anime.rewatchNum} rewatches</Badge>
                                )}
                            </Group>
                            <Stack justify="center" gap={0}>
                                <Text size="xs" fw={700}>
                                    Progress: {anime.watchedNum}/{anime.totalEpisodeNum === 0 ? "?" : anime.totalEpisodeNum}
                                </Text>
                                <Text size="xs" fw={700}>Rating: {anime.rating > 0 ? anime.rating : "Not set"}</Text>
                                <Text size="xs" fw={700}>Start Date: {safeDate(anime.startDate)}</Text>
                                <Text size="xs" fw={700}>Finish Date: {safeDate(anime.finishDate)}</Text>
                                <Text size="xs" c="dimmed" mt={4}>
                                    Last Updated: {" "}
                                    {anime.lastUpdated ? formatDistanceToNow(parseISO(anime.lastUpdated), {addSuffix: true}) : "-"}
                                </Text>
                            </Stack>
                        </Card>
                    </Box>
                </Carousel.Slide>
            ))}
        </Carousel>
    );
}
