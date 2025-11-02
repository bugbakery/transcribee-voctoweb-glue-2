package handlers

import (
	"io/fs"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func StaticWithCustomIndexHtml(fsys fs.FS, customIndexHtml string) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		path := e.Request.PathValue(apis.StaticWildcardParam)

		if path == "" || path == "index.html" {
			// override index.html
			return e.HTML(200, customIndexHtml)
		}

		staticResult := apis.Static(fsys, false)(e)
		if staticResult == router.ErrFileNotFound {
			// fallback to custom index.html
			return e.HTML(200, customIndexHtml)
		}

		return staticResult
   }
}
