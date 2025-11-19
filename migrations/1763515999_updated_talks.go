package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2355380017")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("select2744374011")

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
			"hidden": false,
			"id": "select1109482392",
			"maxSelect": 1,
			"name": "transcribee_state",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"todo",
				"created",
				"in_progress",
				"done"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2355380017")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"hidden": false,
			"id": "select2744374011",
			"maxSelect": 1,
			"name": "state",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"new",
				"auto_transcribed",
				"corrected"
			]
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(9, []byte(`{
			"hidden": false,
			"id": "select1109482392",
			"maxSelect": 1,
			"name": "transcribee_state",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"todo",
				"creating",
				"in_progress",
				"done"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
