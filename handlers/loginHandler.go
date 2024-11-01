package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func validateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func LoginUser(db *sqlx.DB, logger echo.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))
		password := strings.TrimSpace(c.FormValue("password"))

		if !validateEmail(email) || len(password) < 14 {
			logger.Warn("Invalid email or password")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email or password"})
		}

		time.Sleep(2 * time.Second)

		var hashedPassword string
		query := `SELECT password_hash FROM users WHERE email = $1`
		err := db.GetContext(context.Background(), &hashedPassword, query, email)
		if err != nil {
			logger.Warn("Login Failed: email not found")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login Failed: email not found"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			logger.Warn("Login Failed: password incorrect")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login Failed: password incorrect"})
		}

		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		sess.Values["authenticated"] = true
		sess.Values["email"] = email
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			logger.Error("Error saving session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error saving session"})
		}

		logger.Info("Log in successful, Welcome!", email)
		return c.JSON(http.StatusOK, map[string]string{"message": "Log in successful, Welcome!"})
	}
}

func LogoutUser(logger echo.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())

		logger.Info("Logged out")
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}

// Render Login Page
func RenderLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"title": "Login",
		"login": "active",
	})
}
