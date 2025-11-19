package cron

import (
	"encoding/json"
	"errors"
	"log"
	"transcribee-voctoweb/voc_api"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func fetchTalksCron(app *pocketbase.PocketBase, vocApi *voc_api.VocApi) error {
	conferenceRecords, err := app.FindAllRecords("conferences",
		dbx.HashExp{"autocreate_active": true},
	)
	if err != nil {
		return err
	}

	if len(conferenceRecords) == 0 {
		return errors.New("No conferences found")
	}

	talksCollection, err := app.FindCollectionByNameOrId("talks")
	if err != nil {
		return err
	}

	for _, conference := range conferenceRecords {
		log.Printf("Updating conference '%s'", conference.GetString("name"))
		resp, err := vocApi.GetConference(conference.GetString("name"))

		if err != nil {
			log.Println(err)
			continue
		}

		for _, talkSummary := range resp.Talks {
			count, err := app.CountRecords("talks", dbx.HashExp{"conference": conference.Id, "media_talk_id": talkSummary.Guid})
			if err != nil {
				log.Println(err)
				continue
			}

			if count > 0 {
				continue
			}

			log.Printf("Creating talk '%s'", talkSummary.Title)

			talk, err := vocApi.GetTalk(conference.GetString("name"), talkSummary.Guid)
			if err != nil {
				log.Println(err)
				continue
			}

			persons_json, err := json.Marshal(talk.Persons)
			if err != nil {
				log.Println(err)
				continue
			}

			talkRecord := core.NewRecord(talksCollection)
			talkRecord.Set("conference", conference.Id)
			talkRecord.Set("media_talk_id", talk.Guid)
			talkRecord.Set("title", talk.Title)
			talkRecord.Set("subtitle", talk.Subtitle)
			talkRecord.Set("description", talk.Description)
			talkRecord.Set("duration_secs", talk.Duration)
			talkRecord.Set("language", talk.OriginalLanguage)
			talkRecord.Set("date", talk.Date)
			talkRecord.Set("release_date", talk.ReleaseDate)
			talkRecord.Set("persons", persons_json)
			talkRecord.Set("state", "new")
			talkRecord.Set("transcribee_state", "todo")

			err = app.Save(talkRecord)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func RegisterFetchTalksCron(app *pocketbase.PocketBase, vocApi *voc_api.VocApi) error {
	job_finished := true
	app.Cron().MustAdd("fetch_talks", "* * * * *", func() {
		if !job_finished {
			log.Println("fetch_talks Job already running, skipping")
			return
		}
		job_finished = false
		err := fetchTalksCron(app, vocApi)
		if err != nil {
			log.Println(err)
		}
		job_finished = true
	})

	return nil
}
