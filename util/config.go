package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	SiteDomain             string
	SiteSlug               string
	PageTitle              string
	PageDescription        string
	SitemapID              string
	GoogleSiteVerification string
}

var config Config

func setConfig(c Config) {
	config = c
}

// 导出配置字段的 getter
func GetPort() string                   { return config.Port }
func GetSiteDomain() string             { return config.SiteDomain }
func GetSiteSlug() string               { return config.SiteSlug }
func GetPageTitle() string              { return config.PageTitle }
func GetPageDescription() string        { return config.PageDescription }
func GetSitemapID() string              { return config.SitemapID }
func GetGoogleSiteVerification() string { return config.GoogleSiteVerification }

var (
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
	setConfig(Config{
		Port:                   os.Getenv("PORT"),
		SiteDomain:             os.Getenv("SITE_DOMAIN"),
		SiteSlug:               os.Getenv("SITE_SLUG"),
		PageTitle:              os.Getenv("PAGE_TITLE"),
		PageDescription:        os.Getenv("PAGE_DESCRIPTION"),
		SitemapID:              os.Getenv("SITEMAP_ID"),
		GoogleSiteVerification: os.Getenv("GOOGLE_SITE_VERIFICATION"),
	})

	// Log configuration
	LogInfo("Configuration loaded:")
	LogInfo("  Port: %s", GetPort())
	LogInfo("  Site Domain: %s", GetSiteDomain())
	LogInfo("  Site Slug: %s", GetSiteSlug())
	LogInfo("  Page Title: %s", GetPageTitle())
	LogInfo("  Page Description: %s", GetPageDescription())
	LogInfo("  Sitemap ID: %s", GetSitemapID())
	LogInfo("  Google Site Verification: %s", GetGoogleSiteVerification())
}
