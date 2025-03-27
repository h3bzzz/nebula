package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	authenticated, _ := c.Get("authenticated").(bool)

	data := map[string]interface{}{
		"Title":         "Sy-Fi Networks Home",
		"ActivePage":    "home",
		"Authenticated": authenticated,
	}

	if authenticated {
		if userID, ok := c.Get("userID").(uuid.UUID); ok {
			db := c.Get("db").(*sqlx.DB)
			var username string
			err := db.Get(&username, "SELECT username FROM users WHERE id = $1", userID)
			if err == nil {
				data["Username"] = username
			}
		}
	}

	return c.Render(http.StatusOK, "home.html", data)
}
