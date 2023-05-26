package src

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	// you can limit concurrent net request. It's optional
	MaxGoroutines = 1
	// timeout for net requests
	Timeout = 2 * time.Second
)

type SiteStatus struct {
	Name          string
	StatusCode    int
	TimeOfRequest time.Duration
}

type Monitor struct {
	StatusMap        map[string]SiteStatus
	Mtx              *sync.Mutex
	G                errgroup.Group
	Sites            []string
	RequestFrequency time.Duration
}

func NewMonitor(sites []string, requestFrequency time.Duration) *Monitor {
	return &Monitor{
		StatusMap:        make(map[string]SiteStatus),
		Mtx:              &sync.Mutex{},
		Sites:            sites,
		RequestFrequency: requestFrequency,
	}
}

// Run printStatuses and checkSite in different goroutines
func (m *Monitor) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	for _, website := range m.Sites {
		website := website
		wg.Add(1)
		go func() {
			m.checkSite(ctx, website)
			defer wg.Done()
		}()
	}
	wg.Wait()
	m.printStatuses()
	return nil
}

// Check web-site and write result to StatusMap
func (m *Monitor) checkSite(ctx context.Context, site string) {
	client := http.Client{
		Timeout: Timeout, // Timeout specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body.
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, site, nil)
	if err != nil {
		log.Printf("Error while creating request. Error: %s", err)
	}
	start := time.Now()
	response, err := client.Do(request)
	now := time.Now()
	if err != nil {
		log.Printf("Error while receiving data from the website: %s error: %s\n", site, err)
		m.Mtx.Lock()
		m.StatusMap[site] = SiteStatus{Name: site, StatusCode: http.StatusGatewayTimeout, TimeOfRequest: now.Sub(start)}
		m.Mtx.Unlock()
		return
	}
	defer response.Body.Close()
	m.Mtx.Lock()
	m.StatusMap[site] = SiteStatus{Name: site, StatusCode: response.StatusCode, TimeOfRequest: now.Sub(start)}
	m.Mtx.Unlock()
}

// Iterate over the map and print results
func (m *Monitor) printStatuses() {
	for _, status := range m.StatusMap {
		log.Println(status)
	}
}
