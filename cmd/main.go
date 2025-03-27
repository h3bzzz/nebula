package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"nebula/controllers"
	"nebula/handlers"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// Global Variables
var (
	rdb  *redis.Client
	psdb *sqlx.DB
	ctx  = context.Background()
)

// Redis Connection Config
// Need to create a certificate and key for TLS and for securing all paths
// with the server

func rdbInit() {
	// Redis Connection Config
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Error connecting to Redis: %v", err)
		// Continue anyway, as Redis is not critical for basic functionality
	} else {
		log.Println("Connected to Nebula Redis database")
	}
}

// Postgresql Connection Config

func initDB() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
		// Continue anyway, will use environment variables directly
	}

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Error parsing dbUrl: %v", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 5 * time.Minute

	db, err := sqlx.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Connected to Nebula PostgreSQL database")

	return db
}

// Custom Renderer for Echo

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return errors.New("Template not found -> " + name)
	}

	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// Redis middleware to make Redis client available to handlers
func RedisMiddleware(redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("redis", redisClient)
			return next(c)
		}
	}
}

// Database middleware to make database client available to handlers
func DatabaseMiddleware(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

// MAIN
func main() {

	// Initialize Database
	psdb = initDB()
	defer psdb.Close()

	// Initialize Redis
	rdbInit()

	// Initialize Redis Store
	store, err := redistore.NewRediStore(10, "tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		log.Printf("Warning: Error creating Redis Store: %v. Using memory store instead.", err)
		// Continue without Redis store, use cookie store
	} else {
		store.SetMaxAge(86400)
		defer store.Close()
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	cookieStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	cookieStore.Options.SameSite = http.SameSiteLaxMode
	cookieStore.Options.HttpOnly = true
	e.Use(session.Middleware(cookieStore))

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	e.Use(RedisMiddleware(rdb))
	e.Use(DatabaseMiddleware(psdb))
	e.Use(controllers.AuthMiddleware(rdb))
	e.Use(controllers.FlashMiddleware())

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; img-src 'self' data:;",
		ReferrerPolicy:        "same-origin",
	}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:csrf_token",
		CookieSameSite: http.SameSiteLaxMode,
	}))

	templateFiles := []string{
		"home.html",
		"news.html",
		"ttps.html",
		"hof.html",
		"who.html",
		"register.html",
		"login.html",
	}

	templates := make(map[string]*template.Template)

	for _, name := range templateFiles {
		templates[name] = template.Must(
			template.ParseFiles(fmt.Sprintf("templates/%s", name), "templates/base.html"),
		)
	}

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	// Routes
	e.Static("/static", "static")

	e.GET("/", handlers.HomeHandler)
	e.GET("/news", handlers.NewsHandler)
	e.GET("/ttps", handlers.TtpsHandler)
	e.GET("/hof", handlers.HofHandler)
	e.GET("/who", handlers.WhoHandler)

	// S3 Resource Routes
	s3Controller, err := controllers.NewS3ResourcesController()
	if err != nil {
		log.Printf("Warning: S3 controller initialization failed: %v", err)
	} else {
		// Replace the existing news handler with S3-powered articles
		e.GET("/news", s3Controller.ListArticles())
		e.GET("/news/:id", s3Controller.GetArticle())
		e.GET("/images/*", s3Controller.GetImage())

		// Admin routes for article and image management
		adminGroup := e.Group("/admin")
		adminGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// Check if user is authenticated
				if authenticated, ok := c.Get("authenticated").(bool); !ok || !authenticated {
					return c.Redirect(http.StatusSeeOther, "/login")
				}
				return next(c)
			}
		})
		adminGroup.POST("/images/upload", s3Controller.UploadImage())
	}

	// Authentication Routes
	e.POST("/register", controllers.RegisterUser(psdb))
	e.GET("/register", handlers.RenderRegisterPage)
	e.POST("/login", controllers.LoginUser(psdb))
	e.GET("/login", handlers.RenderLoginPage)
	e.GET("/logout", controllers.LogoutUser())

	// End of Routes

	// Start Server
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "7777" // Default port
		}

		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	fmt.Println("Server Nebula is up and running on port", os.Getenv("PORT"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
	fmt.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited gracefully")
}
