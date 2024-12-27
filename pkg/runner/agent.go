package runner

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/projectdiscovery/gologger"
)

const (
	customSearchURL = "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&start=%d&fields=items(link)"
	maxResults      = 100 // Google CSE max results
	resultsPerPage  = 10  // Google CSE returns 10 results per page
)

type SearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

type Agent struct {
	Client         *http.Client
	APIKey         []string
	SearchEngineID []string
	currentKeyIdx  int
}

func NewAgent(apiKey []string, searchId []string, proxy string) *Agent {
	Transport := &http.Transport{}
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		if proxyURL == nil {
			gologger.Warning().Msgf("Invalid proxy URL: %s", proxy)
		} else {
			Transport.Proxy = http.ProxyURL(proxyURL)
		}
	}
	return &Agent{
		Client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: Transport,
		},
		APIKey:         apiKey,
		SearchEngineID: searchId,
		currentKeyIdx:  0,
	}
}

func (a *Agent) rotateAPIKey() bool {
	if a.currentKeyIdx+1 < len(a.APIKey) {
		a.currentKeyIdx++
		return true
	}
	return false
}

func (a *Agent) Dork(query string) ([]string, error) {
	var allURLs []string
	startIndex := 1

	for startIndex <= maxResults {
		encodedQuery := url.QueryEscape(query)
		endPoint := fmt.Sprintf(customSearchURL,
			a.APIKey[a.currentKeyIdx],
			a.SearchEngineID[a.currentKeyIdx],
			encodedQuery,
			startIndex)

		req, err := http.NewRequest(http.MethodGet, endPoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("User-Agent", "pagode (gzip)")

		resp, err := a.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		// Handle API quota errors
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == 429 {
			resp.Body.Close()
			if !a.rotateAPIKey() {
				return allURLs, fmt.Errorf("all API keys exhausted")
			}
			gologger.Warning().Msgf("API key limit reached, rotating to next key")
			continue // Retry with new key
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var searchResp SearchResponse
		gz, _ := gzip.NewReader(resp.Body)
		if err := json.NewDecoder(gz).Decode(&searchResp); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		// If no results returned, we've reached the end
		if len(searchResp.Items) == 0 {
			break
		}

		for _, item := range searchResp.Items {
			allURLs = append(allURLs, item.Link)
		}

		startIndex += resultsPerPage
	}

	return allURLs, nil
}
