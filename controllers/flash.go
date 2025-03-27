package controllers

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SetFlashMessage(c echo.Context, key string, message string) error {
	sess, err := session.Get("flash", c)
	if err != nil {
		return err
	}
	sess.Values[key] = message
	sess.Save(c.Request(), c.Response())
	return nil
}

func GetFlashMessage(c echo.Context, key string) string {
	sess, err := session.Get("flash", c)
	if err != nil {
		return ""
	}

	message, ok := sess.Values[key].(string)
	if !ok {
		return ""
	}

	delete(sess.Values, key)
	sess.Save(c.Request(), c.Response())

	return message
}

func FlashMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			successMessage := GetFlashMessage(c, "success")
			errorMessage := GetFlashMessage(c, "error")

			if successMessage != "" {
				c.Set("SuccessMessage", successMessage)
			}
			if errorMessage != "" {
				c.Set("ErrorMessage", errorMessage)
			}

			return next(c)
		}
	}
}
