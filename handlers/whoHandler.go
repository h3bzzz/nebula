package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func WhoHandler(c echo.Context) error {
	authenticated, _ := c.Get("authenticated").(bool)

	data := map[string]interface{}{
		"Title":         "Whoami",
		"ActivePage":    "who",
		"Authenticated": authenticated,
	}
	return c.Render(http.StatusOK, "who.html", data)
}
