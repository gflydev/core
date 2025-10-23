package utils

import (
	"os"
	"reflect"
	"strconv"
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
	// Create a variable of "any" type and assign it the initial value
	var out any = init

	// Try to retrieve the parameter from the environment
	if value, ok := os.LookupEnv(key); ok {
		switch reflect.TypeOf(init).Name() {
		case "string":
			out = value
		case "int":
			if num, err := strconv.Atoi(value); err == nil {
				out = num
			}
		case "float64":
			if num, err := strconv.ParseFloat(value, 64); err == nil {
				out = num
			}
		case "bool":
			if boolean, err := strconv.ParseBool(value); err == nil {
				out = boolean
			}
		case "float32":
			if num, err := strconv.ParseFloat(value, 32); err == nil {
				out = num
			}
		}
	}

	// Perform type assertion and return the result
	return out.(V)
}
