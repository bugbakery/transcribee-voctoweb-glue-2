package main

import (
	"io/fs"
	"log"
	"os"
	"strings"
	"transcribee-voctoweb/cron"
	"transcribee-voctoweb/handlers"
	"transcribee-voctoweb/voc_api"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "transcribee-voctoweb/migrations"
)

func registerAssigneeHistory(app *pocketbase.PocketBase) {
	trackAssigneeHistory := func(re *core.RecordRequestEvent) error {
		err := re.Next()
		if err != nil {
			return err
		}

		if re.Record.Get("assignee") != re.Record.Original().Get("assignee") {
			hist_collection, err := re.App.FindCollectionByNameOrId("talks_assignee_history")
			if err != nil {
				return err
			}

			hist_record := core.NewRecord(hist_collection)

			hist_record.Set("talk", re.Record.Id)
			hist_record.Set("assignee", re.Record.Get("assignee"))

			if re.Auth.Collection().Name == "users" {
				hist_record.Set("changed_by", re.Auth.Id)
			}

			err = re.App.Save(hist_record)
			if err != nil {
				return err
			}
		}

		return nil
	}

	app.OnRecordCreateRequest("talks").BindFunc(trackAssigneeHistory)
	app.OnRecordUpdateRequest("talks").BindFunc(trackAssigneeHistory)
}

func buildCustomIndexHtml() (error, string) {
	content, err := fs.ReadFile(os.DirFS("./pb_public"), "index.html")
	if err != nil {
		return err, ""
	}
	originalIndex := string(content)

	return nil, strings.Replace(originalIndex, "<base href=\"/\">", "<base href=\"/\">", 1)
}

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	vocApiBaseUrl := os.Getenv("VOC_API_BASEURL")
	if vocApiBaseUrl == "" {
		vocApiBaseUrl = "https://publishing.c3voc.de/api"
	}
	vocApi := voc_api.New(vocApiBaseUrl)

	cron.RegisterFetchTalksCron(app, vocApi)
	registerAssigneeHistory(app)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		err, customIndexHtml := buildCustomIndexHtml()
		if err != nil {
			return err
		}

		// serves static files from the provided public dir and falls back to custom index.html")
		se.Router.GET("/{path...}", handlers.StaticWithCustomIndexHtml(os.DirFS("./pb_public"), customIndexHtml))

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
