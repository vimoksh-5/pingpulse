import React, { useState, useMemo } from "react";
import { MonitorResult } from "../types/api";
import { formatRelativeTime } from "../utils/helpers";
import styles from "./UptimeGraph.module.css";

interface UptimeGraphProps {
  results: MonitorResult[];
  days?: number;
  className?: string;
}

interface UptimeSegment {
  status: "up" | "down" | "degraded" | "unknown";
  timestamp: Date;
  duration: number; // in minutes
  tooltip: string;
}

export const UptimeGraph: React.FC<UptimeGraphProps> = ({
  results,
  days = 90,
  className = "",
}) => {
  const [hoveredSegment, setHoveredSegment] = useState<UptimeSegment | null>(
    null
  );
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

  // Process results into uptime segments
  const uptimeSegments = useMemo(() => {
    if (!results || results.length === 0) return [];

    // Sort results by timestamp (oldest first)
    const sortedResults = [...results].sort((a, b) => {
      return a.checked_at - b.checked_at; // Both are epoch seconds
    });

    const now = new Date();
    const startDate = new Date(now.getTime() - days * 24 * 60 * 60 * 1000);

    // Filter results to only include those within the time range
    const filteredResults = sortedResults.filter((result) => {
      const resultTime = new Date(result.checked_at * 1000); // Convert epoch seconds to milliseconds
      return resultTime >= startDate;
    });

    if (filteredResults.length === 0) return [];

    const segments: UptimeSegment[] = [];
    const totalMinutes = days * 24 * 60;
    const segmentDuration = Math.max(1, Math.floor(totalMinutes / 90)); // Aim for ~90 segments

    // Group results into time segments
    for (let i = 0; i < totalMinutes; i += segmentDuration) {
      const segmentStart = new Date(startDate.getTime() + i * 60 * 1000);
      const segmentEnd = new Date(
        segmentStart.getTime() + segmentDuration * 60 * 1000
      );

      // Find results within this segment
      const segmentResults = filteredResults.filter((result) => {
        const resultTime = new Date(result.checked_at * 1000); // Convert epoch seconds to milliseconds
        return resultTime >= segmentStart && resultTime < segmentEnd;
      });

      let status: UptimeSegment["status"] = "unknown";
      let tooltip = `${formatRelativeTime(
        segmentStart.toISOString()
      )} - No data`;

      if (segmentResults.length > 0) {
        // Determine segment status based on results
        const upCount = segmentResults.filter((r) => r.ok).length;
        const totalCount = segmentResults.length;
        const upPercentage = (upCount / totalCount) * 100;

        if (upPercentage === 100) {
          status = "up";
          tooltip = `${formatRelativeTime(
            segmentStart.toISOString()
          )} - All checks passed (${totalCount} checks)`;
        } else if (upPercentage >= 50) {
          status = "degraded";
          tooltip = `${formatRelativeTime(
            segmentStart.toISOString()
          )} - Partial outage (${upCount}/${totalCount} checks passed)`;
        } else if (upPercentage > 0) {
          status = "degraded";
          tooltip = `${formatRelativeTime(
            segmentStart.toISOString()
          )} - Major outage (${upCount}/${totalCount} checks passed)`;
        } else {
          status = "down";
          tooltip = `${formatRelativeTime(
            segmentStart.toISOString()
          )} - Complete outage (0/${totalCount} checks passed)`;
        }
      }

      segments.push({
        status,
        timestamp: segmentStart,
        duration: segmentDuration,
        tooltip,
      });
    }

    return segments;
  }, [results, days]);

  // Calculate uptime percentage
  const uptimePercentage = useMemo(() => {
    if (uptimeSegments.length === 0) return 0;

    const upSegments = uptimeSegments.filter(
      (segment) => segment.status === "up"
    ).length;
    return (upSegments / uptimeSegments.length) * 100;
  }, [uptimeSegments]);

  const handleMouseEnter = (
    segment: UptimeSegment,
    event: React.MouseEvent
  ) => {
    setHoveredSegment(segment);
    setMousePosition({ x: event.clientX, y: event.clientY });
  };

  const handleMouseMove = (event: React.MouseEvent) => {
    setMousePosition({ x: event.clientX, y: event.clientY });
  };

  const handleMouseLeave = () => {
    setHoveredSegment(null);
  };

  if (uptimeSegments.length === 0) {
    return (
      <div className={`${styles.container} ${className}`}>
        <div className={styles.header}>
          <span className={styles.label}>Uptime ({days} days)</span>
          <span className={styles.percentage}>No data</span>
        </div>
        <div className={styles.graph}>
          <div className={styles.noData}>No uptime data available</div>
        </div>
      </div>
    );
  }

  return (
    <div className={`${styles.container} ${className}`}>
      <div className={styles.header}>
        <span className={styles.label}>Uptime ({days} days)</span>
        <span className={styles.percentage}>
          {uptimePercentage.toFixed(2)}%
        </span>
      </div>

      <div className={styles.graph}>
        <div className={styles.segments}>
          {uptimeSegments.map((segment, index) => (
            <div
              key={index}
              className={`${styles.segment} ${
                styles[
                  `segment${
                    segment.status.charAt(0).toUpperCase() +
                    segment.status.slice(1)
                  }`
                ]
              }`}
              onMouseEnter={(e) => handleMouseEnter(segment, e)}
              onMouseMove={handleMouseMove}
              onMouseLeave={handleMouseLeave}
            />
          ))}
        </div>

        <div className={styles.timeLabels}>
          <span className={styles.timeLabel}>{days} days ago</span>
          <span className={styles.timeLabel}>Today</span>
        </div>
      </div>

      {/* Tooltip */}
      {hoveredSegment && (
        <div
          className={styles.tooltip}
          style={{
            left: mousePosition.x + 10,
            top: mousePosition.y - 10,
          }}
        >
          {hoveredSegment.tooltip}
        </div>
      )}
    </div>
  );
};

export default UptimeGraph;
