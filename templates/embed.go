package templates

import (
    "embed"
	"flag"
	"html/template"
	"io/fs"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

//go:embed **/*.html
var templateFS embed.FS

var patterns = []string{"**/*.html"}

func getDebugHttpRender(funcMap template.FuncMap, delims render.Delims, files ...string) render.HTMLRender {
	return render.HTMLDebug{Files: files,
		Patterns: patterns, FuncMap: funcMap, Delims: delims}
}

func getProductionHttpRender(funcMap template.FuncMap, delims render.Delims) render.HTMLRender {
	templ := template.Must(template.New("").Delims(delims.Left, delims.Right).Funcs(funcMap).ParseFS(
		templateFS, patterns...))
	return render.HTMLProduction{Template: templ.Funcs(funcMap)}
}

func getFiles() []string {
	var files []string
	var dir string
	if flag.Lookup("test.v") == nil {
		dir = "templates" + string(os.PathSeparator)
	} else {
		dir = ""
	}
	var dirFS = os.DirFS(dir)
	for _, pattern := range patterns {
		match, err := fs.Glob(dirFS, pattern)
		if err != nil {
			continue
		}
		for _, file := range match {
			files = append(files, dir+file)
		}
	}
	return files
}

func GetHTMLRender(funcMap template.FuncMap, delims render.Delims) render.HTMLRender {
	if gin.IsDebugging() {
		return getDebugHttpRender(funcMap, delims, getFiles()...)
	}
	return getProductionHttpRender(funcMap, delims)

}