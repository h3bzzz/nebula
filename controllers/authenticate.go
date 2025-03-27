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
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	ctx    = context.Background()
	userID uuid.UUID
)

// Initialize Redis client
func InitRedis(addr string, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// REGISTER USER
func RegisterUser(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Initialize Redis if needed
		redisClient := c.Get("redis").(*redis.Client)
		if redisClient == nil {
			log.Println("Warning: Redis client not initialized, using default")
			redisClient = InitRedis("localhost:6379", "", 3)
		}

		// Parse multipart form
		if err := c.Request().ParseMultipartForm(10 << 20); err != nil { // 10 MB max
			log.Printf("Form parsing error: %v", err)
			return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
				"Error":         "Failed to parse form data",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Get form values
		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")
		passwordConfirm := c.FormValue("password_confirm")

		// Basic validation
		if username == "" || email == "" || password == "" {
			return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
				"Error":         "All fields are required",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Check if passwords match
		if password != passwordConfirm {
			return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
				"Error":         "Passwords do not match",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Check if password is strong enough (min length 8)
		if len(password) < 8 {
			return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
				"Error":         "Password must be at least 8 characters",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Check if user already exists
		var count int
		query := "SELECT COUNT(*) FROM users WHERE username=$1 OR email=$2"
		log.Printf("Executing query: %s with username=%s, email=%s", query, username, email)
		err := db.Get(&count, query, username, email)
		if err != nil {
			log.Printf("Database error when checking for existing users: %v", err)
			// Check if database is connected
			if err := db.Ping(); err != nil {
				log.Printf("Database connection lost or invalid: %v", err)
				return c.Render(http.StatusInternalServerError, "register.html", map[string]interface{}{
					"Error":         "Database connection error. Please try again later.",
					"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
					"Title":         "Register",
					"ActivePage":    "register",
					"Authenticated": false,
				})
			}
			return c.Render(http.StatusInternalServerError, "register.html", map[string]interface{}{
				"Error":         "Failed to check existing users",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		if count > 0 {
			return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
				"Error":         "Username or email already exists",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Hash Password
		hashedPassword, err := HashPassword(password)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			return c.Render(http.StatusInternalServerError, "register.html", map[string]interface{}{
				"Error":         "Failed to secure password",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Get profile picture (if provided)
		var profilePicture []byte = nil
		file, err := c.FormFile("profile_picture")
		if err == nil && file != nil {
			src, err := file.Open()
			if err == nil {
				defer src.Close()

				// Read the file into memory
				buffer := make([]byte, file.Size)
				_, err = src.Read(buffer)
				if err == nil {
					profilePicture = buffer
				} else {
					log.Printf("Error reading profile picture: %v", err)
				}
			} else {
				log.Printf("Error opening profile picture: %v", err)
			}
		}

		// Insert User into PSQLDB
		userID = uuid.New()
		query = "INSERT INTO users (id, username, email, password, profile_picture, role, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err = db.Exec(query, userID, username, email, hashedPassword, profilePicture, "user", time.Now())

		if err != nil {
			log.Printf("Database insertion error: %v", err)
			return c.Render(http.StatusInternalServerError, "register.html", map[string]interface{}{
				"Error":         "Failed to create user",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Register",
				"ActivePage":    "register",
				"Authenticated": false,
			})
		}

		// Generate Session Token
		sessionToken := GenSessionToken()
		if err := StoreSessionToken(redisClient, sessionToken, userID); err != nil {
			log.Printf("Redis error: %v", err)
			// Continue anyway - user was created
		}

		// Set Session Cookie
		SetSessionCookie(c, sessionToken)

		// Set success message
		SetFlashMessage(c, "success", "Registration successful! Welcome to Sy-Fi Networks.")

		return c.Redirect(http.StatusSeeOther, "/")
	}
}

// LOGIN USER
func LoginUser(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Initialize Redis if needed
		redisClient := c.Get("redis").(*redis.Client)
		if redisClient == nil {
			log.Println("Warning: Redis client not initialized, using default")
			redisClient = InitRedis("localhost:6379", "", 3)
		}

		// Parse form
		if err := c.Request().ParseForm(); err != nil {
			return c.Render(http.StatusBadRequest, "login.html", map[string]interface{}{
				"Error":         "Failed to parse form",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Login",
				"ActivePage":    "login",
				"Authenticated": false,
			})
		}

		// Get form values
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Basic validation
		if username == "" || password == "" {
			return c.Render(http.StatusBadRequest, "login.html", map[string]interface{}{
				"Error":         "Username and password are required",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Login",
				"ActivePage":    "login",
				"Authenticated": false,
			})
		}

		var user struct {
			ID       uuid.UUID `db:"id"`
			Username string    `db:"username"`
			PassHash string    `db:"password"`
		}
		query := "SELECT id, username, password FROM users WHERE username=$1"
		err := db.Get(&user, query, username)
		if err != nil {
			log.Printf("Login query error: %v", err)
			return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
				"Error":         "Invalid username or password",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Login",
				"ActivePage":    "login",
				"Authenticated": false,
			})
		}

		// Verify Password
		if err := VerifyPassword(user.PassHash, password); err != nil {
			return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
				"Error":         "Invalid username or password",
				"CSRFToken":     c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
				"Title":         "Login",
				"ActivePage":    "login",
				"Authenticated": false,
			})
		}

		// Generate Session Token
		sessionToken := GenSessionToken()
		if err := StoreSessionToken(redisClient, sessionToken, user.ID); err != nil {
			log.Printf("Redis error: %v", err)
			// Continue anyway - user can still log in
		}

		// Set Session Cookie
		SetSessionCookie(c, sessionToken)

		// Set success message
		SetFlashMessage(c, "success", "Login successful! Welcome back, "+user.Username+".")

		return c.Redirect(http.StatusSeeOther, "/")
	}
}

// LOGOUT USER
func LogoutUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Initialize Redis if needed
		redisClient := c.Get("redis").(*redis.Client)
		if redisClient == nil {
			log.Println("Warning: Redis client not initialized, using default")
			redisClient = InitRedis("localhost:6379", "", 3)
		}

		// Get the session cookie
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// Just redirect to home if no cookie
			return c.Redirect(http.StatusSeeOther, "/")
		}

		// Delete the session from Redis
		if err := redisClient.Del(ctx, cookie.Value).Err(); err != nil {
			log.Printf("Redis error during logout: %v", err)
		}

		// Expire the cookie
		cookie.Expires = time.Now().Add(-1 * time.Hour)
		cookie.MaxAge = -1
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
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Path = "/"

	// Only use Secure in production
	if c.Request().TLS != nil {
		cookie.Secure = true
	}

	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.MaxAge = 86400 // 24 hours in seconds

	c.SetCookie(cookie)
}

// AuthMiddleware checks if a user is authenticated
func AuthMiddleware(redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the cookie
			cookie, err := c.Cookie("session_token")
			if err != nil {
				// No cookie, user is not authenticated
				// Store authentication status in context
				c.Set("authenticated", false)
				return next(c)
			}

			// Validate token from Redis
			userID, err := EvalSessionToken(redisClient, cookie.Value)
			if err != nil || userID == uuid.Nil {
				// Invalid token or user not found
				c.Set("authenticated", false)
				return next(c)
			}

			// User is authenticated, store user ID in context
			c.Set("authenticated", true)
			c.Set("userID", userID)
			return next(c)
		}
	}
}
