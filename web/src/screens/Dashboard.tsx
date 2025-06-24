import { useQuery } from "@tanstack/react-query";
import { Card, Group, Stack, Text, Title } from "@mantine/core";
import { APIClient } from "@api/APIClient";
import { MdMovie, MdRefresh } from "react-icons/md";

export const Dashboard = () => {
    const { data: scrobbleCount, isLoading: scrobbleLoading } = useQuery({
        queryKey: ["scrobbleCount"],
        queryFn: () => APIClient.plex.getScrobbleCount(),
    });

    const { data: animeUpdateCount, isLoading: animeUpdateLoading } = useQuery({
        queryKey: ["animeUpdateCount"],
        queryFn: () => APIClient.animeupdate.getCount(),
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
                        
                        {scrobbleLoading ? (
                            <Text size="xl" fw={700} c="dimmed">Loading...</Text>
                        ) : (
                            <Text size="xl" fw={700}>
                                {scrobbleCount?.count?.toLocaleString() || 0}
                            </Text>
                        )}
                        
                        <Text size="sm" c="dimmed" mt="xs">
                            Total media scrobble events
                        </Text>
                    </Card>

                    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ minWidth: 200 }}>
                        <Group gap="sm" mb="xs">
                            <MdRefresh size={24} />
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
        </main>
    );
};
