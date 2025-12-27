package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("userfiles")

		// Leave rules empty, just for superusers

		collection.Fields.Add(
            &core.TextField{
                Name:     "filename",
                Required: true,
                Max:      100,
				Presentable: true,
            },
            &core.FileField{
                Name:        "file",
                Required: true,
            },
			&core.AutodateField{
				Name: "created",
				OnCreate: true,
			})

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("userfiles")
		if err != nil {
			return err
		}

		return app.Save(collection)
	})
}
