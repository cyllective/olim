package app

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/rtfmkiesel/kisslog"
	"golang.org/x/time/rate"

	"github.com/cyllective/onetim3/internal/db"
)

var (
	//go:embed static
	embedFSStatic embed.FS

	//go:embed templates
	embedFSTemplates embed.FS

	httpServer *http.Server
	database   *db.Database
	log        = logger.New("app")
)

const (
	addr = "0.0.0.0:8080" // Hardcoded for now, as Docker is the intended deployment option
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

// Starts the app (db+http)
func Start() {
	database = db.MustOpen()

	e := echo.New()
	e.Use(middleware.Recover())

	// Rate limiting
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// Validator for requests
	e.Validator = &CustomValidator{validator: validator.New()}

	// Golang template for the HTML pages
	templ, err := template.New("").ParseFS(embedFSTemplates, "templates/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	e.Renderer = &echo.TemplateRenderer{Template: templ}

	// Pages routes
	e.GET("/", getPageIndex)
	e.GET("/view", getPageView)

	// Secret string API routes
	e.POST("/api/string/new", postStringNew)
	e.GET("/api/string/fetch/:id", getStringFetch)

	// Secret file API routes
	e.POST("/api/file/new", postFileNew)
	e.GET("/api/file/fetch/:id", getFileFetch)

	// Static files middleware
	e.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(embedFSStatic))))

	httpServer = &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	go func() {
		log.Info("starting http server on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

// Stops the app
func Stop() {
	log.Debug("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error(err)
		return
	}

	log.Info("http server stopped")

	database.Close()
}
