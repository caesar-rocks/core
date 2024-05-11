package core

import (
	"errors"
	"os"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

// ValidateEnvironmentVariables validates the environment variables,
// and panics if any of the environment variables aren't correctly set.
func ValidateEnvironmentVariables[T interface{}]() *T {
	var env *T = new(T)

	// Fill env with the actual  environment variables
	valueType := reflect.TypeOf(env).Elem()
	value := reflect.ValueOf(env).Elem()

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		envName := field.Tag.Get("env")

		if envName == "" {
			envName = field.Name
		}

		if envValue := os.Getenv(envName); envValue != "" {
			value.Field(i).SetString(envValue)
		}
	}

	// Validate the environment variables
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(env); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		for _, validationError := range validationErrors {
			log.Error(
				"Invalid environment variable",
				"field", validationError.Field(),
				"value", validationError.Value(),
				"tag", validationError.Tag(),
			)
		}

		log.Fatal("Failed to validate environment variables")
	}

	log.Info("Environment variables validated")

	return env
}
