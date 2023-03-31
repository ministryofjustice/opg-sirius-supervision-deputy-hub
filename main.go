package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/server"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
)

func main() {
	logger := logging.New(os.Stdout, "opg-sirius-supervision-deputy-hub ")

	port := getEnv("PORT", "1234")
	webDir := getEnv("WEB_DIR", "web")
	siriusURL := getEnv("SIRIUS_URL", "http://localhost:8080")
	siriusPublicURL := getEnv("SIRIUS_PUBLIC_URL", "")
	prefix := getEnv("PREFIX", "")
	DefaultPaTeam := getEnv("DEFAULT_PA_TEAM", "23")
	DefaultProTeam := getEnv("DEFAULT_PRO_TEAM", "28")
	firmHubURL := getEnv("FIRM_HUB_HOST", "") + "/supervision/deputies/firm"
	features := strings.Split(getEnv("FEATURES", ""), ",")

	layouts, _ := template.
		New("").
		Funcs(map[string]interface{}{
			"join": func(sep string, items []string) string {
				return strings.Join(items, sep)
			},
			"contains": func(xs []string, needle string) bool {
				for _, x := range xs {
					if x == needle {
						return true
					}
				}

				return false
			},
			"prefix": func(s string) string {
				return prefix + s
			},
			"sirius": func(s string) string {
				return siriusPublicURL + s
			},
			"firmhub": func(s string) string {
				return firmHubURL + s
			},
			"translate":       util.Translate,
			"rename_errors":   util.RenameErrors,
			"feature_flagged": util.IsFeatureFlagged(features),
			"is_last":         util.IsLast,
		}).
		ParseGlob(webDir + "/template/*/*.gotmpl")

	files, _ := filepath.Glob(webDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	client, err := sirius.NewClient(http.DefaultClient, siriusURL)
	if err != nil {
		logger.Fatal(err)
	}

	defaultPATeam, err := strconv.Atoi(DefaultPaTeam)
	if err != nil {
		logger.Print("Error converting DEFAULT_PA_TEAM to int")
	}

	defaultPROTeam, err := strconv.Atoi(DefaultProTeam)
	if err != nil {
		logger.Print("Error converting DEFAULT_PRO_TEAM to int")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: server.New(logger, client, tmpls, prefix, siriusPublicURL, webDir, defaultPATeam, defaultPROTeam),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Print("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Print("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Print(err)
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
