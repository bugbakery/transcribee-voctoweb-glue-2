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

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(16, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1079525457",
			"max": 0,
			"min": 0,
			"name": "transcribee_id",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2355380017")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text1079525457")

		return app.Save(collection)
	})
}
