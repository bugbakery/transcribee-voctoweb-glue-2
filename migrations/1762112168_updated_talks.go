package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2355380017")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"updateRule": "@request.body.title:isset = false &&\n@request.body.subtitle:isset = false &&\n@request.body.state:isset = false &&\n@request.body.description:isset = false &&\n@request.body.duration_secs:isset = false &&\n@request.body.transcribee_url:isset = false &&\n@request.body.transcribee_state:isset = false &&\n@request.body.language:isset = false &&\n@request.body.media_talk_id:isset = false &&\n@request.body.date:isset = false &&\n@request.body.conference:isset = false &&\n@request.body.release_date:isset = false"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2355380017")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"updateRule": "@request.body.title:isset = false"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
