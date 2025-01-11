package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	SiteDomain      string
	SiteSlug        string
	PageTitle       string
	PageDescription string
	SitemapID       string
}

var (
	config          Config
	requiredEnvVars = []string{
		"PORT",
		"SITE_DOMAIN",
		"SITE_SLUG",
		"PAGE_TITLE",
		"PAGE_DESCRIPTION",
	}
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Check for required environment variables
	var missingVars []string
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("Error: Required environment variables are missing: %v", missingVars)
	}

	// Load configuration
	config = Config{
		Port:            os.Getenv("PORT"),
		SiteDomain:      os.Getenv("SITE_DOMAIN"),
		SiteSlug:        os.Getenv("SITE_SLUG"),
		PageTitle:       os.Getenv("PAGE_TITLE"),
		PageDescription: os.Getenv("PAGE_DESCRIPTION"),
		SitemapID:       os.Getenv("SITEMAP_ID"),
	}

	// Log configuration
	logInfo("Configuration loaded:")
	logInfo("  Port: %s", config.Port)
	logInfo("  Site Domain: %s", config.SiteDomain)
	logInfo("  Site Slug: %s", config.SiteSlug)
	logInfo("  Page Title: %s", config.PageTitle)
	logInfo("  Page Description: %s", config.PageDescription)
	logInfo("  Sitemap ID: %s", config.SitemapID)
}

func main() {
	logInfo("Server starting at http://0.0.0.0:%s", config.Port)
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil); err != nil {
		logError("Server failed to start: %v", err)
		log.Fatal(err)
	}
}
