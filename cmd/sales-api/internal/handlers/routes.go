package handlers

import (
	"github.com/hejiaji/go-conf/internal/mid"
	"github.com/hejiaji/go-conf/internal/platform/web"
	"log"
	"net/http"
	"os"

)

func API(shutdown chan os.Signal, log *log.Logger) http.Handler {
	app := web.NewApp(shutdown, log, mid.Logger(log), mid.Errors(log), mid.Metrics())
	app.Handle("GET", "/health", health)

	return app
}
