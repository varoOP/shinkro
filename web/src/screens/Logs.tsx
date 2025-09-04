import {Stack, Text, Group, Select, Button, ActionIcon, Container, Paper, Title, SimpleGrid, ScrollArea, TextInput, Divider} from "@mantine/core";
import {useState, useRef} from "react";
import {TfiReload} from "react-icons/tfi";
import {FaDownload, FaArrowUp, FaArrowDown, FaSearch} from "react-icons/fa";
import { DateTimePicker } from "@mantine/dates";
import dayjs from "dayjs";
import { useQuery } from "@tanstack/react-query";
import { LogContentQueryOptions, LogQueryOptions } from "@api/queries";
import {baseUrl} from "@utils";

const LogViewer = () => {
  const [level, setLevel] = useState<string | null>("ALL");
  const [dateFrom, setDateFrom] = useState<string | null>(null);
  const [dateTo, setDateTo] = useState<string | null>(null);
  const [orderDesc, setOrderDesc] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const logRef = useRef<HTMLDivElement | null>(null);

  const { data: logData, isLoading, refetch } = useQuery(LogContentQueryOptions());

  const formatLogLine = (line: string) => {
    try {
      const obj = JSON.parse(line);
      const timestamp = new Date(obj.time).toLocaleString();
      const level = obj.level?.toUpperCase() || 'UNKNOWN';
      const module = obj.module ? `[${obj.module}]` : '';
      const repo = obj.repo ? `[${obj.repo}]` : '';
      
      // Handle special Plex payload cases
      let message = obj.message || '';
      if (!message && (obj.rawPlexPayload || obj.parsedPlexPayload)) {
        if (obj.parsedPlexPayload) {
          message = `Plex Payload: ${obj.parsedPlexPayload}`;
        } else if (obj.rawPlexPayload) {
          message = `Raw Plex Payload: ${JSON.stringify(obj.rawPlexPayload, null, 2)}`;
        }
      }
      
      // Color coding for levels (using ANSI-like approach for better readability)
      let levelColor = '';
      switch (level.toLowerCase()) {
        case 'error': levelColor = '🔴'; break;
        case 'debug': levelColor = '🔵'; break;
        case 'trace': levelColor = '⚪'; break;
        case 'info': levelColor = '🟢'; break;
        default: levelColor = '⚫';
      }
      
      return `${levelColor} ${timestamp} ${level.padEnd(5)} ${module}${repo} ${message}`;
    } catch {
      return line; // Return original line if not valid JSON
    }
  };

  const filterLog = (log: string) => {
    let lines = log.split("\n").filter(Boolean);
    lines = lines.filter((line) => {
      try {
        const obj = JSON.parse(line);
        if (level && level !== "ALL" && (!obj.level || obj.level.toLowerCase() !== level.toLowerCase())) {
          return false;
        }
        if (dateFrom || dateTo) {
          const logDate = new Date(obj.time);
          if (dateFrom && logDate < new Date(dateFrom)) return false;
          if (dateTo && logDate > new Date(dateTo)) return false;
        }
        if (searchTerm) {
          const searchLower = searchTerm.toLowerCase();
          const messageMatch = obj.message?.toLowerCase().includes(searchLower);
          const plexMatch = obj.rawPlexPayload ? JSON.stringify(obj.rawPlexPayload).toLowerCase().includes(searchLower) : false;
          const parsedPlexMatch = obj.parsedPlexPayload?.toLowerCase().includes(searchLower);
          
          if (!messageMatch && !plexMatch && !parsedPlexMatch) {
            return false;
          }
        }
        return true;
      } catch {
        return false;
      }
    });
    if (orderDesc) lines = lines.reverse();
    return lines.map(formatLogLine).join("\n");
  };

  return (
    <Stack>
      <Group>
        <Text fw={600}>Filter by level:</Text>
        <Select
          data={[
            { value: "ALL", label: "All" },
            { value: "DEBUG", label: "DEBUG" },
            { value: "INFO", label: "INFO" },
            { value: "ERROR", label: "ERROR" },
            { value: "TRACE", label: "TRACE" },
          ]}
          value={level}
          onChange={setLevel}
          size="xs"
        />
        <TextInput
          placeholder="Search logs..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          size="xs"
          w={200}
          leftSection={<FaSearch size={12} />}
        />
        <DateTimePicker
          placeholder="From date and time"
          value={dateFrom}
          onChange={setDateFrom}
          valueFormat="YYYY-MM-DD HH:mm:ss"
          clearable
          size="xs"
          mx={4}
          w={180}
          presets={[
            { value: dayjs().subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss'), label: 'Yesterday' },
            { value: dayjs().format('YYYY-MM-DD HH:mm:ss'), label: 'Today' },
            { value: dayjs().add(1, 'day').format('YYYY-MM-DD HH:mm:ss'), label: 'Tomorrow' },
            { value: dayjs().add(1, 'month').format('YYYY-MM-DD HH:mm:ss'), label: 'Next month' },
            { value: dayjs().add(1, 'year').format('YYYY-MM-DD HH:mm:ss'), label: 'Next year' },
            { value: dayjs().subtract(1, 'month').format('YYYY-MM-DD HH:mm:ss'), label: 'Last month' },
            { value: dayjs().subtract(1, 'year').format('YYYY-MM-DD HH:mm:ss'), label: 'Last year' },
          ]}
        />
        <DateTimePicker
          placeholder="To date and time"
          value={dateTo}
          onChange={setDateTo}
          valueFormat="YYYY-MM-DD HH:mm:ss"
          clearable
          size="xs"
          w={180}
          presets={[
            { value: dayjs().subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss'), label: 'Yesterday' },
            { value: dayjs().format('YYYY-MM-DD HH:mm:ss'), label: 'Today' },
            { value: dayjs().add(1, 'day').format('YYYY-MM-DD HH:mm:ss'), label: 'Tomorrow' },
            { value: dayjs().add(1, 'month').format('YYYY-MM-DD HH:mm:ss'), label: 'Next month' },
            { value: dayjs().add(1, 'year').format('YYYY-MM-DD HH:mm:ss'), label: 'Next year' },
            { value: dayjs().subtract(1, 'month').format('YYYY-MM-DD HH:mm:ss'), label: 'Last month' },
            { value: dayjs().subtract(1, 'year').format('YYYY-MM-DD HH:mm:ss'), label: 'Last year' },
          ]}
        />
      </Group>
      <Group mb={6}>
        <Button
          size="xs"
          variant="outline"
          onClick={() => refetch()}
          loading={isLoading}
          leftSection={<TfiReload size={14} />}
        >
          Reload
        </Button>
        <Button
          size="xs"
          variant="light"
          onClick={() => setOrderDesc((v) => !v)}
          rightSection={orderDesc ? <FaArrowDown size={14} /> : <FaArrowUp size={14} />}
        >
          {orderDesc ? "Newest First" : "Oldest First"}
        </Button>
      </Group>
      <Divider />
      <ScrollArea h={"50vh"} type="scroll">
        <div
          ref={logRef}
          tabIndex={-1}
          style={{
            fontFamily: "monospace",
            padding: 12,
            whiteSpace: "pre",
            fontSize: 13,
            lineHeight: 1.4,
            outline: "none",
            boxSizing: "border-box",
          }}
        >
          {isLoading ? "Loading..." : filterLog(logData || "")}
        </div>
      </ScrollArea>
      <Divider />
    </Stack>
  );
};

const LogFiles = () => {
    const {data: logs} = useQuery(LogQueryOptions());
    const [isDownloading, setIsDownloading] = useState(false);
    const isEmpty = !logs || !(logs.length > 0)

    const handleDownload = async (filename: string) => {
        setIsDownloading(true);
        const response = await fetch(`${baseUrl()}api/fs/logs/${filename}`);
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = filename;
        link.click();
        URL.revokeObjectURL(url);
        setIsDownloading(false);
    }

    return (
        <Stack>
            <Title order={3}>
                Download Log Files
            </Title>
            {isEmpty ? (
                <Text size="md" c="dimmed" ta="center">No Logs Found</Text>
            ) : (
                <SimpleGrid cols={3} spacing="md">
                    {logs.map((log) => (
                     
                            <Stack gap={4}>
                                <Group>
                                    <Text fw={600} size="sm" truncate>
                                        {log.name}
                                    </Text>
                                    <ActionIcon
                                        loading={isDownloading}
                                        onClick={() => handleDownload(log.name)}
                                        variant="outline"
                                        size="sm"
                                    >
                                        <FaDownload size={12} />
                                    </ActionIcon>
                                </Group>
                                <Text size={"xs"} c={"dimmed"}>
                                    Size: {log.size_human}
                                </Text>
                                <Text size={"xs"} c={"dimmed"}>
                                    Last Modified: {new Date(log.modified_at).toLocaleString()}
                                </Text>
                                <Group justify="flex-end" mt="xs">
                                </Group>
                            </Stack>
                    ))}
                </SimpleGrid>
            )}
        </Stack>
    );
};

export const Logs = () => {
    return (
        <Container size={1200} px="md" component="main">
            <Stack gap="md">
                <Title order={1}>
                    Logs
                </Title>
                <Paper mt="xs" withBorder p={"md"} h={"100%"} mih={"80vh"}>
                    <Stack gap="lg">
                        <LogViewer />
                        <LogFiles />
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};
