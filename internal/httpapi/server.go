package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pingpulse/internal/monitor"
	"pingpulse/internal/store"
)

// EnhancedSink defines the interface for the enhanced store
type EnhancedSink interface {
	All() map[string][]monitor.Result
	Latest(url string, n int) []monitor.Result
	Query(filter store.QueryFilter) []monitor.Result
	GetProjects() []string
	GetMonitors() []string
	GetMonitorsByProject(project string) []string
}

type Server struct {
	sink     EnhancedSink
	addr     string
	srv      *http.Server
	response *ResponseHelper
}

func NewServer(sink EnhancedSink, addr string) *Server {
	mux := http.NewServeMux()
	s := &Server{
		sink:     sink,
		addr:     addr,
		response: NewResponseHelper(),
	}
	
	// Health check endpoint
	mux.HandleFunc("/healthz", s.handleHealth)
	
	// Legacy results endpoint (backward compatibility)
	mux.HandleFunc("/results", s.handleResults)
	
	// New enhanced API endpoints
	mux.HandleFunc("/api/v1/results", s.handleQueryResults)
	mux.HandleFunc("/api/v1/projects", s.handleProjects)
	mux.HandleFunc("/api/v1/monitors", s.handleMonitors)
	mux.HandleFunc("/api/v1/projects/", s.handleProjectMonitors)
	mux.HandleFunc("/api/v1/statistics", s.handleStatistics)
	mux.HandleFunc("/api/v1/monitors/", s.handleMonitorHistory)
	
	// Wrap with CORS middleware
	corsHandler := s.corsMiddleware(mux)
	
	s.srv = &http.Server{
		Addr:              addr,
		Handler:           corsHandler,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return s
}

func (s *Server) Start() error {
	log.Printf("http listening on %s", s.addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// corsMiddleware adds CORS headers to allow web dashboard access
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.response.Health(w)
}

func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if url == "" {
		_ = enc.Encode(s.sink.All())
		return
	}
	_ = enc.Encode(s.sink.Latest(url, 20))
}

// handleQueryResults handles the enhanced results query endpoint
func (s *Server) handleQueryResults(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := store.QueryFilter{}
	
	if project := r.URL.Query().Get("project"); project != "" {
		filter.Project = project
	}
	
	if monitor := r.URL.Query().Get("monitor"); monitor != "" {
		filter.MonitorName = monitor
	}
	
	if url := r.URL.Query().Get("url"); url != "" {
		filter.URL = url
	}
	
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filter.Status = &status
		} else {
			s.response.BadRequest(w, "Invalid status parameter: must be a valid integer")
			return
		}
	}
	
	if okStr := r.URL.Query().Get("ok"); okStr != "" {
		if ok, err := strconv.ParseBool(okStr); err == nil {
			filter.OK = &ok
		} else {
			s.response.BadRequest(w, "Invalid ok parameter: must be true or false")
			return
		}
	}
	
	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		if since, err := strconv.ParseInt(sinceStr, 10, 64); err == nil {
			filter.Since = &since
		} else {
			s.response.BadRequest(w, "Invalid since parameter: must be epoch timestamp in seconds (e.g., 1757849600)")
			return
		}
	}
	
	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		if until, err := strconv.ParseInt(untilStr, 10, 64); err == nil {
			filter.Until = &until
		} else {
			s.response.BadRequest(w, "Invalid until parameter: must be epoch timestamp in seconds (e.g., 1757849600)")
			return
		}
	}
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		} else {
			s.response.BadRequest(w, "Invalid limit parameter: must be a positive integer")
			return
		}
	}
	
	results := s.sink.Query(filter)
	
	responseData := map[string]interface{}{
		"results": results,
		"count":   len(results),
		"filter":  filter,
	}
	
	s.response.Success(w, responseData, "Results retrieved successfully")
}

// handleProjects returns all available projects
func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	projects := s.sink.GetProjects()
	
	responseData := map[string]interface{}{
		"projects": projects,
		"count":    len(projects),
	}
	
	s.response.Success(w, responseData, "Projects retrieved successfully")
}

// handleMonitors returns all available monitors
func (s *Server) handleMonitors(w http.ResponseWriter, r *http.Request) {
	monitors := s.sink.GetMonitors()
	
	responseData := map[string]interface{}{
		"monitors": monitors,
		"count":    len(monitors),
	}
	
	s.response.Success(w, responseData, "Monitors retrieved successfully")
}

// handleProjectMonitors returns monitors for a specific project
func (s *Server) handleProjectMonitors(w http.ResponseWriter, r *http.Request) {
	// Extract project name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	pathParts := strings.Split(path, "/")
	
	if len(pathParts) < 2 || pathParts[1] != "monitors" {
		s.response.BadRequest(w, "Invalid path. Use /api/v1/projects/{project}/monitors")
		return
	}
	
	project := pathParts[0]
	if project == "" {
		s.response.BadRequest(w, "Project name is required")
		return
	}
	
	monitors := s.sink.GetMonitorsByProject(project)
	
	responseData := map[string]interface{}{
		"project":  project,
		"monitors": monitors,
		"count":    len(monitors),
	}
	
	s.response.Success(w, responseData, "Project monitors retrieved successfully")
}

// Statistics represents monitoring statistics
type Statistics struct {
	TotalProjects    int                    `json:"total_projects"`
	TotalMonitors    int                    `json:"total_monitors"`
	HealthyMonitors  int                    `json:"healthy_monitors"`
	UnhealthyMonitors int                   `json:"unhealthy_monitors"`
	ProjectStats     map[string]ProjectStat `json:"project_stats"`
	StatusCodes      map[int]int            `json:"status_codes"`
	AverageLatency   time.Duration          `json:"average_latency"`
}

type ProjectStat struct {
	MonitorCount      int `json:"monitor_count"`
	HealthyCount      int `json:"healthy_count"`
	UnhealthyCount    int `json:"unhealthy_count"`
}

// handleStatistics returns monitoring statistics
func (s *Server) handleStatistics(w http.ResponseWriter, r *http.Request) {
	// Get all results
	allResults := s.sink.All()
	
	stats := Statistics{
		ProjectStats: make(map[string]ProjectStat),
		StatusCodes:  make(map[int]int),
	}
	
	projects := make(map[string]bool)
	monitors := make(map[string]bool)
	var totalLatency time.Duration
	var resultCount int
	
	// Process all results
	for _, results := range allResults {
		if len(results) == 0 {
			continue
		}
		
		// Get latest result for each URL
		latest := results[len(results)-1]
		
		// Count projects and monitors
		projects[latest.Project] = true
		monitors[latest.Name] = true
		
		// Update project stats
		projectStat := stats.ProjectStats[latest.Project]
		projectStat.MonitorCount++
		if latest.OK {
			projectStat.HealthyCount++
			stats.HealthyMonitors++
		} else {
			projectStat.UnhealthyCount++
			stats.UnhealthyMonitors++
		}
		stats.ProjectStats[latest.Project] = projectStat
		
		// Count status codes
		stats.StatusCodes[latest.Status]++
		
		// Calculate average latency
		totalLatency += latest.Latency
		resultCount++
	}
	
	stats.TotalProjects = len(projects)
	stats.TotalMonitors = len(monitors)
	
	if resultCount > 0 {
		stats.AverageLatency = totalLatency / time.Duration(resultCount)
	}
	
	s.response.Success(w, stats, "Statistics retrieved successfully")
}

// handleMonitorHistory returns historical data for a specific monitor
func (s *Server) handleMonitorHistory(w http.ResponseWriter, r *http.Request) {
	// Extract monitor name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/monitors/")
	pathParts := strings.Split(path, "/")
	
	if len(pathParts) < 2 || pathParts[1] != "history" {
		s.response.BadRequest(w, "Invalid path. Use /api/v1/monitors/{monitor}/history")
		return
	}
	
	monitorName := pathParts[0]
	if monitorName == "" {
		s.response.BadRequest(w, "Monitor name is required")
		return
	}
	
	// Parse query parameters for time range
	filter := store.QueryFilter{
		MonitorName: monitorName,
	}
	
	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		if since, err := strconv.ParseInt(sinceStr, 10, 64); err == nil {
			filter.Since = &since
		} else {
			s.response.BadRequest(w, "Invalid since parameter: must be epoch timestamp in seconds (e.g., 1757849600)")
			return
		}
	}
	
	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		if until, err := strconv.ParseInt(untilStr, 10, 64); err == nil {
			filter.Until = &until
		} else {
			s.response.BadRequest(w, "Invalid until parameter: must be epoch timestamp in seconds (e.g., 1757849600)")
			return
		}
	}
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		} else {
			s.response.BadRequest(w, "Invalid limit parameter: must be a positive integer")
			return
		}
	}
	
	// If no time range specified, default to last 90 days
	if filter.Since == nil {
		since := time.Now().AddDate(0, 0, -90).Unix()
		filter.Since = &since
	}
	
	results := s.sink.Query(filter)
	
	responseData := map[string]interface{}{
		"monitor": monitorName,
		"results": results,
		"count":   len(results),
		"filter":  filter,
	}
	
	s.response.Success(w, responseData, "Monitor history retrieved successfully")
}