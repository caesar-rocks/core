package core

import (
	"strings"
)

// Route is a route that can be added to a Router.
type Route struct {
	Method     string
	Pattern    string
	Handler    Handler
	Middleware []Handler
	Name       string
}

// Use adds middleware to the route.
func (route *Route) Use(handler Handler) *Route {
	route.Middleware = append(route.Middleware, handler)
	return route
}

// As sets the name of the route.
func (route *Route) As(name string) *Route {
	route.Name = name
	return route
}

// MakeURL generates a URL for the route.
func (route *Route) MakeURL(params map[string]string) string {
	url := route.Pattern

	// Replace named parameters in the URL.
	for key, value := range params {
		url = strings.Replace(url, ":"+key, value, 1)
	}

	return ""
}
