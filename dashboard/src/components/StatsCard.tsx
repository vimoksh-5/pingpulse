import React from "react";
import { LoadingSpinner } from "./LoadingSpinner";
import styles from "./StatsCard.module.css";

interface StatsCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: React.ReactNode;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  loading?: boolean;
  className?: string;
}

export const StatsCard: React.FC<StatsCardProps> = ({
  title,
  value,
  subtitle,
  icon,
  trend,
  loading = false,
  className = "",
}) => {
  return (
    <div className={`${styles.card} ${className}`}>
      <div className={styles.content}>
        <div className={styles.main}>
          <p className={styles.title}>{title}</p>
          {loading ? (
            <div className={styles.loadingContainer}>
              <LoadingSpinner size="sm" />
              <span className={styles.loadingText}>Loading...</span>
            </div>
          ) : (
            <>
              <p className={styles.value}>{value}</p>
              {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
              {trend && (
                <div className={styles.trendContainer}>
                  <span
                    className={`${styles.trendValue} ${
                      trend.isPositive
                        ? styles.trendPositive
                        : styles.trendNegative
                    }`}
                  >
                    {trend.isPositive ? "↗" : "↘"} {Math.abs(trend.value)}%
                  </span>
                  <span className={styles.trendLabel}>vs last hour</span>
                </div>
              )}
            </>
          )}
        </div>
        {icon && (
          <div className={styles.iconContainer}>
            <div className={styles.icon}>{icon}</div>
          </div>
        )}
      </div>
    </div>
  );
};

export default StatsCard;
