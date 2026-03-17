package config

import "os"

// ResolvedConfig holds the final resolved configuration after layering
// flags > env vars > profile values.
type ResolvedConfig struct {
	URL      string
	Token    string
	Username string
	Password string
	OrgID    int64
	Output   string
}

// Resolve merges flag values, environment variables, and profile values.
// Priority: flags > env vars > profile values.
func Resolve(flagURL, flagToken, flagUsername, flagPassword string, flagOrgID int64, profile *Profile, defaults Defaults) *ResolvedConfig {
	rc := &ResolvedConfig{}

	// URL
	rc.URL = firstNonEmpty(flagURL, os.Getenv("GRAFANA_URL"))
	if rc.URL == "" && profile != nil {
		rc.URL = profile.URL
	}

	// Token
	rc.Token = firstNonEmpty(flagToken, os.Getenv("GRAFANA_TOKEN"))
	if rc.Token == "" && profile != nil {
		rc.Token = profile.Token
	}

	// Username
	rc.Username = firstNonEmpty(flagUsername, os.Getenv("GRAFANA_USERNAME"))
	if rc.Username == "" && profile != nil {
		rc.Username = profile.Username
	}

	// Password
	rc.Password = firstNonEmpty(flagPassword, os.Getenv("GRAFANA_PASSWORD"))
	if rc.Password == "" && profile != nil {
		rc.Password = profile.Password
	}

	// OrgID
	if flagOrgID > 0 {
		rc.OrgID = flagOrgID
	} else if envOrgID := os.Getenv("GRAFANA_ORG_ID"); envOrgID != "" {
		var id int64
		for _, c := range envOrgID {
			if c >= '0' && c <= '9' {
				id = id*10 + int64(c-'0')
			}
		}
		if id > 0 {
			rc.OrgID = id
		}
	} else if profile != nil {
		rc.OrgID = profile.OrgID
	}

	// Output
	rc.Output = defaults.Output
	if rc.Output == "" {
		rc.Output = "table"
	}

	return rc
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
