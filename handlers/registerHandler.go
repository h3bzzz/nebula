package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RenderRegisterPage(c echo.Context) error {
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	data := map[string]interface{}{
		"CSRFToken": csrfToken,
		"Title":     "Register",
		"Register":  "active",
	}
	return c.Render(http.StatusOK, "register.html", data)
}
