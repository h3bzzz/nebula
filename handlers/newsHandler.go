package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewsHandler(c echo.Context) error {
	authenticated, _ := c.Get("authenticated").(bool)

	data := map[string]interface{}{
		"Title":         "News",
		"ActivePage":    "news",
		"Authenticated": authenticated,
	}
	return c.Render(http.StatusOK, "news.html", data)
}
