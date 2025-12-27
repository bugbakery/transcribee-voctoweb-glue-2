package hooks

import (
	"github.com/pocketbase/pocketbase/core"
)

func BindAppHooks(app core.App) {

	app.OnRecordAfterCreateSuccess("userfiles").BindFunc(CreateUsers)
}