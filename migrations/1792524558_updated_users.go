package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"oauth2": {
				"mappedFields": {
					"avatarURL": "",
					"name": "",
					"username": ""
				}
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"oauth2": {
				"mappedFields": {
					"avatarURL": "avatar",
					"name": "name",
					"username": "username"
				}
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
