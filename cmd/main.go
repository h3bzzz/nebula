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

	"nebula/handlers"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// Database Connection Config
var db *sqlx.DB

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

func main() {

	// Initialize Database
	db = initDB()
	defer db.Close()

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

	e.POST("/register", handlers.RegisterUser(db))
	e.GET("/register", handlers.RenderRegisterPage)
	e.POST("/login", handlers.LoginUser(db, e.Logger))
	e.GET("/login", handlers.RenderLoginPage)

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
