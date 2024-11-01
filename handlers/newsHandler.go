package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewsHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "news.html", map[string]interface{}{
		"Title": "What Now !?",
		"News":  "active",
	})
}
