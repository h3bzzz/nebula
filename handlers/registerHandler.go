package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RenderRegisterPage(c echo.Context) error {
	if authenticated, ok := c.Get("authenticated").(bool); ok && authenticated {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	data := map[string]interface{}{
		"CSRFToken":     csrfToken,
		"Title":         "Register",
		"ActivePage":    "register",
		"Authenticated": false,
	}
	return c.Render(http.StatusOK, "register.html", data)
}
