import axios, { AxiosInstance, AxiosResponse } from "axios";
import {
  StandardResponse,
  ResultsResponse,
  ProjectsResponse,
  MonitorsResponse,
  ProjectMonitorsResponse,
  Statistics,
  HealthResponse,
  QueryFilter,
  MonitorResult,
} from "../types/api";

class PingPulseAPI {
  private client: AxiosInstance;

  constructor(
    baseURL: string = process.env.REACT_APP_API_URL || "http://localhost:8080"
  ) {
    this.client = axios.create({
      baseURL,
      timeout: 10000,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        console.error("API Error:", error);
        throw error;
      }
    );
  }

  // Health check
  async getHealth(): Promise<HealthResponse> {
    const response: AxiosResponse<StandardResponse<HealthResponse>> =
      await this.client.get("/healthz");
    return response.data.data;
  }

  // Get all projects
  async getProjects(): Promise<string[]> {
    const response: AxiosResponse<StandardResponse<ProjectsResponse>> =
      await this.client.get("/api/v1/projects");
    return response.data.data.projects;
  }

  // Get all monitors
  async getMonitors(): Promise<string[]> {
    const response: AxiosResponse<StandardResponse<MonitorsResponse>> =
      await this.client.get("/api/v1/monitors");
    return response.data.data.monitors;
  }

  // Get monitors for a specific project
  async getProjectMonitors(project: string): Promise<string[]> {
    const encodedProject = encodeURIComponent(project);
    const response: AxiosResponse<StandardResponse<ProjectMonitorsResponse>> =
      await this.client.get(`/api/v1/projects/${encodedProject}/monitors`);
    return response.data.data.monitors;
  }

  // Get monitoring results with optional filters
  async getResults(filter?: QueryFilter): Promise<ResultsResponse> {
    const params = new URLSearchParams();

    if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }

    const response: AxiosResponse<StandardResponse<ResultsResponse>> =
      await this.client.get(`/api/v1/results?${params.toString()}`);
    return response.data.data;
  }

  // Get statistics
  async getStatistics(): Promise<Statistics> {
    const response: AxiosResponse<StandardResponse<Statistics>> =
      await this.client.get("/api/v1/statistics");
    return response.data.data;
  }

  // Get latest results for all monitors
  async getLatestResults(): Promise<MonitorResult[]> {
    const results = await this.getResults({ limit: 1000 });

    // Group by monitor and get the latest result for each
    const latestByMonitor = new Map<string, MonitorResult>();

    results.results.forEach((result) => {
      const key = `${result.project}:${result.name}`;
      const existing = latestByMonitor.get(key);

      // Handle epoch timestamps (always numbers now)
      const resultTime = result.checked_at * 1000; // Convert epoch seconds to milliseconds
      const existingTime = existing ? existing.checked_at * 1000 : 0;

      if (!existing || resultTime > existingTime) {
        latestByMonitor.set(key, result);
      }
    });

    // Convert to array and sort by config_order if available, otherwise preserve API order
    const latestResults = Array.from(latestByMonitor.values());

    // Sort by config_order if available (backend should already provide this order)
    if (latestResults.some((r) => r.config_order !== undefined)) {
      latestResults.sort(
        (a, b) => (a.config_order || 0) - (b.config_order || 0)
      );
    }

    return latestResults;
  }

  // Get results for a specific project
  async getProjectResults(
    project: string,
    limit?: number
  ): Promise<MonitorResult[]> {
    const results = await this.getResults({ project, limit });
    return results.results;
  }

  // Get results for a specific monitor
  async getMonitorResults(
    monitor: string,
    limit?: number
  ): Promise<MonitorResult[]> {
    const results = await this.getResults({ monitor, limit });
    return results.results;
  }

  // Get failed results
  async getFailedResults(limit?: number): Promise<MonitorResult[]> {
    const results = await this.getResults({ ok: false, limit });
    return results.results;
  }

  // Get results within a time range
  async getResultsInTimeRange(
    since: string,
    until?: string,
    limit?: number
  ): Promise<MonitorResult[]> {
    const results = await this.getResults({ since, until, limit });
    return results.results;
  }

  // Get results by status code
  async getResultsByStatus(
    status: number,
    limit?: number
  ): Promise<MonitorResult[]> {
    const results = await this.getResults({ status, limit });
    return results.results;
  }

  // Get historical data for a specific monitor
  async getMonitorHistory(
    monitor: string,
    since?: string,
    until?: string,
    limit?: number
  ): Promise<MonitorResult[]> {
    const params = new URLSearchParams();

    if (since) params.append("since", since);
    if (until) params.append("until", until);
    if (limit) params.append("limit", limit.toString());

    const encodedMonitor = encodeURIComponent(monitor);
    const response: AxiosResponse<StandardResponse<ResultsResponse>> =
      await this.client.get(
        `/api/v1/monitors/${encodedMonitor}/history?${params.toString()}`
      );
    return response.data.data.results;
  }

  // Update base URL (useful for configuration changes)
  updateBaseURL(baseURL: string): void {
    this.client.defaults.baseURL = baseURL;
  }

  // Get client instance for custom requests
  getClient(): AxiosInstance {
    return this.client;
  }
}

// Create and export a singleton instance
export const pingPulseAPI = new PingPulseAPI();
export default PingPulseAPI;
