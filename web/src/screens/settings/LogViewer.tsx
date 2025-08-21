import { useState, useEffect, useRef } from "react";
import { Select, Stack, Group, Text, Button } from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import dayjs from "dayjs";

export const LogViewer = ({ filename }: { filename: string }) => {
  const [level, setLevel] = useState<string | null>("ALL");
  const [dateFrom, setDateFrom] = useState<string | null>(null);
  const [dateTo, setDateTo] = useState<string | null>(null);
  const [live, setLive] = useState(false);
  const [orderDesc, setOrderDesc] = useState(true);
  const [logData, setLogData] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const logRef = useRef<HTMLDivElement | null>(null);

  const fetchLog = async () => {
    if (!filename) return;
    setLoading(true);
    try {
      const res = await fetch(`/api/fs/logs/${filename}`);
      if (!res.ok) throw new Error("Failed to fetch log");
      setLogData(await res.text());
    } catch {
      setLogData("");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLog();
    if (live) {
      intervalRef.current = setInterval(fetchLog, 2000);
    } else if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [filename, live]);

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
        return true;
      } catch {
        return false;
      }
    });
    if (orderDesc) lines = lines.reverse();
    return lines.join("\n");
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
          variant={live ? "filled" : "outline"}
          color={live ? "green" : "gray"}
          onClick={() => setLive((v) => !v)}
        >
          {live ? "Live Log: ON" : "Live Log: OFF"}
        </Button>
        <Button
          size="xs"
          variant={orderDesc ? "filled" : "outline"}
          color={orderDesc ? "blue" : "gray"}
          onClick={() => setOrderDesc((v) => !v)}
        >
          {orderDesc ? "Newest first" : "Oldest first"}
        </Button>
      </Group>
      <div
        ref={logRef}
        tabIndex={-1}
        style={{
          fontFamily: "monospace",
          background: "#18181a",
          color: "#fff",
          borderRadius: 6,
          border: "1px solid #333",
          minHeight: 320,
          maxHeight: 480,
          height: 320,
          overflow: "auto",
          padding: 12,
          whiteSpace: "pre",
          fontSize: 13,
          outline: "none",
          boxSizing: "border-box",
          WebkitOverflowScrolling: "touch",
          pointerEvents: "auto",
        }}
        onWheel={e => e.stopPropagation()}
        onScroll={e => e.stopPropagation()}
      >
        {loading ? "Loading..." : filterLog(logData || "")}
      </div>
    </Stack>
  );
};
