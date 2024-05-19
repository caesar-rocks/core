package core

import (
	"errors"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strconv"

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
			switch field.Type.Kind() {
			case reflect.Int:
				envValueInt, err := strconv.Atoi(envValue)
				if err != nil {
					log.Fatalf("Failed to convert environment variable to int field: %s value: %s", envName, envValue)
				}
				value.Field(i).SetInt(int64(envValueInt))
			case reflect.String:
				value.Field(i).SetString(envValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(envValue)
				if err != nil {
					log.Fatalf("Failed to convert environment variable to bool field: %s value: %s", envName, envValue)
				}
				value.Field(i).SetBool(boolValue)
			default:
				log.Fatalf("Unsupported type field: %s type: %v", envName, field.Type.Kind())
			}
		}
	}

	// Validate the environment variables
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(env); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		for _, validationError := range validationErrors {
			slog.Error(
				"Invalid environment variable",
				"field", validationError.Field(),
				"value", validationError.Value(),
				"tag", validationError.Tag(),
			)
		}

		log.Fatal("Failed to validate environment variables")
	}

	return env
}
