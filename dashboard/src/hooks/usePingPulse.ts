import { useState, useEffect, useCallback } from "react";
import { pingPulseAPI } from "../services/api";
import {
  MonitorResult,
  Statistics,
  QueryFilter,
  MonitorStatus,
} from "../types/api";
import { resultToStatus, sortResultsByTime } from "../utils/helpers";

export interface UsePingPulseState {
  // Data
  projects: string[];
  monitors: string[];
  results: MonitorResult[];
  statistics: Statistics | null;
  latestResults: MonitorStatus[];

  // Loading states
  loading: boolean;
  projectsLoading: boolean;
  monitorsLoading: boolean;
  resultsLoading: boolean;
  statisticsLoading: boolean;

  // Error states
  error: string | null;
  projectsError: string | null;
  monitorsError: string | null;
  resultsError: string | null;
  statisticsError: string | null;

  // Connection status
  isConnected: boolean;
  lastUpdated: Date | null;
}

export interface UsePingPulseActions {
  // Data fetching
  fetchProjects: () => Promise<void>;
  fetchMonitors: () => Promise<void>;
  fetchResults: (filter?: QueryFilter) => Promise<void>;
  fetchStatistics: () => Promise<void>;
  fetchLatestResults: () => Promise<void>;

  // Specific queries
  fetchProjectResults: (project: string) => Promise<void>;
  fetchMonitorResults: (monitor: string) => Promise<void>;
  fetchFailedResults: () => Promise<void>;

  // Utility
  refresh: () => Promise<void>;
  clearErrors: () => void;
  updateBaseURL: (url: string) => void;
}

export const usePingPulse = (
  autoRefresh = true,
  refreshInterval = 30000
): [UsePingPulseState, UsePingPulseActions] => {
  const [state, setState] = useState<UsePingPulseState>({
    projects: [],
    monitors: [],
    results: [],
    statistics: null,
    latestResults: [],

    loading: false,
    projectsLoading: false,
    monitorsLoading: false,
    resultsLoading: false,
    statisticsLoading: false,

    error: null,
    projectsError: null,
    monitorsError: null,
    resultsError: null,
    statisticsError: null,

    isConnected: false,
    lastUpdated: null,
  });

  // Fetch projects
  const fetchProjects = useCallback(async () => {
    setState((prev) => ({
      ...prev,
      projectsLoading: true,
      projectsError: null,
    }));

    try {
      const projects = await pingPulseAPI.getProjects();
      setState((prev) => ({
        ...prev,
        projects,
        projectsLoading: false,
        isConnected: true,
        lastUpdated: new Date(),
      }));
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to fetch projects";
      setState((prev) => ({
        ...prev,
        projectsError: errorMessage,
        projectsLoading: false,
        isConnected: false,
      }));
    }
  }, []);

  // Fetch monitors
  const fetchMonitors = useCallback(async () => {
    setState((prev) => ({
      ...prev,
      monitorsLoading: true,
      monitorsError: null,
    }));

    try {
      const monitors = await pingPulseAPI.getMonitors();
      setState((prev) => ({
        ...prev,
        monitors,
        monitorsLoading: false,
        isConnected: true,
        lastUpdated: new Date(),
      }));
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to fetch monitors";
      setState((prev) => ({
        ...prev,
        monitorsError: errorMessage,
        monitorsLoading: false,
        isConnected: false,
      }));
    }
  }, []);

  // Fetch results
  const fetchResults = useCallback(async (filter?: QueryFilter) => {
    setState((prev) => ({ ...prev, resultsLoading: true, resultsError: null }));

    try {
      const response = await pingPulseAPI.getResults(filter);
      const sortedResults = sortResultsByTime(response.results);

      setState((prev) => ({
        ...prev,
        results: sortedResults,
        resultsLoading: false,
        isConnected: true,
        lastUpdated: new Date(),
      }));
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to fetch results";
      setState((prev) => ({
        ...prev,
        resultsError: errorMessage,
        resultsLoading: false,
        isConnected: false,
      }));
    }
  }, []);

  // Fetch statistics
  const fetchStatistics = useCallback(async () => {
    setState((prev) => ({
      ...prev,
      statisticsLoading: true,
      statisticsError: null,
    }));

    try {
      const statistics = await pingPulseAPI.getStatistics();
      setState((prev) => ({
        ...prev,
        statistics,
        statisticsLoading: false,
        isConnected: true,
        lastUpdated: new Date(),
      }));
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to fetch statistics";
      setState((prev) => ({
        ...prev,
        statisticsError: errorMessage,
        statisticsLoading: false,
        isConnected: false,
      }));
    }
  }, []);

  // Fetch latest results
  const fetchLatestResults = useCallback(async () => {
    setState((prev) => ({ ...prev, resultsLoading: true, resultsError: null }));

    try {
      const results = await pingPulseAPI.getLatestResults();
      const latestResults = results.map((result) => resultToStatus(result));

      setState((prev) => ({
        ...prev,
        latestResults,
        resultsLoading: false,
        isConnected: true,
        lastUpdated: new Date(),
      }));
    } catch (error) {
      const errorMessage =
        error instanceof Error
          ? error.message
          : "Failed to fetch latest results";
      setState((prev) => ({
        ...prev,
        resultsError: errorMessage,
        resultsLoading: false,
        isConnected: false,
      }));
    }
  }, []);

  // Fetch project results
  const fetchProjectResults = useCallback(
    async (project: string) => {
      await fetchResults({ project, limit: 100 });
    },
    [fetchResults]
  );

  // Fetch monitor results
  const fetchMonitorResults = useCallback(
    async (monitor: string) => {
      await fetchResults({ monitor, limit: 100 });
    },
    [fetchResults]
  );

  // Fetch failed results
  const fetchFailedResults = useCallback(async () => {
    await fetchResults({ ok: false, limit: 100 });
  }, [fetchResults]);

  // Refresh all data
  const refresh = useCallback(async () => {
    setState((prev) => ({ ...prev, loading: true, error: null }));

    try {
      await Promise.all([
        fetchProjects(),
        fetchMonitors(),
        fetchStatistics(),
        fetchLatestResults(),
      ]);

      setState((prev) => ({ ...prev, loading: false }));
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to refresh data";
      setState((prev) => ({
        ...prev,
        error: errorMessage,
        loading: false,
        isConnected: false,
      }));
    }
  }, [fetchProjects, fetchMonitors, fetchStatistics, fetchLatestResults]);

  // Clear errors
  const clearErrors = useCallback(() => {
    setState((prev) => ({
      ...prev,
      error: null,
      projectsError: null,
      monitorsError: null,
      resultsError: null,
      statisticsError: null,
    }));
  }, []);

  // Update base URL
  const updateBaseURL = useCallback(
    (url: string) => {
      pingPulseAPI.updateBaseURL(url);
      refresh();
    },
    [refresh]
  );

  // Auto-refresh effect
  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(() => {
      refresh();
    }, refreshInterval);

    return () => clearInterval(interval);
  }, [autoRefresh, refreshInterval, refresh]);

  // Initial data fetch
  useEffect(() => {
    refresh();
  }, [refresh]);

  const actions: UsePingPulseActions = {
    fetchProjects,
    fetchMonitors,
    fetchResults,
    fetchStatistics,
    fetchLatestResults,
    fetchProjectResults,
    fetchMonitorResults,
    fetchFailedResults,
    refresh,
    clearErrors,
    updateBaseURL,
  };

  return [state, actions];
};

export default usePingPulse;
