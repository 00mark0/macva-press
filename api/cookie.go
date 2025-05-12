package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (server *Server) deleteCookie(ctx echo.Context) error {
	ctx.SetCookie(&http.Cookie{
		Name:   "content_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}
