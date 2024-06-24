package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// Context is a wrapper around http.ResponseWriter and *http.Request,
// that is augmented with some Caesar-specific methods.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	statusCode int
	nextCalled bool
	router     *Router
}

func (c *Context) Render(component templ.Component) error {
	return component.Render(c.Request.Context(), c.ResponseWriter)
}

func (c *Context) SendJSON(v interface{}, statuses ...int) error {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(status)

	return json.NewEncoder(c.ResponseWriter).Encode(v)
}

func (c *Context) SendText(text string, statuses ...int) error {
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}

	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.WithStatus(status)

	_, err := c.ResponseWriter.Write([]byte(text))
	return err
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}

func (c *Context) WithStatus(statusCode int) *Context {
	c.statusCode = statusCode
	return c
}

// PathValue returns the value of a path parameter.
func (c *Context) PathValue(key string) string {
	return c.Request.PathValue(key)
}

// DecodeJSON decodes the JSON body of the request into the provided value.
func (c *Context) DecodeJSON(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// Redirect redirects the client to the provided URL.
func (ctx *Context) Redirect(to string) error {
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.WithStatus(http.StatusSeeOther).SetHeader("HX-Redirect", to)
		return nil
	}
	http.Redirect(ctx.ResponseWriter, ctx.Request, to, http.StatusSeeOther)
	return nil
}

// RedirectBack redirects the client to the previous page.
func (ctx *Context) RedirectBack() error {
	return ctx.Redirect(ctx.Request.Referer())
}

// Validate validates the request body or form values.
// It returns the data, the validation errors, and a boolean indicating if the data is valid.
func Validate[T interface{}](ctx *Context) (data *T, validationErrors map[string]string, ok bool) {
	data = new(T)
	var errs validator.ValidationErrors

	if ctx.Request.Header.Get("Content-Type") == "application/json" {
		ctx.DecodeJSON(&data)
	} else {
		ctx.Request.ParseForm()

		decoder := form.NewDecoder()
		if err := decoder.Decode(&data, ctx.Request.Form); err != nil {
			errors.As(err, &errs)
		}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		errors.As(err, &errs)
	}

	// Turn `errs` into a map.
	var errors = make(map[string]string)
	for _, err := range errs {
		errors[err.StructField()] = fmt.Sprintf("This field does not meet the following rule: \"%s\".", err.Tag())
	}

	return data, errors, len(errs) == 0
}

// SetHeader sets a header in the response.
func (ctx *Context) SetHeader(key string, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// GetHeader returns the value of a header in the request.
func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

// Next marks the middleware as having called the next handler in the chain.
func (c *Context) Next() {
	c.nextCalled = true
}

// WantsJSON returns true if the client accepts JSON responses.
func (c *Context) WantsJSON() bool {
	return c.Request.Header.Get("Accept") == "application/json"
}

// SetSSEHeaders sets the headers for Server-Sent Events
func (ctx *Context) SetSSEHeaders() {
	ctx.SetHeader("Content-Type", "text/event-stream")
	ctx.SetHeader("Cache-Control", "no-cache")
	ctx.SetHeader("Connection", "keep-alive")
	ctx.SetHeader("Access-Control-Allow-Origin", "*")
	ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// SendSSE sends a Server-Sent Event to the client
func (ctx *Context) SendSSE(name string, data string) error {
	ctx.ResponseWriter.Write([]byte("event: " + name + "\n"))
	ctx.ResponseWriter.Write([]byte("data: " + data + "\n\n"))
	flusher, ok := ctx.ResponseWriter.(http.Flusher)
	if !ok {
		return errors.New("streaming unsupported")
	}
	flusher.Flush()
	return nil
}

// MakeURL returns the URL for a route with the given name.
func (ctx *Context) MakeURL(name string, params map[string]string) (string, error) {
	return ctx.router.MakeURL(name, params)
}
