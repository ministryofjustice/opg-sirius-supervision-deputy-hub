package main

import (
	"context"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/server"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	logger := telemetry.NewLogger("opg-sirius-supervision-deputy-hub ")
	// manually set time zone
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	envVars, err := server.NewEnvironmentVars()
	if err != nil {
		logger.Info(err.Error())
	}

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
				return envVars.Prefix + s
			},
			"sirius": func(s string) string {
				return envVars.SiriusPublicURL + s
			},
			"firmhub": func(s string) string {
				return envVars.FirmHubURL + s
			},
			"translate":       util.Translate,
			"feature_flagged": util.IsFeatureFlagged(envVars.Features),
			"is_last":         util.IsLast,
		}).
		ParseGlob(envVars.WebDir + "/template/*/*.gotmpl")

	layouts, _ = layouts.Parse(paginate.Template)

	files, _ := filepath.Glob(envVars.WebDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	client, err := sirius.NewClient(http.DefaultClient, envVars.SiriusURL)
	if err != nil {
		logger.Info(err.Error())
	}

	//templates := createTemplates(envVars)

	server := &http.Server{
		Addr:    ":" + envVars.Port,
		Handler: server.New(logger, client, tmpls, envVars),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Info(err.Error())
		}
	}()

	logger.Info("Running at :" + envVars.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info(fmt.Sprint("signal received: ", sig))

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Info(err.Error())
	}
}

//func createTemplates(envVars server.EnvironmentVars) map[string]*template.Template {
//	templates := map[string]*template.Template{}
//	templateFunctions := map[string]interface{}{
//		"contains": func(xs []string, needle string) bool {
//			for _, x := range xs {
//				if x == needle {
//					return true
//				}
//			}
//
//			return false
//		},
//		"title": func(s string) string {
//			r := []rune(s)
//			r[0] = unicode.ToUpper(r[0])
//
//			return string(r)
//		},
//		"prefix": func(s string) string {
//			return envVars.Prefix + s
//		},
//		"sirius": func(s string) string {
//			return envVars.SiriusPublicURL + s
//		},
//	}
//
//	templateDirPath := envVars.WebDir + "/template"
//	templateDir, _ := os.Open(templateDirPath)
//	templateDirs, _ := templateDir.Readdir(0)
//	_ = templateDir.Close()
//
//	mainTemplates, _ := filepath.Glob(templateDirPath + "/*.gotmpl")
//
//	for _, file := range mainTemplates {
//		tmpl := template.New(filepath.Base(file)).Funcs(templateFunctions)
//		for _, dir := range templateDirs {
//			if dir.IsDir() {
//				tmpl, _ = tmpl.ParseGlob(templateDirPath + "/" + dir.Name() + "/*.gotmpl")
//			}
//		}
//		tmpl, _ = tmpl.Parse(paginate.Template)
//		templates[tmpl.Name()] = template.Must(tmpl.ParseFiles(file))
//	}
//
//	return templates
//}
