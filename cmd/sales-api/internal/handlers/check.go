package handlers

import (
	"context"
	"github.com/hejiaji/go-conf/internal/platform/web"
	"net/http"
)

type User struct {
	Name string
	Email string
	Alias string
}

func (u *User) sendEmail() {

}

func health(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
