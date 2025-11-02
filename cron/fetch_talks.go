package cron

import (
	"log"
	"transcribee-voctoweb/voc_api"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterFetchTalksCron(app *pocketbase.PocketBase, vocApi *voc_api.VocApi) error {
	app.Cron().MustAdd("fetch_talks", "* * * * *", func() {
		conferenceRecords, err := app.FindAllRecords("conferences",
			dbx.HashExp{"autocreate_active": true},
		)
		if err != nil {
			log.Println(err)
			return
		}

		if len(conferenceRecords) == 0 {
			log.Println("No conferences found")
			return
		}

		talksCollection, err := app.FindCollectionByNameOrId("talks")
		if err != nil {
			log.Println(err)
			return
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
				talkRecord.Set("state", "new")
				talkRecord.Set("transcribee_state", "todo")

				err = app.Save(talkRecord)

				if err != nil {
					log.Println(err)
				}
			}
		}
	})

	return nil
}
