package handlers

import (
  "net/http"
  "github.com/labstack/echo/v4"
)


func WhoHandler(c echo.Context) error {
  return c.Render(http.StatusOK, "who.html", map[string]interface{}{
    "Title": "whoami",
    "H1": "whoami?",
    "Who": "active",
  })
}
