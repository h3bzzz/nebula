package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RenderLoginPage(c echo.Context) error {
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	data := map[string]interface{}{
		"CSRFToken": csrfToken,
		"Title":     "Login",
		"Login":     "active",
	}
	return c.Render(http.StatusOK, "login.html", data)
}
