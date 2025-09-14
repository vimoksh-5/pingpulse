package store

import (
	"sort"
	"strings"
	"sync"
	"time"

	"pingpulse/internal/monitor"
)

// QueryFilter represents filters for querying results
type QueryFilter struct {
	Project     string    // Filter by project name
	MonitorName string    // Filter by monitor name
	URL         string    // Filter by URL
	Status      *int      // Filter by HTTP status code
	OK          *bool     // Filter by OK status
	Since       time.Time // Filter results after this time
	Until       time.Time // Filter results before this time
	Limit       int       // Limit number of results
}

type MemoryStore struct {
	mu     sync.RWMutex
	cap    int
	data   map[string][]monitor.Result
	// Additional indexes for faster querying
	projectIndex map[string][]string // project -> list of URLs
	monitorIndex map[string]string   // monitor name -> URL
}

func NewMemoryStore(capacity int) *MemoryStore {
	return &MemoryStore{
		cap:          capacity,
		data:         make(map[string][]monitor.Result),
		projectIndex: make(map[string][]string),
		monitorIndex: make(map[string]string),
	}
}

func (s *MemoryStore) Append(r monitor.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	buf := append(s.data[r.URL], r)
	if len(buf) > s.cap {
		buf = buf[len(buf)-s.cap:]
	}
	s.data[r.URL] = buf
	
	// Update indexes
	s.updateIndexes(r)
}

func (s *MemoryStore) updateIndexes(r monitor.Result) {
	// Update project index
	if r.Project != "" {
		urls := s.projectIndex[r.Project]
		found := false
		for _, url := range urls {
			if url == r.URL {
				found = true
				break
			}
		}
		if !found {
			s.projectIndex[r.Project] = append(urls, r.URL)
		}
	}
	
	// Update monitor name index
	if r.Name != "" {
		s.monitorIndex[r.Name] = r.URL
	}
}

func (s *MemoryStore) Latest(url string, n int) []monitor.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	buf := s.data[url]
	if n <= 0 || n > len(buf) {
		n = len(buf)
	}
	start := len(buf) - n
	return append([]monitor.Result(nil), buf[start:]...)
}

func (s *MemoryStore) All() map[string][]monitor.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]monitor.Result, len(s.data))
	for k, v := range s.data {
		cp := make([]monitor.Result, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

// Query returns results based on the provided filter
func (s *MemoryStore) Query(filter QueryFilter) []monitor.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []monitor.Result
	
	// Determine which URLs to check based on filters
	var urlsToCheck []string
	
	if filter.URL != "" {
		// Direct URL filter
		urlsToCheck = []string{filter.URL}
	} else if filter.MonitorName != "" {
		// Filter by monitor name
		if url, exists := s.monitorIndex[filter.MonitorName]; exists {
			urlsToCheck = []string{url}
		}
	} else if filter.Project != "" {
		// Filter by project
		if urls, exists := s.projectIndex[filter.Project]; exists {
			urlsToCheck = urls
		}
	} else {
		// No specific filter, check all URLs
		for url := range s.data {
			urlsToCheck = append(urlsToCheck, url)
		}
	}
	
	// Collect results from the determined URLs
	for _, url := range urlsToCheck {
		if urlResults, exists := s.data[url]; exists {
			for _, result := range urlResults {
				if s.matchesFilter(result, filter) {
					results = append(results, result)
				}
			}
		}
	}
	
	// Sort results by timestamp (newest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].CheckedAt.After(results[j].CheckedAt)
	})
	
	// Apply limit if specified
	if filter.Limit > 0 && len(results) > filter.Limit {
		results = results[:filter.Limit]
	}
	
	return results
}

// matchesFilter checks if a result matches the given filter
func (s *MemoryStore) matchesFilter(result monitor.Result, filter QueryFilter) bool {
	// Project filter
	if filter.Project != "" && !strings.EqualFold(result.Project, filter.Project) {
		return false
	}
	
	// Monitor name filter
	if filter.MonitorName != "" && !strings.EqualFold(result.Name, filter.MonitorName) {
		return false
	}
	
	// URL filter
	if filter.URL != "" && result.URL != filter.URL {
		return false
	}
	
	// Status filter
	if filter.Status != nil && result.Status != *filter.Status {
		return false
	}
	
	// OK filter
	if filter.OK != nil && result.OK != *filter.OK {
		return false
	}
	
	// Time range filters
	if !filter.Since.IsZero() && result.CheckedAt.Before(filter.Since) {
		return false
	}
	
	if !filter.Until.IsZero() && result.CheckedAt.After(filter.Until) {
		return false
	}
	
	return true
}

// GetProjects returns all unique project names
func (s *MemoryStore) GetProjects() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	projects := make([]string, 0, len(s.projectIndex))
	for project := range s.projectIndex {
		projects = append(projects, project)
	}
	
	sort.Strings(projects)
	return projects
}

// GetMonitors returns all unique monitor names
func (s *MemoryStore) GetMonitors() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	monitors := make([]string, 0, len(s.monitorIndex))
	for monitor := range s.monitorIndex {
		monitors = append(monitors, monitor)
	}
	
	sort.Strings(monitors)
	return monitors
}

// GetMonitorsByProject returns all monitor names for a specific project
func (s *MemoryStore) GetMonitorsByProject(project string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var monitors []string
	if urls, exists := s.projectIndex[project]; exists {
		for _, url := range urls {
			if results := s.data[url]; len(results) > 0 {
				// Get the monitor name from the latest result
				monitors = append(monitors, results[len(results)-1].Name)
			}
		}
	}
	
	sort.Strings(monitors)
	return monitors
}