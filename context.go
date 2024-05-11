package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
)

// CaesarCtx is a wrapper around http.ResponseWriter and *http.Request,
// that is augmented with some Caesar-specific methods.
type CaesarCtx struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	statusCode int
}

func (c *CaesarCtx) Render(component templ.Component) error {
	return component.Render(c.Request.Context(), c.ResponseWriter)
}

func (c *CaesarCtx) SendJSON(v interface{}, statuses ...int) error {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(status)

	return json.NewEncoder(c.ResponseWriter).Encode(v)
}

func (c *CaesarCtx) SendText(text string, statuses ...int) error {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.WithStatus(status)

	_, err := c.ResponseWriter.Write([]byte(text))
	return err
}

func (c *CaesarCtx) Context() context.Context {
	return c.Request.Context()
}

func (c *CaesarCtx) WithStatus(statusCode int) *CaesarCtx {
	c.statusCode = statusCode
	return c
}
