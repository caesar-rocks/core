package core

import (
	"embed"
	"errors"
	"net/http"

	"github.com/a-h/templ"
)

// Handler is a function that can be used as middleware, or as a route handler.
type Handler func(*Context) error

// Router is a router that can be used to add routes.
type Router struct {
	Routes         []*Route
	Mux            *http.ServeMux
	Middleware     []Handler
	StandardRoutes map[string]http.Handler
}

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{
		Routes:         []*Route{},
		StandardRoutes: map[string]http.Handler{},
	}
}

// Handle adds a standard http.Handler to the router.
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.StandardRoutes[pattern] = handler
}

// Any adds a route that matches any HTTP method.
func (r *Router) Any(pattern string, handler Handler) *Route {
	route := &Route{Pattern: pattern, Handler: handler}
	r.Routes = append(r.Routes, route)
	return route
}

// Get adds a GET route to the router.
func (r *Router) Get(pattern string, handler Handler) *Route {
	route := &Route{Method: "GET", Pattern: pattern, Handler: handler}
	r.Routes = append(r.Routes, route)
	return route
}

// Post adds a route that only matches POST requests.
func (r *Router) Post(pattern string, handler Handler) *Route {
	route := &Route{Method: "POST", Pattern: pattern, Handler: handler}
	r.Routes = append(r.Routes, route)
	return route
}

// Patch adds a route that uses the PUT method.
func (r *Router) Patch(pattern string, handler Handler) *Route {
	route := &Route{Method: "PATCH", Pattern: pattern, Handler: handler}
	r.Routes = append(r.Routes, route)
	return route
}

// Put adds a route that handles DELETE requests.
func (r *Router) Put(pattern string, handler Handler) *Route {
	route := &Route{Method: "PUT", Pattern: pattern, Handler: handler}
	r.Routes = append(r.Routes, route)
	return route
}

// Delete adds a route that handles DELETE requests.
func (r *Router) Delete(pattern string, handler Handler) *Route {
	route := &Route{Method: "DELETE", Pattern: pattern, Handler: handler}
	r.Routes = append(r.Routes, route)
	return route
}

// Render adds a route (with a GET method) that renders a component.
func (r *Router) Render(pattern string, component templ.Component) *Route {
	return r.Get(pattern, func(ctx *Context) error {
		return ctx.Render(component)
	})
}

// Use adds middleware to the whole router (all routes).
func (r *Router) Use(handler Handler) {
	r.Middleware = append(r.Middleware, handler)
}

// ServeStaticAssets returns a provider,
// which serves static assets from the embed.FS.
func ServeStaticAssets(fs embed.FS) func(r *Router) {
	return func(r *Router) {
		fileServer := http.FileServer(http.FS(fs))
		r.Handle(
			"GET /assets/*",
			http.StripPrefix("/", fileServer),
		)
	}
}

// MakeURL returns the URL for a route.
func (r *Router) MakeURL(name string, values map[string]string) (string, error) {
	for _, route := range r.Routes {
		if route.Name == name {
			url := route.MakeURL(values)
			if url != "" {
				return url, nil
			}
		}
	}
	return "", errors.New("route not found")
}
