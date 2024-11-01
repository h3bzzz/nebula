package handlers

import (
  "net/http"

  "github.com/labstack/echo/v4"
)


func TtpsHandler(c echo.Context) error {
  return c.Render(http.StatusOK, "ttps.html", map[string]interface{}{
    "Title": "Tactics, Techniques, Procedures, & Tools",
    "Ttps": "active",
  })
}
