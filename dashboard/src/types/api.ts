// API Response Types
export interface StandardResponse<T = any> {
  data: T;
  status_code: number; // 1 for success, 0 for failure
  status_message: string;
}

// Monitor Result Types
export interface MonitorResult {
  project: string;
  name: string;
  url: string;
  ok: boolean;
  status: number;
  latency: string | number; // Can be string (e.g., "1.5s") or number (nanoseconds)
  checked_at: string;
}

export interface QueryFilter {
  project?: string;
  monitor?: string;
  url?: string;
  status?: number;
  ok?: boolean;
  since?: string;
  until?: string;
  limit?: number;
}

export interface ResultsResponse {
  results: MonitorResult[];
  count: number;
  filter: QueryFilter;
}

// Project Types
export interface ProjectsResponse {
  projects: string[];
  count: number;
}

export interface ProjectMonitorsResponse {
  project: string;
  monitors: string[];
  count: number;
}

// Monitor Types
export interface MonitorsResponse {
  monitors: string[];
  count: number;
}

// Statistics Types
export interface ProjectStat {
  monitor_count: number;
  healthy_count: number;
  unhealthy_count: number;
}

export interface Statistics {
  total_projects: number;
  total_monitors: number;
  healthy_monitors: number;
  unhealthy_monitors: number;
  project_stats: Record<string, ProjectStat>;
  status_codes: Record<number, number>;
  average_latency: string | number; // Can be string (e.g., "1.5s") or number (nanoseconds)
}

// Health Check Types
export interface HealthResponse {
  status: string;
}

// Dashboard specific types
export interface DashboardStats {
  totalProjects: number;
  totalMonitors: number;
  healthyMonitors: number;
  unhealthyMonitors: number;
  uptime: number;
  averageLatency: string;
}

export interface AlertSeverity {
  critical: number;
  warning: number;
  info: number;
}

export interface MonitorStatus {
  name: string;
  project: string;
  url: string;
  status: "healthy" | "warning" | "critical" | "unknown";
  lastCheck: string;
  latency: string | number; // Can be string (e.g., "1.5s") or number (nanoseconds)
  uptime: number;
  statusCode: number;
}
