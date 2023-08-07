package server

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type EnvironmentVars struct {
	Port            string
	WebDir          string
	SiriusURL       string
	SiriusPublicURL string
	FirmHubURL      string
	Prefix          string
	DefaultPaTeam   int
	DefaultProTeam  int
	Features        []string
}

func NewEnvironmentVars() (EnvironmentVars, error) {
	defaultPaTeamId, err := strconv.Atoi(getEnv("DEFAULT_PA_TEAM", "23"))
	if err != nil {
		return EnvironmentVars{}, errors.New("error converting DEFAULT_PA_TEAM to int")
	}

	defaultProTeamId, err := strconv.Atoi(getEnv("DEFAULT_PRO_TEAM", "28"))
	if err != nil {
		return EnvironmentVars{}, errors.New("error converting DEFAULT_PRO_TEAM to int")
	}

	return EnvironmentVars{
		Port:            getEnv("PORT", "1234"),
		WebDir:          getEnv("WEB_DIR", "web"),
		SiriusURL:       getEnv("SIRIUS_URL", "http://localhost:8080"),
		SiriusPublicURL: getEnv("SIRIUS_PUBLIC_URL", ""),
		FirmHubURL:      getEnv("FIRM_HUB_HOST", "") + "/supervision/deputies/firm",
		Prefix:          getEnv("PREFIX", ""),
		DefaultPaTeam:   defaultPaTeamId,
		DefaultProTeam:  defaultProTeamId,
		Features:        strings.Split(getEnv("FEATURES", ""), ","),
	}, nil
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
