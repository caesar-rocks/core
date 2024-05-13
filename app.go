package core

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/log"

	"go.uber.org/fx"
)

type App struct {
	Providers []any
	Invokers  []any
	Config    *AppConfig
}

type AppConfig struct {
	Addr string
}

func NewApp(cfg *AppConfig) *App {
	return &App{
		Config: cfg,
	}
}

func (app *App) RegisterProviders(providers ...any) {
	app.Providers = append(app.Providers, providers...)
}

func (app *App) RegisterInvokers(invokers ...any) {
	app.Invokers = append(app.Invokers, invokers...)
}

func buildHTTPHandler(router *Router, route *Route, errorHandler *ErrorHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		timeStart := time.Now()

		ctx := &CaesarCtx{ResponseWriter: w, Request: r, statusCode: http.StatusOK}

		// Run the global middleware
		for _, middleware := range router.Middleware {
			err := middleware(ctx)
			if err != nil {
				break
			}
		}

		// Run the route-related middleware, if no error occurred before
		if err == nil {
			for _, middleware := range route.Middleware {
				err = middleware(ctx)
				if err != nil {
					break
				}
			}
		}

		defer func() {
			if err != nil {
				// Defer the error handling
				log.Error("Error handling request", "err", err, "duration", time.Since(timeStart))
				errorHandler.Handle(ctx, err)
			} else {
				// Log the request completion
				log.Info(
					"Request completed",
					"method", r.Method,
					"path", r.URL.Path,
					"status", strconv.Itoa(ctx.statusCode),
					"statusText", http.StatusText(ctx.statusCode),
					"duration", time.Since(timeStart),
				)
			}
		}()

		// Run the route handler, if no error occurred
		// during the middleware execution
		if err == nil {
			err = route.Handler(ctx)
		}
	}
}

func NewHTTPMux(router *Router, errorHandler *ErrorHandler) *http.ServeMux {
	router.Mux = http.NewServeMux()

	for _, route := range router.Routes {
		var handler http.HandlerFunc

		log.Info("Register route", "method", route.Method, "pattern", route.Pattern)

		handler = buildHTTPHandler(router, route, errorHandler)

		if route.Method == "" {
			router.Mux.HandleFunc(route.Pattern, handler)
			if route.Pattern != "/" {
				router.Mux.HandleFunc(route.Pattern+"/", handler)
			}
		} else {
			router.Mux.HandleFunc(route.Method+" "+route.Pattern, handler)
			if route.Pattern != "/" {
				router.Mux.HandleFunc(route.Method+" "+route.Pattern+"/", handler)
			}
		}
	}

	// Map the home route to 404 if not already mapped
	router.Mux.HandleFunc("/", buildHTTPHandler(router, &Route{
		Handler: func(ctx *CaesarCtx) error {
			return NewError(http.StatusNotFound)
		},
	}, errorHandler))

	return router.Mux
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	// Create the HTTP error logger
	logger := log.NewWithOptions(os.Stderr, log.Options{Prefix: "http"})
	stdlog := logger.StandardLog(log.StandardLogOptions{
		ForceLevel: log.ErrorLevel,
	})

	// Create the server
	srv := &http.Server{
		Addr:     os.Getenv("ADDR"),
		ErrorLog: stdlog,
		Handler:  mux,
	}

	// Register the server with the lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("HTTP server started", "addr", os.Getenv("ADDR"))
			go func() {
				err := srv.ListenAndServe()
				if err != nil {
					log.Fatal("Error starting HTTP server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func (app *App) Run() {
	fx.New(
		fx.Provide(NewHTTPMux, NewHTTPServer),
		fx.Provide(app.Providers...),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(app.Invokers...),
	).Run()
}
