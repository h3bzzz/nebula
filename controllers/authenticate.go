package controllers

import (
	"context"
	"crypto"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
	psdb        *sqlx.DB
	userID      uuid.UUID
)

// REGISTER USER
func RegisterUser(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Username       string `json:"username" validate:"required"`
			Email          string `json:"email" validate:"required,email"`
			Password       string `json:"password" validate:"required,min=14"`
			ProfilePicture []byte `json:"profile_picture" validate:"required"`
		}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		// Check if user already exists
		var count int
		query := "SELECT COUNT(*) FROM users WHERE username=$1 OR email=$2"
		err := db.Get(&count, query, req.Username, req.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to query user",
			})
		}

		if count > 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "User already exists",
			})
		}

		// Hash Password
		HashPassword, err := HashPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to hash password",
			})
		}

		// Insert User into PSQLDB
		userID = uuid.New()
		query = "INSERT INTO users (id, username, email, password_hash, profile_picture, role, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err = db.Exec(query, userID, req.Username, req.Email, HashPassword, req.ProfilePicture, "user", time.Now())

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to create user",
			})
		}

		// Generate Session Token
		sessionToken := GenSessionToken()
		if err := StoreSessionToken(redisClient, sessionToken, userID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to store session token",
			})
		}

		// Set Session Cookie, have redis manage session tokens
		SetSessionCookie(c, sessionToken)
		StoreSessionToken(redisClient, sessionToken, userID)

		return c.Redirect(http.StatusSeeOther, "/")

		// Give User "user" role

		// Store in S3 Bucket Profile Picture

	}
}

// LOGIN USER
func LoginUser(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		var user struct {
			ID       uuid.UUID `db:"id"`
			Username string    `db:"username"`
			PassHash string    `db:"password_hash"`
		}
		query := "SELECT id, username, password_hash FROM users WHERE username=$1"
		err := db.Get(&user, query, req.Username)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid username",
			})
		}

		// Verify Password
		if err := VerifyPassword(user.PassHash, req.Password); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid password",
			})
		}

		// Generate Session Token
		sessionToken := GenSessionToken()
		if err := StoreSessionToken(redisClient, sessionToken, user.ID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to store session token",
			})
		}

		// Set Session Cookie
		SetSessionCookie(c, sessionToken)
		StoreSessionToken(redisClient, sessionToken, user.ID)

		return c.Redirect(http.StatusSeeOther, "/")
	}
}

// LOGOUT USER
func LogoutUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(userID.String())
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Failed to get cookie",
			})
		}

		if err := redisClient.Del(ctx, cookie.Value).Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to delete session token",
			})
		}

		cookie.Expires = time.Now()
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/")
	}
}

// Generate Token with salting and hashing for session management
func GenSessionToken() string {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	h := crypto.SHA256.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Eliminate verbosity on errors for security
func StoreSessionToken(redisClient *redis.Client, token string, userID uuid.UUID) error {
	err := redisClient.Set(ctx, token, userID, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Check Session Token
func EvalSessionToken(redisClient *redis.Client, token string) (uuid.UUID, error) {
	// Check if token exists
	userID, err := redisClient.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	// Parse UUID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return parsedUserID, nil
}

// Hash Password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// Verify Password
func VerifyPassword(HashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(password))
}

// Set Session Cookie
func SetSessionCookie(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = userID.String()
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(24 * time.Hour)

	c.SetCookie(cookie)
}
