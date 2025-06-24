import { useQuery, useQueries } from "@tanstack/react-query";
import { Card, Group, Stack, Text, Title, Badge, Image, Anchor, Timeline, TimelineItem } from "@mantine/core";
import { MdMovie, MdStar } from "react-icons/md";
import { SiMyanimelist } from "react-icons/si";
import { formatDistanceToNow, parseISO } from "date-fns";
import { recentPlexPayloadsQueryOptions, animeUpdateByPlexIdQueryOptions, plexCountsQueryOptions, animeUpdateCountQueryOptions, recentAnimeUpdatesQueryOptions } from "@api/queries";
import { AuthContext } from "@utils/Context";
import { Navigate } from "@tanstack/react-router";

export const Dashboard = () => {
    const isLoggedIn = AuthContext.useSelector((s) => s.isLoggedIn);
    if (!isLoggedIn) {
        return <Navigate to="/login" />;
    }

    const { data: plexCounts, isLoading: plexLoading } = useQuery(plexCountsQueryOptions());
    const { data: animeUpdateCount, isLoading: animeUpdateLoading } = useQuery(animeUpdateCountQueryOptions());
    const { data: recentAnime, isLoading: recentLoading } = useQuery(recentAnimeUpdatesQueryOptions(5));
    const { data: recentPlex, isLoading: recentPlexLoading } = useQuery(recentPlexPayloadsQueryOptions(20));

    // Fetch AnimeUpdates for each Plex payload
    const animeUpdateQueries = useQueries({
        queries: (recentPlex || []).map((plex: any) => animeUpdateByPlexIdQueryOptions(plex.id)),
    });

    return (
        <main>
            <Title order={1} mb="lg">Dashboard</Title>
            
            <Stack gap="md">
                <Title order={2} size="h3">Statistics</Title>
                
                <Group gap="md">
                    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ minWidth: 200 }}>
                        <Group gap="sm" mb="xs">
                            <MdMovie size={24} />
                            <Text size="lg" fw={500}>Plex Scrobbles</Text>
                        </Group>
                        
                        {plexLoading ? (
                            <Text size="xl" fw={700} c="dimmed" style={{ textAlign: 'center' }}>Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700} style={{ textAlign: 'center' }}>
                                {plexCounts?.countScrobble?.toLocaleString() || 0}
                            </Text>
                        )}
                        
                        <Text size="sm" c="dimmed" mt="xs">
                            Total media.scrobble events
                        </Text>
                    </Card>

                    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ minWidth: 200 }}>
                        <Group gap="sm" mb="xs">
                            <MdStar size={24} />
                            <Text size="lg" fw={500}>Plex Ratings</Text>
                        </Group>
                        {plexLoading ? (
                            <Text size="xl" fw={700} c="dimmed" style={{ textAlign: 'center' }}>Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700} style={{ textAlign: 'center' }}>
                                {plexCounts?.countRate?.toLocaleString() || 0}
                            </Text>
                        )}
                        <Text size="sm" c="dimmed" mt="xs">
                            Total media.rate events
                        </Text>
                    </Card>

                    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ minWidth: 200 }}>
                        <Group gap="sm" mb="xs">
                            <SiMyanimelist size={30} />
                            <Text size="lg" fw={500}>Anime Updates</Text>
                        </Group>
                        {animeUpdateLoading ? (
                            <Text size="xl" fw={700} c="dimmed" style={{ textAlign: 'center' }}>Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700} style={{ textAlign: 'center' }}>
                                {animeUpdateCount?.count?.toLocaleString() || 0}
                            </Text>
                        )}
                        <Text size="sm" c="dimmed" mt="xs">
                            Total MyAnimeList updates
                        </Text>
                    </Card>
                </Group>
            </Stack>

            <Stack gap="md" mt="xl">
                <Title order={2} size="h3">Recently Updated Anime</Title>
                <Group gap="md">
                    {recentLoading ? (
                        <Text>Loading...</Text>
                    ) : recentAnime && recentAnime.length > 0 ? (
                        recentAnime.map((anime) => (
                            <Card key={anime.malId} shadow="sm" padding="md" radius="md" withBorder style={{ minWidth: 260, maxWidth: 300 }}>
                                <Anchor href={`https://myanimelist.net/anime/${anime.malId}`} target="_blank" underline="always">
                                    <Image src={anime.pictureUrl} alt={anime.title} height={180} width="100%" fit="cover" radius="sm" mb="sm" style={{ objectFit: 'cover', width: '100%', height: 350 }} />
                                </Anchor>
                                <Anchor href={`https://myanimelist.net/anime/${anime.malId}`} target="_blank" underline="hover">
                                    <Text fw={600} size="lg" mb={4} style={{ textAlign: 'center', width: '100%' }}>{anime.title}</Text>
                                </Anchor>
                                <Group gap="xs" mb={4} mr="auto">
                                    <Badge color={
                                        anime.animeStatus === 'watching' ? 'green' :
                                        anime.animeStatus === 'completed' ? 'blue' :
                                        anime.animeStatus === 'on_hold' ? 'yellow' :
                                        anime.animeStatus === 'dropped' ? 'red' :
                                        anime.animeStatus === 'plan_to_watch' ? 'plex' :
                                        'gray'
                                    } variant="light">
                                        {anime.animeStatus}
                                    </Badge>
                                    {anime.rewatchNum > 0 && (
                                        <Badge color="yellow" variant="light">Rewatched x{anime.rewatchNum}</Badge>
                                    )}
                                </Group>
                                <Text size="sm" mb={2}>
                                    Progress: {anime.watchedNum}/{anime.totalEpisodeNum === 0 ? '?' : anime.totalEpisodeNum}
                                </Text>
                                {anime.rating > 0 && (
                                    <Text size="sm" >Rating: {anime.rating}</Text>
                                )}
                                {anime.startDate && (
                                    <Text size="xs" >Start Date: {anime.startDate}</Text>
                                )}
                                {anime.finishDate && (
                                    <Text size="xs">Finish Date: {anime.finishDate}</Text>
                                )}
                                <Text size="xs" c="dimmed" mt={4}>
                                    Last Updated: {anime.lastUpdated ? formatDistanceToNow(parseISO(anime.lastUpdated), { addSuffix: true }) : "-"}
                                </Text>
                            </Card>
                        ))
                    ) : (
                        <Text>No Updates yet! Start watching some anime and sync your progress to MyAnimeList.</Text>
                    )}
                </Group>
            </Stack>

            {/* Timeline Section */}
            <Stack gap="md" mt="xl">
                <Title order={2} size="h3">Recent Activity Timeline</Title>
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
                                        <Text fw={700} ml="xs">{plex.metadata?.grandparentTitle || plex.Metadata?.grandparentTitle || plex.grandparentTitle || 'Unknown Title'}</Text>
                                            {plex.Metadata?.librarySectionTitle && (
                                                <Badge color="gray" variant="light">{plex.Metadata.librarySectionTitle}</Badge>
                                            )}
                                            <Text size="xs" c="dimmed">
                                                {plex.timestamp ? formatDistanceToNow(new Date(plex.timestamp), { addSuffix: true }) : "-"}
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
                                                <Text size="sm">Progress: {animeUpdate?.listStatus?.num_episodes_watched}/{animeUpdate?.listDetails?.totalEpisodeNum || '?'}</Text>
                                                {animeUpdate?.listStatus?.score > 0 && (
                                                    <Text size="sm">Score: {animeUpdate.listStatus.score}</Text>
                                                )}
                                                {animeUpdate?.listStatus?.status && (
                                                   <Text size="sm">Status: {animeUpdate.listStatus.status}</Text>
                                                )}
                                                <Anchor href={`https://myanimelist.net/anime/${animeUpdate.malId}`} target="_blank" underline="hover" c="mal" fw={700}>
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
                    <Text>No recent Plex activity found.</Text>
                )}
            </Stack>
        </main>
    );
};
