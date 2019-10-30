package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

// ctxKey represents the type of value for the context key.
type ctxKey string

// KeyValues is how request values or stored/retrieved.
const KeyValues = ctxKey("hejiaji")

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

type App struct {
	*httptreemux.TreeMux
	shutdown chan os.Signal
	log      *log.Logger
	mv       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, log *log.Logger, mv ...Middleware) *App {
	app := App{
		TreeMux:  httptreemux.New(),
		shutdown: shutdown,
		log:      log,
		mv:       mv,
	}

	return &app
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.log.Println("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGSTOP
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mv ...Middleware) {

	handler = wrapMiddleware(mv, handler)

	handler = wrapMiddleware(a.mv, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := handler(ctx, w, r, params); err != nil {
			a.log.Printf("critical error: %s", err)
			a.SignalShutdown()
			return
		}
	}

	// Add this handler for the specified verb and route.
	a.TreeMux.Handle(verb, path, h)
}
