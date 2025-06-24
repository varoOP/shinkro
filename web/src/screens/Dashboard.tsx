import { useQuery } from "@tanstack/react-query";
import { Card, Group, Stack, Text, Title, Badge, Image, Anchor } from "@mantine/core";
import { APIClient } from "@api/APIClient";
import { MdMovie, MdStar } from "react-icons/md";
import { SiMyanimelist } from "react-icons/si";
import { formatDistanceToNow, parseISO } from "date-fns";

export const Dashboard = () => {
    const { data: plexCounts, isLoading: plexLoading } = useQuery({
        queryKey: ["plexCounts"],
        queryFn: () => APIClient.plex.getCounts(),
    });

    const { data: animeUpdateCount, isLoading: animeUpdateLoading } = useQuery({
        queryKey: ["animeUpdateCount"],
        queryFn: () => APIClient.animeupdate.getCount(),
    });

    const { data: recentAnime, isLoading: recentLoading } = useQuery({
        queryKey: ["recentAnimeUpdates"],
        queryFn: () => APIClient.animeupdate.getRecent(5),
    });

    return (
        <main>
            <Title order={1} mb="lg">Dashboard</Title>
            
            <Stack gap="md">
                <Title order={2} size="h3">Sync Statistics</Title>
                
                <Group gap="md">
                    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ minWidth: 200 }}>
                        <Group gap="sm" mb="xs">
                            <MdMovie size={24} />
                            <Text size="lg" fw={500}>Plex Scrobbles</Text>
                        </Group>
                        
                        {plexLoading ? (
                            <Text size="xl" fw={700} c="dimmed">Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700}>
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
                            <Text size="xl" fw={700} c="dimmed">Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700}>
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
                            <Text size="xl" fw={700} c="dimmed">Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700}>
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
                                    <Image src={anime.pictureUrl} alt={anime.title} height={180} fit="cover" radius="sm" mb="sm" />
                                </Anchor>
                                <Anchor href={`https://myanimelist.net/anime/${anime.malId}`} target="_blank" underline="always">
                                    <Text fw={600} size="lg" mb={4}>{anime.title}</Text>
                                </Anchor>
                                <Group gap="xs" mb={4}>
                                    <Badge color="blue" variant="light">{anime.animeStatus}</Badge>
                                    {anime.rewatchNum > 0 && (
                                        <Badge color="yellow" variant="light">Rewatched x{anime.rewatchNum}</Badge>
                                    )}
                                </Group>
                                <Text size="sm" c="dimmed" mb={2}>
                                    Progress: {anime.watchedNum}/{anime.totalEpisodeNum}
                                </Text>
                                {anime.rating > 0 && (
                                    <Text size="sm" c="yellow.7" mb={2}>Rating: {anime.rating}</Text>
                                )}
                                {anime.startDate && (
                                    <Text size="xs" c="dimmed">Start: {anime.startDate}</Text>
                                )}
                                {anime.finishDate && (
                                    <Text size="xs" c="dimmed">Finish: {anime.finishDate}</Text>
                                )}
                                <Text size="xs" c="dimmed" mt={4}>
                                    Updated {anime.lastUpdated ? formatDistanceToNow(parseISO(anime.lastUpdated), { addSuffix: true }) : "-"}
                                </Text>
                            </Card>
                        ))
                    ) : (
                        <Text>No Updates yet! Start watching some anime and sync your progress to MyAnimeList.</Text>
                    )}
                </Group>
            </Stack>
        </main>
    );
};
