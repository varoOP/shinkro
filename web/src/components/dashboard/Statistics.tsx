import {Card, Group, Text} from "@mantine/core";
import {MdMovie, MdStar} from "react-icons/md";
import {SiMyanimelist} from "react-icons/si";

type StatisticsProps = {
    plexCounts?: { countScrobble?: number; countRate?: number };
    animeUpdateCount?: { count?: number };
    plexLoading: boolean;
    animeUpdateLoading: boolean;
};

export const StatisticsSection = ({
    plexCounts,
    animeUpdateCount,
    plexLoading,
    animeUpdateLoading,
}: StatisticsProps) => {
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
