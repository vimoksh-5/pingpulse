import { MonitorResult, MonitorStatus } from "../types/api";

// Format duration string or nanoseconds to milliseconds
export const parseDuration = (
  duration: string | number | undefined | null
): number => {
  // Handle null/undefined
  if (!duration) {
    return 0;
  }

  // If it's a number, assume it's nanoseconds (from Go backend)
  if (typeof duration === "number") {
    return duration / 1000000; // Convert nanoseconds to milliseconds
  }

  // If it's a string, parse the duration format
  if (typeof duration === "string") {
    const match = duration.match(/^(\d+(?:\.\d+)?)(ms|s|m|h)$/);
    if (!match) return 0;

    const value = parseFloat(match[1]);
    const unit = match[2];

    switch (unit) {
      case "ms":
        return value;
      case "s":
        return value * 1000;
      case "m":
        return value * 60 * 1000;
      case "h":
        return value * 60 * 60 * 1000;
      default:
        return 0;
    }
  }

  return 0;
};

// Format milliseconds to human readable duration
export const formatDuration = (ms: number): string => {
  if (ms < 1000) return `${Math.round(ms)}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  if (ms < 3600000) return `${(ms / 60000).toFixed(1)}m`;
  return `${(ms / 3600000).toFixed(1)}h`;
};

// Format timestamp to relative time (epoch timestamps only)
export const formatRelativeTime = (timestamp: string | number): string => {
  const now = new Date();
  const time =
    typeof timestamp === "number"
      ? new Date(timestamp * 1000) // Convert epoch seconds to milliseconds
      : new Date(timestamp); // For backward compatibility with string ISO dates
  const diff = now.getTime() - time.getTime();

  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (seconds < 60) return `${seconds}s ago`;
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  return `${days}d ago`;
};

// Format timestamp to local time (epoch timestamps only)
export const formatLocalTime = (timestamp: string | number): string => {
  const date =
    typeof timestamp === "number"
      ? new Date(timestamp * 1000) // Convert epoch seconds to milliseconds
      : new Date(timestamp); // For backward compatibility with string ISO dates
  return date.toLocaleString();
};

// Determine monitor status based on result
export const getMonitorStatus = (
  result: MonitorResult
): "healthy" | "warning" | "critical" | "unknown" => {
  if (!result.ok) return "critical";

  const latencyMs = parseDuration(result.latency);
  if (latencyMs > 5000) return "critical";
  if (latencyMs > 2000) return "warning";

  if (result.status >= 500) return "critical";
  if (result.status >= 400) return "warning";
  if (result.status >= 200 && result.status < 300) return "healthy";

  return "unknown";
};

// Calculate uptime percentage from results
export const calculateUptime = (results: MonitorResult[]): number => {
  if (results.length === 0) return 0;

  const successfulResults = results.filter((r) => r.ok).length;
  return (successfulResults / results.length) * 100;
};

// Convert MonitorResult to MonitorStatus
export const resultToStatus = (
  result: MonitorResult,
  historicalResults: MonitorResult[] = []
): MonitorStatus => {
  return {
    name: result.name,
    project: result.project,
    url: result.url,
    status: getMonitorStatus(result),
    lastCheck: result.checked_at,
    latency: result.latency,
    uptime: calculateUptime([...historicalResults, result]),
    statusCode: result.status,
  };
};

// Group results by project
export const groupResultsByProject = (
  results: MonitorResult[]
): Record<string, MonitorResult[]> => {
  return results.reduce((acc, result) => {
    if (!acc[result.project]) {
      acc[result.project] = [];
    }
    acc[result.project].push(result);
    return acc;
  }, {} as Record<string, MonitorResult[]>);
};

// Group results by monitor name
export const groupResultsByMonitor = (
  results: MonitorResult[]
): Record<string, MonitorResult[]> => {
  return results.reduce((acc, result) => {
    const key = `${result.project}:${result.name}`;
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(result);
    return acc;
  }, {} as Record<string, MonitorResult[]>);
};

// These functions are no longer needed with CSS Modules
// Status colors are now handled in CSS files

// Debounce function
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;

  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};

// Sort results by timestamp (newest first) - epoch timestamps only
export const sortResultsByTime = (
  results: MonitorResult[]
): MonitorResult[] => {
  return [...results].sort((a, b) => {
    return b.checked_at - a.checked_at; // Both are epoch seconds, newest first
  });
};

// Filter results by time range - epoch timestamps only
export const filterResultsByTimeRange = (
  results: MonitorResult[],
  hours: number
): MonitorResult[] => {
  const cutoffEpoch = Math.floor(Date.now() / 1000) - hours * 60 * 60;
  return results.filter((result) => {
    return result.checked_at >= cutoffEpoch; // Both are epoch seconds
  });
};
