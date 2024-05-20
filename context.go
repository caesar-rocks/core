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

// PathValue returns the value of a path parameter.
func (c *CaesarCtx) PathValue(key string) string {
	return c.Request.PathValue(key)
}

// DecodeJSON decodes the JSON body of the request into the provided value.
func (c *CaesarCtx) DecodeJSON(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// Redirect redirects the client to the provided URL.
func (ctx *CaesarCtx) Redirect(to string) error {
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.WithStatus(http.StatusSeeOther).SetHeader("HX-Redirect", to)
		return nil
	}
	http.Redirect(ctx.ResponseWriter, ctx.Request, to, http.StatusSeeOther)
	return nil
}

// Validate validates the request body or form values.
// It returns the data, the validation errors, and a boolean indicating if the data is valid.
func Validate[T interface{}](ctx *CaesarCtx) (data *T, validationErrors map[string]string, ok bool) {
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
func (ctx *CaesarCtx) SetHeader(key string, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// GetHeader returns the value of a header in the request.
func (ctx *CaesarCtx) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}
