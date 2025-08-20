import {Anchor, Badge, Button, Card, Code, Flex, Group, Stack, Text, Timeline, TimelineItem} from "@mantine/core";
import {formatDistanceToNow} from "date-fns";
import {UseInfiniteQueryResult, InfiniteData} from "@tanstack/react-query";
import type { PlexHistoryItem, PlexHistoryResponse } from "@app/types/Plex";
import { formatStatusLabel } from "@utils/index";

export type RecentTimelineProps = {
    timelineQuery: UseInfiniteQueryResult<InfiniteData<PlexHistoryResponse, unknown>, Error>;
    timelineItems: PlexHistoryItem[];
};

export function RecentTimeline({timelineQuery, timelineItems}: RecentTimelineProps) {
    if (timelineQuery.isLoading) {
        return <Text>Loading timeline...</Text>;
    }

    if (!timelineItems || timelineItems.length === 0) {
        return <Text>No recent activity found.</Text>;
    }

    return (
        <>
            <Timeline active={-1} bulletSize={24} lineWidth={2}>
                {timelineItems.map((item) => {
                    const plex = item.plex;
                    const status = item.status;
                    const animeUpdate = item.animeUpdate;
                    const score = animeUpdate?.listStatus?.score;
                    const statusText = animeUpdate?.listStatus?.status;
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
                                        {typeof score === 'number' && score > 0 && (
                                            <Text size="sm">Score: {score}</Text>
                                        )}
                                        {statusText && (
                                            <Text size="sm">Status: {formatStatusLabel(statusText, true)}</Text>
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
    );
}
