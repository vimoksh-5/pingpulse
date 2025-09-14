import React from "react";
import styles from "./StatusBadge.module.css";

interface StatusBadgeProps {
  status: "healthy" | "warning" | "critical" | "unknown";
  size?: "sm" | "md" | "lg";
  showDot?: boolean;
  className?: string;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({
  status,
  size = "md",
  showDot = true,
  className = "",
}) => {
  const sizeClass = {
    sm: styles.small,
    md: styles.medium,
    lg: styles.large,
  }[size];

  const dotSizeClass = {
    sm: styles.dotSmall,
    md: styles.dotMedium,
    lg: styles.dotLarge,
  }[size];

  const statusClass = {
    healthy: styles.healthy,
    warning: styles.warning,
    critical: styles.critical,
    unknown: styles.unknown,
  }[status];

  const dotStatusClass = {
    healthy: styles.dotHealthy,
    warning: styles.dotWarning,
    critical: styles.dotCritical,
    unknown: styles.dotUnknown,
  }[status];

  const statusText = {
    healthy: "Healthy",
    warning: "Warning",
    critical: "Critical",
    unknown: "Unknown",
  }[status];

  return (
    <span
      className={`${styles.badge} ${sizeClass} ${statusClass} ${className}`}
    >
      {showDot && (
        <span className={`${styles.dot} ${dotSizeClass} ${dotStatusClass}`} />
      )}
      {statusText}
    </span>
  );
};

export default StatusBadge;
