package cron

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"transcribee-voctoweb/transcribee_api"
	"transcribee-voctoweb/voc_api"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func createTranscribeeDocumentCron(app *pocketbase.PocketBase, vocApi *voc_api.VocApi, transcribeeApiBaseUrl string) error {
	conferenceRecords, err := app.FindAllRecords("conferences",
		dbx.HashExp{"autocreate_active": true},
	)
	if err != nil {
		return err
	}

	for _, conference := range conferenceRecords {
		talkRecords, err := app.FindRecordsByFilter("talks", "state = 'new' && conference = {:conference}", "date ASC", 0, 0,
			dbx.Params{"conference": conference.Id},
		)
		if err != nil {
			return fmt.Errorf("Error finding talks: %v", err)
		}

		transcribeeToken := conference.GetString("transcribee_user_token")
		transcribeeApi := transcribee_api.New(transcribeeApiBaseUrl, transcribeeToken)

		// get the existing transcribee documents from the transcribee api
		// if we already have a document with a fitting title, we will populate it into our db
		transcribeeDocuments, err := transcribeeApi.GetTranscribeeDocuments()
		if err != nil {
			return fmt.Errorf("Error getting transcribee documents: %v", err)
		}

		for _, talkRecord := range talkRecords {
			log.Printf("Processing talk for adding to transcribee '%s'", talkRecord.GetString("title"))

			for _, document := range transcribeeDocuments {
				if document.Name == talkRecord.GetString("title") {
					log.Printf("Found existing transcribee document for talk '%s'", talkRecord.GetString("title"))
					talkRecord.Set("transcribee_id", document.ID)
					transcribeeUrl, err := transcribeeApi.CreateShareUrl(document.ID)
					if err != nil {
						return fmt.Errorf("Error creating share URL for transcribee document: %v", err)
					}
					talkRecord.Set("transcribee_url", transcribeeUrl)
					talkRecord.Set("state", "auto_transcribed")
					err = app.Save(talkRecord)
					if err != nil {
						return fmt.Errorf("Error updating talk record: %v", err)
					}
					break
				}
			}

			// Check autocreate_limit - if unset or null, treat as unlimited (skip check)
			autocreateLimit := conference.GetInt("autocreate_limit")
			if conference.Get("autocreate_limit") != nil && autocreateLimit <= 0 {
				log.Printf("Conference '%s' has reached autocreate limit (limit: %d), skipping talk '%s'",
					conference.GetString("name"), autocreateLimit, talkRecord.GetString("title"))
				continue
			}

			// Get detailed talk information from VOC API
			talk, err := vocApi.GetTalk(conference.GetString("name"), talkRecord.GetString("media_talk_id"))
			if err != nil {
				log.Printf("Error getting talk details for %s: %v", talkRecord.Id, err)
				continue
			}

			// Find the best video recording
			var selectedRecording *voc_api.Recording
			for i := range talk.Recordings {
				recording := &talk.Recordings[i]
				if strings.HasPrefix(recording.MimeType, "video/") {
					if selectedRecording == nil || !recording.HighQuality {
						selectedRecording = recording
					}
				}
			}
			if selectedRecording == nil {
				log.Printf("No video recording found for talk %s", talkRecord.GetString("title"))
				continue
			}
			log.Printf("Selected video recording: %s (type: %s, quality: %v) for talk %s",
				selectedRecording.Filename, selectedRecording.MimeType, selectedRecording.HighQuality, talkRecord.GetString("title"))

			// Download the media file
			resp, err := http.Get(selectedRecording.RecordingUrl)
			if err != nil {
				log.Printf("Error downloading recording for talk %s: %v", talkRecord.GetString("title"), err)
				continue
			}

			if resp.StatusCode != 200 {
				log.Printf("Error downloading recording for talk %s: HTTP %d", talkRecord.GetString("title"), resp.StatusCode)
				resp.Body.Close()
				continue
			}

			log.Printf("Downloaded video recording %s (HTTP %d) for talk %s", selectedRecording.RecordingUrl, resp.StatusCode, talkRecord.GetString("title"))

			var language string
			switch talk.OriginalLanguage {
			case "deu":
				language = "de"
			case "eng":
				language = "en"
			default:
				language = "auto"
			}

			documentBody := &transcribee_api.DocumentBodyWithFile{
				Language: language,
				Model:    "large",
				Name:     talk.Title,
				FileName: selectedRecording.Filename,
				File:     resp.Body,
			}

			log.Printf("Creating transcribee document for talk: %s", talkRecord.GetString("title"))
			document, err := transcribeeApi.CreateDocument(documentBody)
			resp.Body.Close()

			if err != nil {
				log.Printf("Error creating transcribee document for talk %s: %v", talkRecord.GetString("title"), err)
				continue
			}

			// Update talk record with transcribee information
			talkRecord.Set("transcribee_id", document.ID)
			transcribeeUrl, err := transcribeeApi.CreateShareUrl(document.ID)
			if err != nil {
				return fmt.Errorf("Error creating share URL for transcribee document: %v", err)
			}

			talkRecord.Set("transcribee_url", transcribeeUrl)
			talkRecord.Set("state", "auto_transcribed")

			err = app.Save(talkRecord)
			if err != nil {
				log.Printf("Error saving transcribee information for talk %s: %v", talkRecord.Id, err)
				continue
			}

			log.Printf("Successfully created transcribee document for talk '%s' (ID: %s, URL: %s)",
				talkRecord.GetString("title"), document.ID, transcribeeUrl)

			// Decrement the autocreate_limit for the conference
			conference.Set("autocreate_limit", autocreateLimit-1)
			err = app.Save(conference)
			if err != nil {
				log.Printf("Error updating autocreate_limit for conference %s: %v", conference.GetString("name"), err)
				// Don't continue here as the talk was successfully created, just log the error
			} else {
				log.Printf("Decremented autocreate_limit for conference '%s' to %d", conference.GetString("name"), autocreateLimit-1)
			}
		}
	}
	return nil
}

func RegisterCreateTranscribeeDocumentsCron(app *pocketbase.PocketBase, vocApi *voc_api.VocApi, transcribeeApiBaseUrl string) error {
	job_finished := true
	app.Cron().MustAdd("create_transcribee_documents", "* * * * *", func() {
		if !job_finished {
			log.Println("create_transcribee_documents Job already running, skipping")
			return
		}
		job_finished = false
		err := createTranscribeeDocumentCron(app, vocApi, transcribeeApiBaseUrl)
		if err != nil {
			log.Println(err)
		}
		job_finished = true
	})

	return nil
}
