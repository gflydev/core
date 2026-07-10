package utils

import (
	"os"
	"reflect"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// Getenv retrieves an environment variable and parses it into the specified type,
// falling back to a default value if the variable is not found.
//
// Parameters:
//   - key (string): The name of the environment variable.
//   - init (V): The default value to use if the environment variable is not set.
//
// Returns:
//   - V: The value of the environment variable parsed into the specified type,
//     or the default value if the variable is not found.
func Getenv[V any](key string, init V) V {
	// Try to retrieve the parameter from the environment.
	value, ok := os.LookupEnv(key)
	if !ok {
		return init
	}

	// Switch on the reflect Kind (not the type Name): this also handles named
	// types such as `type Level int` or time.Duration whose Name() would not be
	// "int". A nil interface would make reflect.TypeOf panic, so guard for it.
	t := reflect.TypeOf(init)
	if t == nil {
		return init
	}

	var out any = init
	switch t.Kind() {
	case reflect.String:
		out = value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, err := strconv.ParseInt(value, 10, 64); err == nil {
			out = reflect.ValueOf(num).Convert(t).Interface()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if num, err := strconv.ParseUint(value, 10, 64); err == nil {
			out = reflect.ValueOf(num).Convert(t).Interface()
		}
	case reflect.Float32, reflect.Float64:
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			out = reflect.ValueOf(num).Convert(t).Interface()
		}
	case reflect.Bool:
		if boolean, err := strconv.ParseBool(value); err == nil {
			out = reflect.ValueOf(boolean).Convert(t).Interface()
		}
	}

	// Perform type assertion and return the result
	return out.(V)
}
