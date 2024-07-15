package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	// charmbracelt log
	"github.com/charmbracelet/log"
)

type App struct {
	Config       *AppConfig
	Router       *Router
	ErrorHandler *ErrorHandler
}

type AppConfig struct {
	Addr  string
	Debug bool
}

func NewApp(cfg *AppConfig) *App {
	return &App{
		Config: cfg,
	}
}

func (app *App) buildHTTPHandler(router *Router, route *Route, errorHandler *ErrorHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		timeStart := time.Now()

		ctx := &Context{ResponseWriter: w, Request: r, statusCode: http.StatusOK, nextCalled: true, router: router}

		// Run the global middleware
		for _, middleware := range router.Middleware {
			ctx.nextCalled = false
			err = middleware(ctx)
			if err != nil || !ctx.nextCalled {
				break
			}
		}

		// Run the route-related middleware, if no error occurred before
		if err == nil {
			for _, middleware := range route.Middleware {
				ctx.nextCalled = false
				err = middleware(ctx)
				if err != nil || !ctx.nextCalled {
					break
				}
			}
		}

		defer func() {
			if err != nil {
				// Defer the error handling
				log.Error("error handling request", "err", err, "duration", time.Since(timeStart))
				errorHandler.Handle(ctx, err)
			} else {
				// Log the request completion
				if app.Config.Debug {
					log.Info(
						"request completed",
						"method", r.Method,
						"path", r.URL.Path,
						"status", strconv.Itoa(ctx.statusCode),
						"statusText", http.StatusText(ctx.statusCode),
						"duration", time.Since(timeStart),
					)
				}
			}
		}()

		// Run the route handler, if no error occurred during the middleware execution
		if err == nil && ctx.nextCalled {
			err = route.Handler(ctx)
		}
	}
}

func (app *App) newHTTPMux(router *Router, errorHandler *ErrorHandler) *http.ServeMux {
	router.Mux = http.NewServeMux()

	for _, route := range router.Routes {
		handler := app.buildHTTPHandler(router, route, errorHandler)

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

	// Map the standard routes
	for pattern, handler := range router.StandardRoutes {
		router.Mux.Handle(pattern, handler)
	}

	// Map the home route to 404 if not already mapped
	router.Mux.HandleFunc("/", app.buildHTTPHandler(router, &Route{
		Handler: func(ctx *Context) error {
			if ctx.Request.URL.Path == "/" {
				return NewError(http.StatusNotFound)
			}
			ctx.Next()
			return nil
		},
	}, errorHandler))

	return router.Mux
}

func NewHTTPServer(addr string, mux *http.ServeMux) *http.Server {
	// Create the server
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return srv
}

func (app *App) Run() {
	// Ensure the router and error handler are set
	if app.Router == nil {
		log.Fatal("Router must be set")
	}

	if app.ErrorHandler == nil {
		log.Fatal("ErrorHandler must be set")
	}

	// Create HTTP mux
	mux := app.newHTTPMux(app.Router, app.ErrorHandler)

	// Create and start the HTTP server
	server := NewHTTPServer(app.Config.Addr, mux)
	go func() {
		log.Info("starting http server", "addr", app.Config.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("error starting http server", "err", err)
		}
	}()

	// Handle shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	<-shutdownChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
