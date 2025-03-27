package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TtpsHandler(c echo.Context) error {
	authenticated, _ := c.Get("authenticated").(bool)

	data := map[string]interface{}{
		"Title":         "TTPS",
		"ActivePage":    "ttps",
		"Authenticated": authenticated,
	}
	return c.Render(http.StatusOK, "ttps.html", data)
}
