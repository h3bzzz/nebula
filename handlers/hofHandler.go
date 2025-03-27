package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HofHandler(c echo.Context) error {
	authenticated, _ := c.Get("authenticated").(bool)

	data := map[string]interface{}{
		"Title":         "Hacks of Fame",
		"ActivePage":    "hof",
		"Authenticated": authenticated,
	}
	return c.Render(http.StatusOK, "hof.html", data)
}
