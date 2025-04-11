package utils

import (
	"fmt"
	"net/url"
)

// RequestPath parse request path string.
//

// RequestPath parses the request path string from a given URL.
//
// Parameters:
//   - urlStr string: The URL string from which the path needs to be extracted.
//     Example: "https://902-local.s3.us-west-1.amazonaws.com/news/Avatar2023-jcer.jpeg"
//
// Eg:
//   - Get `/news/Avatar2023-jcer.jpeg` from `https://902-local.s3.us-west-1.amazonaws.com/news/Avatar2023-jcer.jpeg`
//   - Get `/storage/tmp/63e85ba1.jpeg` from `https://www.dancefitvn.com/storage/tmp/63e85ba1.jpeg`
//
// Returns:
// - string: The extracted path from the given URL string.
// - error: An error if the URL parsing fails.
//
// Variables in function:
// - u *url.URL: Parsed URL object obtained from the input string.
// - err error: Stores any error that occurs during URL parsing.
func RequestPath(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return "", err
	}

	return u.Path, nil
}

// RequestParam parses a specific parameter from a given URL string.
//
// Parameters:
//   - urlStr string: The URL string from which the parameter needs to be extracted.
//     Example: "http://localhost:7789/api/v1/uploads?G-Key=e2c32ed0807ce086083281f131&G-Time=20240803130551"
//   - param string: The name of the query parameter to fetch.
//     Example: "G-Key"
//
// Eg:
//   - Get `G-Key=e2c32ed0807ce086083281f131&G-Time=20240803130551`
//     from `http://localhost:7789/api/v1/uploads?G-Key=e2c32ed0807ce086083281f131&G-Time=20240803130551`
//
// Returns:
//   - string: The value of the specified query parameter. If the parameter is not found, an empty string is returned.
//   - error: An error if the URL parsing fails.
//
// Variables in function:
//   - u *url.URL: Parsed URL object obtained from the input string.
//   - err error: Stores any error that occurs during the URL parsing process.
func RequestParam(urlStr, param string) (string, error) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return "", err
	}

	return u.Query().Get(param), nil
}

// RequestURL parses the request URL string and constructs the URL using the scheme, host, and path.
//
// Parameters:
//   - urlStr string: The input URL string that needs to be parsed and reconstructed.
//     Example: "https://example.com:8080/path/to/resource"
//
// Returns:
//   - string: The reconstructed URL string in the format "scheme://host/path".
//   - error: An error if the URL parsing fails.
func RequestURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path), nil
}
