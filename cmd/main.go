package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/signal"
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
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       3,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	log.Println("Connected to Nebula Redis database")
}

// Postgresql Connection Config

func initDB() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	return psdb
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

// MAIN
func main() {

	// Initialize Database
	psdb = initDB()
	defer psdb.Close()

	// Initialize Redis
	rdbInit()

	// Initialize Redis Store
	store, err := redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("Error creating Redis Store: %v", err)
	}
	defer store.Close()
	// Session lasts 24 hrs
	store.SetMaxAge(86400)

	// Start Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// Get Secure
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "same-origin",
	}))

	// CSRF Protection
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf_token",
	}))

	// Cookie Store
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

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

	// Authentication Routes
	e.POST("/register", controllers.RegisterUser(psdb))
	e.GET("/register", handlers.RenderRegisterPage)
	e.POST("/login", controllers.LoginUser(psdb))
	e.GET("/login", handlers.RenderLoginPage)
	e.GET("/logout", controllers.LogoutUser())

	// End of Routes

	// Start Server
	go func() {
		if err := e.Start(":7777"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	fmt.Println("Server Nebula is up and running on port 7777")

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
