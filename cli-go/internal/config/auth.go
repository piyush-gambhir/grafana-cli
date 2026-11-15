package config

import (
	"os"
	"strconv"
	"strings"
)

type ResolvedConfig struct {
	URL      string
	Token    string
	Username string
	Password string
	OrgID    int64
	Output   string
	ReadOnly bool
}

func Resolve(flagURL, flagToken, flagUsername, flagPassword string, flagOrgID int64, profile *Profile, defaults Defaults) *ResolvedConfig {
	rc := &ResolvedConfig{}

	rc.URL = firstNonEmpty(flagURL, os.Getenv("GRAFANA_URL"))
	if rc.URL == "" && profile != nil {
		rc.URL = profile.URL
	}

	rc.Token = firstNonEmpty(flagToken, os.Getenv("GRAFANA_TOKEN"))
	if rc.Token == "" && profile != nil {
		rc.Token = profile.Token
	}

	rc.Username = firstNonEmpty(flagUsername, os.Getenv("GRAFANA_USERNAME"))
	if rc.Username == "" && profile != nil {
		rc.Username = profile.Username
	}

	rc.Password = firstNonEmpty(flagPassword, os.Getenv("GRAFANA_PASSWORD"))
	if rc.Password == "" && profile != nil {
		rc.Password = profile.Password
	}

	if flagOrgID > 0 {
		rc.OrgID = flagOrgID
	} else if v := os.Getenv("GRAFANA_ORG_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			rc.OrgID = id
		}
	} else if profile != nil {
		rc.OrgID = profile.OrgID
	}

	rc.Output = defaults.Output
	if rc.Output == "" {
		rc.Output = "table"
	}

	if profile != nil {
		rc.ReadOnly = profile.ReadOnly
	}
	if envRO := os.Getenv("GRAFANA_READ_ONLY"); envRO != "" {
		rc.ReadOnly = strings.EqualFold(envRO, "true") || envRO == "1"
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
