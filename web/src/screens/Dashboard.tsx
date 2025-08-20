import {useQuery, useQueries} from "@tanstack/react-query";
import {
    Card,
    Group,
    Stack,
    Text,
    Title,
    Badge,
    Image,
    Anchor,
    Timeline,
    TimelineItem,
    Container,
    Box
} from "@mantine/core";
import {Carousel} from '@mantine/carousel';
import {MdMovie, MdStar} from "react-icons/md";
import {SiMyanimelist} from "react-icons/si";
import {formatDistanceToNow, parseISO} from "date-fns";
import {
    recentPlexPayloadsQueryOptions,
    animeUpdateByPlexIdQueryOptions,
    plexCountsQueryOptions,
    animeUpdateCountQueryOptions,
    recentAnimeUpdatesQueryOptions
} from "@api/queries";
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
    const {data: recentPlex, isLoading: recentPlexLoading} = useQuery(recentPlexPayloadsQueryOptions(5));

    // Fetch AnimeUpdates for each Plex payload
    const animeUpdateQueries = useQueries({
        queries: (recentPlex || []).map((plex: any) => animeUpdateByPlexIdQueryOptions(plex.id)),
    });

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
                {recentPlexLoading ? (
                    <Text>Loading timeline...</Text>
                ) : recentPlex && recentPlex.length > 0 ? (
                    <Timeline active={-1} bulletSize={24} lineWidth={2}>
                        {recentPlex.map((plex: any, idx: number) => {
                            const animeUpdateQuery = animeUpdateQueries[idx];
                            const animeUpdate = animeUpdateQuery?.data;
                            return (
                                <TimelineItem key={plex.id}>
                                    <Card shadow="xs" padding="sm" radius="md" withBorder>
                                        <Group>
                                            <Text fw={700}
                                                  ml="xs">{plex.metadata?.grandparentTitle || plex.Metadata?.grandparentTitle || plex.grandparentTitle || 'Unknown Title'}</Text>
                                            {plex.Metadata?.librarySectionTitle && (
                                                <Badge color="gray"
                                                       variant="light">{plex.Metadata.librarySectionTitle}</Badge>
                                            )}
                                            <Text size="xs" c="dimmed">
                                                {plex.timestamp ? formatDistanceToNow(new Date(plex.timestamp), {addSuffix: true}) : "-"}
                                            </Text>
                                        </Group>
                                        <Group gap="sm" mb={4} mt="xs">
                                            <Badge color="plex">
                                                {plex.event}
                                            </Badge>
                                            {animeUpdate ? (
                                                <Badge color="green" variant="filled">Successful</Badge>
                                            ) : (
                                                <Badge color="red" variant="filled">Failed</Badge>
                                            )}
                                        </Group>
                                        {animeUpdateQuery.isLoading ? (
                                            <Text size="sm" c="dimmed">Loading update...</Text>
                                        ) : animeUpdate ? (
                                            <Stack gap={2} mt={4}>
                                                <Text size="sm" fw={700}>MyAnimeList Update Details</Text>
                                                <Text
                                                    size="sm">Progress: {animeUpdate?.listStatus?.num_episodes_watched}/{animeUpdate?.listDetails?.totalEpisodeNum || '?'}</Text>
                                                {animeUpdate?.listStatus?.score > 0 && (
                                                    <Text size="sm">Score: {animeUpdate.listStatus.score}</Text>
                                                )}
                                                {animeUpdate?.listStatus?.status && (
                                                    <Text size="sm">Status: {animeUpdate.listStatus.status}</Text>
                                                )}
                                                <Anchor href={`https://myanimelist.net/anime/${animeUpdate.malid}`}
                                                        target="_blank" underline="hover" c="mal" fw={700}>
                                                    View on MAL
                                                </Anchor>
                                            </Stack>
                                        ) : (
                                            <Text size="sm" c="dimmed">No MyAnimeList update for this event.</Text>
                                        )}
                                    </Card>
                                </TimelineItem>
                            );
                        })}
                    </Timeline>
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
            withIndicators
            withControls={false}
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

                            <Group gap="xs" mb={4} mr="auto">
                                <Badge color={statusColor(anime.animeStatus)} variant="light">
                                    {anime.animeStatus}
                                </Badge>
                                {anime.rewatchNum > 0 && (
                                    <Badge color="yellow" variant="light">Rewatched x{anime.rewatchNum}</Badge>
                                )}
                            </Group>

                            <Text size="sm" mb={2}>
                                Progress: {anime.watchedNum}/{anime.totalEpisodeNum === 0 ? "?" : anime.totalEpisodeNum}
                            </Text>
                            <Text size="sm">Rating: {anime.rating > 0 ? anime.rating : "Not set"}</Text>
                            <Text size="xs">Start Date: {safeDate(anime.startDate)}</Text>
                            <Text size="xs">Finish Date: {safeDate(anime.finishDate)}</Text>
                            <Text size="xs" c="dimmed" mt={4}>
                                Last Updated:{" "}
                                {anime.lastUpdated ? formatDistanceToNow(parseISO(anime.lastUpdated), {addSuffix: true}) : "-"}
                            </Text>
                        </Card>
                    </Box>
                </Carousel.Slide>
            ))}
        </Carousel>
    );
}
