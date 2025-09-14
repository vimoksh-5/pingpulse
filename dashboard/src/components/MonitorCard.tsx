import React from "react";
import { MonitorStatus } from "../types/api";
import { StatusBadge } from "./StatusBadge";
import {
  formatRelativeTime,
  formatDuration,
  parseDuration,
} from "../utils/helpers";
import styles from "./MonitorCard.module.css";

interface MonitorCardProps {
  monitor: MonitorStatus;
  onClick?: () => void;
  className?: string;
}

export const MonitorCard: React.FC<MonitorCardProps> = ({
  monitor,
  onClick,
  className = "",
}) => {
  const latencyMs = parseDuration(monitor.latency);

  const getStatusCodeClass = (statusCode: number) => {
    if (statusCode >= 200 && statusCode < 300) return styles.statusSuccess;
    if (statusCode >= 400) return styles.statusError;
    return styles.statusNeutral;
  };

  const getLatencyClass = (latencyMs: number) => {
    if (latencyMs > 2000) return styles.statusError;
    if (latencyMs > 1000) return styles.statusWarning;
    return styles.statusSuccess;
  };

  const getUptimeClass = (uptime: number) => {
    if (uptime >= 99) return styles.statusSuccess;
    if (uptime >= 95) return styles.statusWarning;
    return styles.statusError;
  };

  return (
    <div
      className={`${styles.card} ${
        onClick ? styles.clickable : ""
      } ${className}`}
      onClick={onClick}
    >
      <div className={styles.content}>
        <div className={styles.main}>
          <div className={styles.header}>
            <h3 className={styles.title}>{monitor.name}</h3>
            <StatusBadge status={monitor.status} size="sm" />
          </div>

          <p className={styles.info}>
            <span className={styles.infoLabel}>Project:</span> {monitor.project}
          </p>

          <p className={`${styles.info} ${styles.url}`}>
            <span className={styles.infoLabel}>URL:</span> {monitor.url}
          </p>

          <div className={styles.metrics}>
            <div className={styles.metric}>
              <span className={styles.metricLabel}>Status Code</span>
              <p
                className={`${styles.metricValue} ${getStatusCodeClass(
                  monitor.statusCode
                )}`}
              >
                {monitor.statusCode}
              </p>
            </div>

            <div className={styles.metric}>
              <span className={styles.metricLabel}>Latency</span>
              <p
                className={`${styles.metricValue} ${getLatencyClass(
                  latencyMs
                )}`}
              >
                {formatDuration(latencyMs)}
              </p>
            </div>

            <div className={styles.metric}>
              <span className={styles.metricLabel}>Uptime</span>
              <p
                className={`${styles.metricValue} ${getUptimeClass(
                  monitor.uptime
                )}`}
              >
                {monitor.uptime.toFixed(1)}%
              </p>
            </div>

            <div className={styles.metric}>
              <span className={styles.metricLabel}>Last Check</span>
              <p className={`${styles.metricValue} ${styles.statusNeutral}`}>
                {formatRelativeTime(monitor.lastCheck)}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MonitorCard;
