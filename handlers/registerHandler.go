package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Error hashing password: %v", err)
		return ""
	}
	return string(hashedPassword)
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func RegisterUser(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := strings.TrimSpace(c.FormValue("username"))
		email := strings.TrimSpace(c.FormValue("email"))
		password := strings.TrimSpace(c.FormValue("password"))

		if len(username) < 5 || len(password) < 14 || !ValidateEmail(email) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid username, email, or password",
			})
		}

		hashedPassword := HashPassword(password)
		if hashedPassword == "" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		query := `INSERT INTO users (id, username, email, password, created_at)
			VALUES (gen_random_uuid(), $1, $2, $3, $4)`

		_, err := db.ExecContext(context.Background(), query, username, email, hashedPassword, time.Now())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal server error",
			})
		}

		return c.Redirect(http.StatusSeeOther, "/login")

	}
}

// Serve Register Page
func RenderRegisterPage(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"title":    "Register",
		"register": "active",
	})
}
