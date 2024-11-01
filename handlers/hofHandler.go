package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HofHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "hof.html", map[string]interface{}{
		"Title": "Hacks of Fame",
		"Hof":   "active",
	})
}
