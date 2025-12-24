import {Anchor, Badge, Box, Card, Group, Image, Stack, Text} from "@mantine/core";
import {Carousel} from '@mantine/carousel';
import {formatDistanceToNow, parseISO} from "date-fns";
import type { RecentAnimeItem } from "@app/types/Anime";
import { safeDate, statusColor, formatStatusLabel } from "@utils/index";
import classes from './RecentAnimeUpdates.module.css';

export const RecentlyUpdatedAnimeCarousel = ({
    items,
    loading,
}: {
    items?: RecentAnimeItem[];
    loading: boolean;
}) => {
    if (loading) return <Text>Loading...</Text>;

    if (!items || items.length === 0) {
        return <Text>No Updates yet! Start watching some anime and sync your progress to MyAnimeList.</Text>;
    }

    return (
        <Box className={classes.carouselWrapper}>
            <Carousel
                style={{width: "100%"}}
                slideSize={{base: "100%", sm: "50%", md: "33.333%", lg: "25%"}}
                slideGap="md"
                controlSize={40}
                withIndicators
                draggable={true}
                height={420}
                styles={{
                    control: {
                        transition: 'background-color 0.2s ease',
                        backgroundColor: 'rgba(0, 0, 0, 0.5)',
                        borderRadius: '50%',
                        width: 40,
                        height: 40,
                        color: 'white',
                        border: 'none',
                        boxShadow: '0 2px 4px rgba(0, 0, 0, 0.2)',
                        '&:hover': {
                            backgroundColor: 'rgba(0, 0, 0, 0.85)',
                        },
                    },
                }}
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
                                    height={240}
                                    fit="contain"
                                    radius="sm"
                                    mb="sm"
                                    style={{width: "100%", objectFit: "contain"}}
                                />
                            </Anchor>

                            <Group gap="xs" justify="center" mb={4}>
                                <Badge color={statusColor(anime.animeStatus)} variant="light">
                                    {formatStatusLabel(anime.animeStatus, false)}
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
        </Box>
    );
}
