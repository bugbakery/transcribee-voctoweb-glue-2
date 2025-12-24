package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"transcribee-voctoweb/cron"
	"transcribee-voctoweb/handlers"
	"transcribee-voctoweb/transcribee_api"

	// "transcribee-voctoweb/transcribee_api"
	"transcribee-voctoweb/voc_api"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
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
	vocApi := voc_api.New(vocApiBaseUrl, os.Getenv("VOC_API_TOKEN"))

	transcribeeApiBaseUrl := os.Getenv("TRANSCRIBEE_API_BASEURL")
	if transcribeeApiBaseUrl == "" {
		transcribeeApiBaseUrl = "https://beta.transcribee.net"
	}

	cron.RegisterFetchTalksCron(app, vocApi)
	cron.RegisterCreateTranscribeeDocumentsCron(app, vocApi, transcribeeApiBaseUrl)
	cron.RegisterMonitorTranscriptionProgressCron(app, transcribeeApiBaseUrl)
	registerAssigneeHistory(app)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		err, customIndexHtml := buildCustomIndexHtml()
		if err != nil {
			return err
		}

		se.Router.POST("/api/talks/{id}/publish", func(e *core.RequestEvent) error {
			id := e.Request.PathValue("id")

			talkRecord, err := e.App.FindRecordById("talks", id)
			if err != nil {
				return err
			}

			errs := app.ExpandRecord(talkRecord, []string{"conference"}, nil)
			if len(errs) > 0 {
			    return fmt.Errorf("failed to expand: %v", errs)
			}

			conference := talkRecord.ExpandedOne("conference")

			transcribeeToken := conference.GetString("transcribee_user_token")
			transcribeeApi := transcribee_api.New(transcribeeApiBaseUrl, transcribeeToken)

			vtt, err := transcribeeApi.Export(talkRecord.GetString("transcribee_id"), "VTT", true, false, 60)
			if err != nil {
				return err
			}

			err = vocApi.UploadVtt(talkRecord.GetString("media_talk_id"), conference.GetString("name"), []byte(vtt), talkRecord.GetString("language"))
			if err != nil {
				return err
			}

			// transcribee_api.Export
        	return e.JSON(http.StatusOK, map[string]any{"success": vtt})
    	}).Bind(apis.RequireAuth())

		// serves static files from the provided public dir and falls back to custom index.html")
		se.Router.GET("/{path...}", handlers.StaticWithCustomIndexHtml(os.DirFS("./pb_public"), customIndexHtml))

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
