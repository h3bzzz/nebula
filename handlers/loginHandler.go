package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RenderLoginPage(c echo.Context) error {
	if authenticated, ok := c.Get("authenticated").(bool); ok && authenticated {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	data := map[string]interface{}{
		"CSRFToken":     csrfToken,
		"Title":         "Login",
		"ActivePage":    "login",
		"Authenticated": false,
	}
	return c.Render(http.StatusOK, "login.html", data)
}
