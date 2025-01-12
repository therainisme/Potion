package proxy

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/therainisme/potion/util"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

func loadDatabasePages(r *http.Request) ([]URL, error) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	// Build request payload
	// TODO: Handle pagination
	payload := map[string]interface{}{
		"page": map[string]interface{}{
			"id": util.GetSitemapID(),
		},
		"limit":           30,
		"cursor":          map[string]interface{}{"stack": []interface{}{}},
		"verticalColumns": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v3/loadCachedPageChunkV2", util.GetSiteDomain()),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Get view ID from the response
	var viewID string
	if recordMap, ok := result["recordMap"].(map[string]interface{}); ok {
		if block, ok := recordMap["block"].(map[string]interface{}); ok {
			if blockValue, ok := block[util.GetSitemapID()].(map[string]interface{}); ok {
				if value, ok := blockValue["value"].(map[string]interface{}); ok {
					if viewIds, ok := value["view_ids"].([]interface{}); ok && len(viewIds) > 0 {
						viewID = viewIds[0].(string)
					}
				}
			}
		}
	}

	var urls []URL
	if recordMap, ok := result["recordMap"].(map[string]interface{}); ok {
		if collectionView, ok := recordMap["collection_view"].(map[string]interface{}); ok {
			if view, ok := collectionView[viewID].(map[string]interface{}); ok {
				if value, ok := view["value"].(map[string]interface{}); ok {
					if pageSort, ok := value["page_sort"].([]interface{}); ok {
						// ???? Skip the first ID and process the rest
						// TODO: i := 1; i < len(pageSort); i++
						for i := 0; i < len(pageSort); i++ {
							if pageId, ok := pageSort[i].(string); ok {
								urls = append(urls, URL{
									Loc:        fmt.Sprintf("%s/%s", baseURL, strings.ReplaceAll(pageId, "-", "")),
									LastMod:    time.Now().Format("2006-01-02"),
									ChangeFreq: "daily",
									Priority:   0.8,
								})
							}
						}
					}
				}
			}
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("no pages found in the collection view")
	}

	return urls, nil
}

func handleSitemap(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	// Create base URL set
	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []URL{
			{
				Loc:        fmt.Sprintf("%s/%s", baseURL, util.GetSiteSlug()),
				LastMod:    time.Now().Format("2006-01-02"),
				ChangeFreq: "daily",
				Priority:   1.0,
			},
		},
	}

	// Load database pages
	if dbUrls, err := loadDatabasePages(r); err == nil {
		urlset.URLs = append(urlset.URLs, dbUrls...)
	} else {
		util.LogError("Failed to load database pages: %v", err)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlset); err != nil {
		util.LogError("Failed to encode sitemap: %v", err)
		http.Error(w, "Failed to generate sitemap", http.StatusInternalServerError)
		return
	}
}
