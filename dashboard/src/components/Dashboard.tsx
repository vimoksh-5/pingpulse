import React, { useState, useMemo } from "react";
import { usePingPulse } from "../hooks/usePingPulse";
import { StatsCard } from "./StatsCard";
import { MonitorCard } from "./MonitorCard";
import { StatusBadge } from "./StatusBadge";
import { LoadingSpinner } from "./LoadingSpinner";
import { formatDuration, parseDuration } from "../utils/helpers";
import styles from "./Dashboard.module.css";

// Icons (using simple SVG icons)
const ServerIcon = () => (
  <svg
    className={styles.icon}
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2"
    />
  </svg>
);

const ActivityIcon = () => (
  <svg
    className={styles.icon}
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
    />
  </svg>
);

const CheckIcon = () => (
  <svg
    className={styles.icon}
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const ClockIcon = () => (
  <svg
    className={styles.icon}
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const RefreshIcon = () => (
  <svg
    className={styles.iconSmall}
    fill="none"
    stroke="currentColor"
    viewBox="0 0 24 24"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
    />
  </svg>
);

export const Dashboard: React.FC = () => {
  const [state, actions] = usePingPulse(true, 30000);
  const [selectedProject, setSelectedProject] = useState<string>("all");
  const [searchTerm, setSearchTerm] = useState("");

  // Filter monitors based on selected project and search term
  const filteredMonitors = useMemo(() => {
    let filtered = state.latestResults;

    if (selectedProject !== "all") {
      filtered = filtered.filter(
        (monitor) => monitor.project === selectedProject
      );
    }

    if (searchTerm) {
      filtered = filtered.filter(
        (monitor) =>
          monitor.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          monitor.url.toLowerCase().includes(searchTerm.toLowerCase()) ||
          monitor.project.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    return filtered;
  }, [state.latestResults, selectedProject, searchTerm]);

  // Calculate dashboard statistics
  const dashboardStats = useMemo(() => {
    const stats = state.statistics;
    if (!stats) return null;

    // Parse average latency (handles both string and number formats)
    const avgLatencyMs = parseDuration(stats.average_latency);
    const uptime =
      stats.total_monitors > 0
        ? (stats.healthy_monitors / stats.total_monitors) * 100
        : 0;

    return {
      totalProjects: stats.total_projects || 0,
      totalMonitors: stats.total_monitors || 0,
      healthyMonitors: stats.healthy_monitors || 0,
      unhealthyMonitors: stats.unhealthy_monitors || 0,
      uptime: uptime.toFixed(1),
      averageLatency: formatDuration(avgLatencyMs),
    };
  }, [state.statistics]);

  // Group monitors by status for quick overview
  const monitorsByStatus = useMemo(() => {
    const grouped = {
      healthy: 0,
      warning: 0,
      critical: 0,
      unknown: 0,
    };

    filteredMonitors.forEach((monitor) => {
      grouped[monitor.status]++;
    });

    return grouped;
  }, [filteredMonitors]);

  const handleRefresh = () => {
    actions.refresh();
  };

  if (state.loading && !state.statistics) {
    return (
      <div className={`${styles.dashboard} ${styles.loadingCenter}`}>
        <div>
          <LoadingSpinner size="lg" />
          <p className={styles.loadingText}>Loading PingPulse Dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.dashboard}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerInner}>
            <div className={styles.headerText}>
              <h1>PingPulse Dashboard</h1>
              <p>Monitor your services in real-time</p>
            </div>

            <div className={styles.headerActions}>
              {/* Connection Status */}
              <div className={styles.connectionStatus}>
                <div
                  className={`${styles.connectionDot} ${
                    state.isConnected ? styles.connected : styles.disconnected
                  }`}
                />
                <span className={styles.connectionText}>
                  {state.isConnected ? "Connected" : "Disconnected"}
                </span>
              </div>

              {/* Last Updated */}
              {state.lastUpdated && (
                <span className={styles.lastUpdated}>
                  Updated {state.lastUpdated.toLocaleTimeString()}
                </span>
              )}

              {/* Refresh Button */}
              <button
                onClick={handleRefresh}
                disabled={state.loading}
                className={styles.refreshButton}
              >
                <RefreshIcon />
                {state.loading ? "Refreshing..." : "Refresh"}
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className={styles.main}>
        {/* Error Display */}
        {state.error && (
          <div className={styles.errorAlert}>
            <div className={styles.errorHeader}>
              <div>⚠️</div>
              <p className={styles.errorTitle}>Error loading data</p>
            </div>
            <p className={styles.errorMessage}>{state.error}</p>
            <button
              onClick={actions.clearErrors}
              className={styles.errorDismiss}
            >
              Dismiss
            </button>
          </div>
        )}

        {/* Stats Cards */}
        <div className={styles.statsGrid}>
          <StatsCard
            title="Total Projects"
            value={dashboardStats?.totalProjects || 0}
            icon={<ServerIcon />}
            loading={state.statisticsLoading}
          />
          <StatsCard
            title="Total Monitors"
            value={dashboardStats?.totalMonitors || 0}
            subtitle={`${monitorsByStatus.healthy} healthy, ${monitorsByStatus.critical} critical`}
            icon={<ActivityIcon />}
            loading={state.statisticsLoading}
          />
          <StatsCard
            title="Overall Uptime"
            value={`${dashboardStats?.uptime || 0}%`}
            icon={<CheckIcon />}
            loading={state.statisticsLoading}
          />
          <StatsCard
            title="Avg Latency"
            value={dashboardStats?.averageLatency || "0ms"}
            icon={<ClockIcon />}
            loading={state.statisticsLoading}
          />
        </div>

        {/* Filters */}
        <div className={styles.filtersCard}>
          <div className={styles.filtersContent}>
            {/* Project Filter */}
            <div className={styles.filterGroup}>
              <label htmlFor="project-filter" className={styles.filterLabel}>
                Filter by Project
              </label>
              <select
                id="project-filter"
                value={selectedProject}
                onChange={(e) => setSelectedProject(e.target.value)}
                className="form-select"
              >
                <option value="all">All Projects</option>
                {state.projects.map((project) => (
                  <option key={project} value={project}>
                    {project}
                  </option>
                ))}
              </select>
            </div>

            {/* Search */}
            <div className={styles.filterGroup}>
              <label htmlFor="search" className={styles.filterLabel}>
                Search Monitors
              </label>
              <input
                id="search"
                type="text"
                placeholder="Search by name, URL, or project..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="form-input"
              />
            </div>
          </div>

          {/* Status Overview */}
          <div className={styles.statusOverview}>
            <div className={styles.statusItem}>
              <StatusBadge status="healthy" size="sm" />
              <span className={styles.statusCount}>
                {monitorsByStatus.healthy}
              </span>
            </div>
            <div className={styles.statusItem}>
              <StatusBadge status="warning" size="sm" />
              <span className={styles.statusCount}>
                {monitorsByStatus.warning}
              </span>
            </div>
            <div className={styles.statusItem}>
              <StatusBadge status="critical" size="sm" />
              <span className={styles.statusCount}>
                {monitorsByStatus.critical}
              </span>
            </div>
            <div className={styles.statusItem}>
              <StatusBadge status="unknown" size="sm" />
              <span className={styles.statusCount}>
                {monitorsByStatus.unknown}
              </span>
            </div>
          </div>
        </div>

        {/* Monitors Grid */}
        <div className={styles.monitorsSection}>
          <div className={styles.monitorsHeader}>
            <h2 className={styles.monitorsTitle}>
              Monitors ({filteredMonitors.length})
            </h2>
            {state.resultsLoading && <LoadingSpinner size="sm" />}
          </div>

          {filteredMonitors.length === 0 ? (
            <div className={styles.emptyState}>
              <div className={styles.emptyIcon}>
                <ServerIcon />
              </div>
              <h3 className={styles.emptyTitle}>No monitors found</h3>
              <p className={styles.emptyMessage}>
                {searchTerm || selectedProject !== "all"
                  ? "Try adjusting your filters to see more monitors."
                  : "No monitors are currently configured."}
              </p>
            </div>
          ) : (
            <div className={styles.monitorsGrid}>
              {filteredMonitors.map((monitor, index) => (
                <MonitorCard
                  key={`${monitor.project}-${monitor.name}-${index}`}
                  monitor={monitor}
                />
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;
