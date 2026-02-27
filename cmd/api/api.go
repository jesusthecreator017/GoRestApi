package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jesusthecreator017/fswithgo/cmd/api/middleware"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr       string
	db         dbConfig
	corsOrigin string
	jwtSecret  string
}

// dbConfig holds database connection settings.
// These are read from environment variables in main.go.
type dbConfig struct {
	dsn         string        // PostgreSQL connection string
	maxConns    int           // max open connections in the pool
	maxIdleTime time.Duration // close connections idle longer than this
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	// Issue
	mux.Handle("GET /v1/issues", middleware.RequiredAuth(http.HandlerFunc(app.listIssueHandler)))
	mux.HandleFunc("GET /v1/issues/{id}", app.getIssueHandler)

	// Protected Issue routed
	mux.Handle("POST /v1/issues", middleware.RequiredAuth(http.HandlerFunc(app.createIssueHandler)))
	mux.Handle("DELETE /v1/issues/{id}", middleware.RequiredAuth(http.HandlerFunc(app.deleteIssueHandler)))
	mux.Handle("PATCH /v1/issues/{id}/status", middleware.RequiredAuth(http.HandlerFunc(app.updateIssueStatusHandler)))

	// Users
	mux.HandleFunc("POST /v1/users/register", app.registerUserHandler)
	mux.HandleFunc("POST /v1/users/login", app.loginUserHandler)

	// Admin
	mux.Handle("GET /v1/admin/stats", middleware.RequiredAuth(middleware.RequiredAdmin(http.HandlerFunc(app.adminStatsHandler))))

	return mux
}

func (app *application) run(mux *http.ServeMux) error {

	// Global Middleware
	stack := middleware.CreateStack(
		middleware.CORS(app.config.corsOrigin),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logging,
		middleware.GlobalAuth(app.config.jwtSecret),
		middleware.Timeout(time.Second*25),
	)

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      stack(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("Listening on port%s\n", app.config.addr)
	return srv.ListenAndServe()
}
