package core

import (
	"encoding/json"
	"fmt"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/valyala/fasthttp"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ====================================================================
//                               Context
// ====================================================================

// Ctx HTTP request context
type Ctx struct {
	app    *GFly                // Reference to gFly application.
	root   *fasthttp.RequestCtx // Reference to the root request context of fasthttp.
	router *Router              // Reference to the router.
	data   Data                 // Data holds request context data.
}

// Root returns the original HTTP request context from fasthttp.
//
// Returns:
//   - *fasthttp.RequestCtx: The root HTTP request context.
func (c *Ctx) Root() *fasthttp.RequestCtx {
	return c.root
}

// Router returns the original router associated with the context.
//
// Returns:
//   - *Router: The router instance associated with the current request context.
func (c *Ctx) Router() *Router {
	return c.router
}

// ====================================================================
//                         Ctx - Header
// ====================================================================

// RequestHandler A wrapper of fasthttp.RequestHandler.
//
// Parameters:
//   - ctx (*Ctx): The current HTTP request context.
//
// Returns:
//   - error: An error if the handling fails, otherwise nil.
type RequestHandler func(ctx *Ctx) error

// IHandler Interface for a handler request.
type IHandler interface {
	// Validate validates the request context.
	//
	// Parameters:
	//   - c (*Ctx): The current HTTP request context.
	//
	// Returns:
	//   - error: An error if validation fails, otherwise nil.
	Validate(c *Ctx) error

	// Handle handles the request context.
	//
	// Parameters:
	//   - c (*Ctx): The current HTTP request context.
	//
	// Returns:
	//   - error: An error if handling fails, otherwise nil.
	Handle(c *Ctx) error
}

// Endpoint Default handler.
//
// Implements:
//   - Validate: A no-op validation method.
//   - Handle: A no-op handling method.
type Endpoint struct{}

// Validate validates the request context.
//
// Parameters:
//   - c (*Ctx): The current HTTP request context.
//
// Returns:
//   - error: Always nil (no validation logic implemented).
func (e *Endpoint) Validate(c *Ctx) error {
	return nil
}

// Handle handles the request context.
//
// Parameters:
//   - c (*Ctx): The current HTTP request context.
//
// Returns:
//   - error: Always nil (no handling logic implemented).
func (e *Endpoint) Handle(c *Ctx) error {
	return nil
}

// Page Abstract web page.
//
// Inherits:
//   - Endpoint: Base endpoint for page handling logic.
type Page struct {
	Endpoint
}

// Api Abstract API.
//
// Inherits:
//   - Endpoint: Base endpoint for API handling logic.
type Api struct {
	Endpoint
}

// ====================================================================
//                         Ctx - Header Data
// ====================================================================

type IHeader interface {
	// Status sets the response's HTTP status code.
	//
	// Parameters:
	//   - code (int): The HTTP status code to set.
	//
	// Returns:
	//   - *Ctx: The current HTTP context.
	Status(code int) *Ctx
	// ContentType sets the response's HTTP content type.
	//
	// Parameters:
	//   - mime (string): The content type to set.
	//
	// Returns:
	//   - *Ctx: The current HTTP context.
	ContentType(mime string) *Ctx
	// SetHeader sets the response's HTTP header field with the specified key and value.
	//
	// Parameters:
	//   - key (string): The header key.
	//   - val (string): The header value.
	//
	// Returns:
	//   - *Ctx: The current HTTP context.
	SetHeader(key, val string) *Ctx
	// SetCookie sets a cookie in the response's HTTP header.
	//
	// Parameters:
	//   - key (string): The cookie name.
	//   - value (string): The cookie value.
	//
	// Returns:
	//   - *Ctx: The current HTTP context.
	SetCookie(key, value string) *Ctx
	// GetCookie retrieves a cookie value from the request's HTTP header.
	//
	// Parameters:
	//   - key (string): The name of the cookie.
	//
	// Returns:
	//   - string: The value of the cookie; an empty string if not found.
	GetCookie(key string) string
	// GetReqHeaders returns all HTTP request headers.
	//
	// Returns:
	//   - map[string][]string: A map of header keys and their respective values.
	GetReqHeaders() map[string][]string
	// Path retrieves the URI path of the current request.
	//
	// Returns:
	//   - string: The URI path as a string.
	Path() string
}

// Status sets the response's HTTP status code.
//
// Parameters:
//   - code (int): The HTTP status code to set.
//
// Returns:
//   - *Ctx: The current HTTP context.
func (c *Ctx) Status(code int) *Ctx {
	c.root.Response.SetStatusCode(code)

	return c
}

// ContentType sets the response's HTTP content type.
//
// Parameters:
//   - mime (string): The content type to set.
//
// Returns:
//   - *Ctx: The current HTTP context.
func (c *Ctx) ContentType(mime string) *Ctx {
	c.root.Response.Header.SetContentType(mime)

	return c
}

// SetHeader sets the response's HTTP header field with the specified key and value.
//
// Parameters:
//   - key (string): The header key.
//   - val (string): The header value.
//
// Returns:
//   - *Ctx: The current HTTP context.
func (c *Ctx) SetHeader(key, val string) *Ctx {
	c.root.Response.Header.Set(key, val)

	return c
}

// SetCookie sets a cookie in the response's HTTP header.
//
// Parameters:
//   - key (string): The cookie name.
//   - value (string): The cookie value.
//
// Returns:
//   - *Ctx: The current HTTP context.
func (c *Ctx) SetCookie(key, value string) *Ctx {
	cook := fasthttp.Cookie{}
	cook.SetKey(key)
	cook.SetValue(value)
	cook.SetMaxAge(3600000)

	c.root.Response.Header.SetCookie(&cook)

	return c
}

// GetCookie retrieves a cookie value from the request's HTTP header.
//
// Parameters:
//   - key (string): The name of the cookie.
//
// Returns:
//   - string: The value of the cookie; an empty string if not found.
func (c *Ctx) GetCookie(key string) string {
	return string(c.root.Request.Header.Cookie(key))
}

// GetReqHeaders returns all HTTP request headers.
//
// Returns:
//   - map[string][]string: A map of header keys and their respective values.
func (c *Ctx) GetReqHeaders() map[string][]string {
	headers := make(map[string][]string)
	c.root.Request.Header.VisitAll(func(k, v []byte) {
		key := utils.UnsafeStr(k)
		headers[key] = append(headers[key], utils.UnsafeStr(v))
	})

	return headers
}

// Path retrieves the URI path of the current request.
//
// Returns:
//   - string: The URI path as a string.
func (c *Ctx) Path() string {
	return string(c.root.URI().Path())
}

// ====================================================================
//                         Ctx - Response Data
// ====================================================================

type IResponse interface {
	// Success sends a successful JSON response.
	//
	// Parameters:
	//   - data (any): The data to include in the JSON response.
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	Success(data any) error

	// Error sends an error JSON response.
	//
	// Parameters:
	//   - data (any): The data to include in the JSON response.
	//   - httpCode (...int): Optional HTTP status code to set (default: 400 Bad Request).
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	Error(data any, httpCode ...int) error

	// NoContent sends a response with no content.
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	NoContent() error

	// View renders and sends a template page response.
	//
	// Parameters:
	//   - template (string): The name or path of the template to render.
	//   - data (Data): The data to inject into the template.
	//
	// Returns:
	//   - error: An error if template rendering or response generation fails, otherwise nil.
	View(template string, data Data) error

	// JSON sends a JSON response.
	//
	// Parameters:
	//   - data (Data): The data to include in the JSON response.
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	JSON(data Data) error

	// HTML sends an HTML response.
	//
	// Parameters:
	//   - body (string): The HTML content to include in the response.
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	HTML(body string) error

	// String sends a plain text response.
	//
	// Parameters:
	//   - body (string): The plain text content to include in the response.
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	String(body string) error

	// Raw sends a raw byte response.
	//
	// Parameters:
	//   - body ([]byte): The raw byte content to include in the response.
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	Raw(body []byte) error

	// Stream sends a streaming response.
	//
	// Parameters:
	//   - stream (io.Reader): The data stream to send in the response.
	//   - size (...int): Optional body size for the stream (default: -1 for unknown size).
	//
	// Returns:
	//   - error: An error if the response generation fails, otherwise nil.
	Stream(stream io.Reader, size ...int) error

	// Redirect sends a redirect response.
	//
	// Parameters:
	//   - path (string): The URL path or full URL to redirect to.
	//
	// Returns:
	//   - error: An error if the redirect fails, otherwise nil.
	Redirect(path string) error

	// Download sends a file as an attachment.
	//
	// Parameters:
	//   - file (string): The path of the file to send.
	//   - filename (...string): Optional custom filename for the attachment.
	//
	// Returns:
	//   - error: An error if the file transfer fails, otherwise nil.
	Download(file string, filename ...string) error

	// File sends a file as a response.
	//
	// Parameters:
	//   - file (string): The path of the file to send.
	//   - compress (...bool): Optional compression flag (default: false).
	//
	// Returns:
	//   - error: An error if the file transfer fails, otherwise nil.
	File(file string, compress ...bool) error
}

// Success sends a successful JSON response.
//
// Parameters:
//   - data (any): The data to include in the JSON response.
//
// Returns:
//   - error: An error if the response generation fails, otherwise nil.
func (c *Ctx) Success(data any) error {
	c.root.Response.SetStatusCode(StatusOK)
	return c.JSON(data)
}

// Error sends an error JSON response with a specific HTTP status code.
//
// Parameters:
//   - data (any): The data to include in the JSON response.
//   - httpCode (...int): Optional HTTP status code to set (default: 400 Bad Request).
//
// Returns:
//   - error: An error if the response generation fails, otherwise nil.
func (c *Ctx) Error(data any, httpCode ...int) error {
	if len(httpCode) > 0 {
		c.root.Response.SetStatusCode(httpCode[0])
	} else {
		c.root.Response.SetStatusCode(StatusBadRequest)
	}

	// Set response content
	_ = c.JSON(data)

	return errors.UnknownError
}

// NoContent sends a response with no content.
//
// Returns:
//   - error: An error if the response generation fails, otherwise nil.
func (c *Ctx) NoContent() error {
	c.root.Response.SetStatusCode(StatusNoContent)
	return nil
}

// View renders a template file and sends the resulting HTML as the HTTP response.
//
// Parameters:
//   - template (string): The name or path of the template file to render.
//   - data (Data): The data to pass into the template for rendering.
//
// Returns:
//   - error: An error if the template rendering or sending the response fails.
func (c *Ctx) View(template string, data Data) error {
	c.ContentType(MIMETextHTMLCharsetUTF8)

	data["appName"] = AppName
	data["appURL"] = AppURL
	data["appCode"] = AppCode
	data["appDebug"] = AppDebug

	return view.Writer(template, data, c.root.Response.BodyWriter())
}

// JSON serializes the given data into JSON and sends it as the HTTP response.
//
// Parameters:
//   - data (Data): The data to serialize into JSON.
//
// Returns:
//   - error: An error if the JSON serialization or sending the response fails.
func (c *Ctx) JSON(data any) error {
	c.root.Response.Header.SetContentType(MIMEApplicationJSONCharsetUTF8)

	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.Raw(marshal)
}

// HTML sets the HTTP response body content to the provided HTML string.
//
// Parameters:
//   - body (string): The HTML content to send as the response body.
//
// Returns:
//   - error: An error if setting the response body fails.
func (c *Ctx) HTML(body string) error {
	c.root.Response.Header.SetContentType(MIMETextHTMLCharsetUTF8)

	return c.Raw([]byte(body))
}

// String sets the HTTP response body content to the provided plain text string.
//
// Parameters:
//   - body (string): The plain text content to send as the response body.
//
// Returns:
//   - error: An error if setting the response body fails.
func (c *Ctx) String(body string) error {
	c.root.Response.Header.SetContentType(MIMETextPlain)

	return c.Raw([]byte(body))
}

// Raw sets the HTTP response body directly using the given raw byte slice.
//
// Parameters:
//   - body ([]byte): The raw byte slice to send as the response body.
//
// Returns:
//   - error: An error if setting the response body fails.
func (c *Ctx) Raw(body []byte) error {
	c.root.Response.SetBodyRaw(body)

	return nil
}

// Stream sets the HTTP response body to a data stream with an optional size.
//
// Parameters:
//   - stream (io.Reader): The data stream to send as the response body.
//   - size (...int): Optional argument to specify the size of the stream body. Defaults to -1 (unknown size).
//
// Returns:
//   - error: An error if setting the response body stream fails.
func (c *Ctx) Stream(stream io.Reader, size ...int) error {
	if len(size) > 0 && size[0] >= 0 {
		c.root.Response.SetBodyStream(stream, size[0])
	} else {
		c.root.Response.SetBodyStream(stream, -1)
	}

	return nil
}

// Redirect sends a redirect response to the specified path.
//
// Parameters:
//   - path (string): The URL path or full URL to redirect to.
//
// Returns:
//   - error: An error if setting the redirect fails.
func (c *Ctx) Redirect(path string) error {
	// Check the `path` is relative URI. Build full internal URL
	if !strings.HasPrefix(path, SchemaHTTP) {
		path = AppURL + path
	}

	c.root.Redirect(path, StatusMovedPermanently)
	log.Tracef("redirect to %s", path)

	return errors.UnknownError
}

// Download transfers the file from the provided path as an attachment.
// Typically, browsers will prompt the user to download the file.
// By default, the Content-Disposition header's filename parameter is set to the provided file path's basename.
// Optionally, this behavior can be overridden with a custom filename.
//
// Parameters:
//   - file (string): The path of the file to download.
//   - filename (...string): Optional custom filename for the Content-Disposition header.
//
// Returns:
//   - error: An error if the file transfer fails, otherwise nil.
func (c *Ctx) Download(file string, filename ...string) error {
	var fName string

	if len(filename) > 0 {
		fName = filename[0]
	} else {
		fName = filepath.Base(file)
	}
	c.root.Response.Header.Set(HeaderContentDisposition, `attachment; filename="`+utils.QuoteStr(fName)+`"`)

	return c.File(file)
}

var (
	sendFileOnce    sync.Once
	sendFileFS      *fasthttp.FS
	sendFileHandler fasthttp.RequestHandler
)

// File transfers the file from the given path.
//
// The file is not compressed by default, enable this option by passing 'true' as an argument.
//
// Parameters:
//   - file (string): The path of the file to be transferred.
//   - compress (...bool): Optional argument to enable compression. Default is 'false'. Pass 'true' to enable.
//
// Returns:
//   - error: An error if the file transfer fails, otherwise nil.
//
// Behavior:
//   - Automatically sets the Content-Type response HTTP header field based on the file's extension.
//   - If the file path is relative, it converts it to an absolute path.
//   - Handles errors when file paths are invalid or if the file is not found.
//   - Configures HTTP headers for compression if explicitly enabled.
func (c *Ctx) File(file string, compress ...bool) error {
	// Save the filename, we will need it in the error message if the file isn't found
	filename := file

	// https://github.com/valyala/fasthttp/blob/c7576cc10cabfc9c993317a2d3f8355497bea156/fs.go#L129-L134
	sendFileOnce.Do(func() {
		const cacheDuration = 10 * time.Second
		sendFileFS = &fasthttp.FS{
			Root:                 "",
			AllowEmptyRoot:       true,
			GenerateIndexPages:   false,
			AcceptByteRange:      true,
			Compress:             true,
			CompressedFileSuffix: c.app.config.CompressedFileSuffix,
			CacheDuration:        cacheDuration,
			IndexNames:           []string{"index.html"},
			PathNotFound: func(ctx *fasthttp.RequestCtx) {
				ctx.Response.SetStatusCode(StatusNotFound)
			},
		}
		sendFileHandler = sendFileFS.NewRequestHandler()
	})

	// Disable compression
	if len(compress) == 0 || !compress[0] {
		// https://github.com/valyala/fasthttp/blob/7cc6f4c513f9e0d3686142e0a1a5aa2f76b3194a/fs.go#L55
		c.root.Request.Header.Del(HeaderAcceptEncoding)
	}
	// copy of https://github.com/valyala/fasthttp/blob/7cc6f4c513f9e0d3686142e0a1a5aa2f76b3194a/fs.go#L103-L121 with small adjustments
	if file == "" || !filepath.IsAbs(file) {
		// extend relative path to absolute path
		hasTrailingSlash := file != "" && (file[len(file)-1] == '/' || file[len(file)-1] == '\\')

		var err error
		file = filepath.FromSlash(file)
		if file, err = filepath.Abs(file); err != nil {
			return fmt.Errorf("failed to determine abs file path: %w", err)
		}
		if hasTrailingSlash {
			file += "/"
		}
	}
	// convert the path to forward slashes regardless the OS in order to set the URI properly
	// the handler will convert back to OS path separator before opening the file
	file = filepath.ToSlash(file)

	// Restore the original requested URL
	originalURL := utils.CopyStr(c.OriginalURL())
	defer c.root.Request.SetRequestURI(originalURL)
	// Set new URI for fileHandler
	c.root.Request.SetRequestURI(file)
	// Save status code
	status := c.root.Response.StatusCode()
	// Serve file
	sendFileHandler(c.root)
	// Get the status code which is set by fasthttp
	fsStatus := c.root.Response.StatusCode()
	// Set the status code set by the user if it is different from the fasthttp status code and 200
	if status != fsStatus && status != StatusOK {
		c.Status(status)
	}
	// Check for error
	if status != StatusNotFound && fsStatus == StatusNotFound {
		return errors.FileNotFound{
			FileName: filename,
		}
	}

	return nil
}

// Compress compresses the HTTP response body using the best available compression method
// based on the "Accept-Encoding" request header.
//
// # TODO Need more checking
//
// Parameters:
//   - body ([]byte): The raw response body to be compressed.
//
// Returns:
//   - error: An error if the compression process fails or encounters issues.
func (c *Ctx) Compress(body []byte) error {
	ctx := c.root

	switch {
	case ctx.Request.Header.HasAcceptEncodingBytes([]byte(StrGzip)):
		_, err := fasthttp.WriteGzip(ctx.Response.BodyWriter(), body)
		if err != nil {
			return err
		}
	case ctx.Request.Header.HasAcceptEncodingBytes([]byte(StrBrotli)):
		_, err := fasthttp.WriteBrotli(ctx.Response.BodyWriter(), body)
		if err != nil {
			return err
		}
	case ctx.Request.Header.HasAcceptEncodingBytes([]byte(StrDeflate)):
		_, err := fasthttp.WriteDeflate(ctx.Response.BodyWriter(), body)
		if err != nil {
			return err
		}
	default:
		ctx.Response.SetBodyRaw(body)
	}

	return nil
}

// ====================================================================
//  Ctx - Request Data (Form|MultipartForm|FormFile, Query, Path, RAW)
// ====================================================================

// UploadedFile represents the information of an uploaded file.
type UploadedFile struct {
	Field string // Field name in the form where the file was uploaded.
	Name  string // Name of the uploaded file.
	Path  string // Path where the uploaded file is saved.
	Size  int64  // Size of the uploaded file in bytes.
}

// IRequestData provides an interface for handling HTTP request data, such as form values, queries,
// and uploaded files, as well as parsing body content into specific data structures.
type IRequestData interface {
	// ParseBody parses the HTTP request body into the provided struct.
	//
	// Parameters:
	//   - data (any): A pointer to a struct where the body data will be unmarshalled.
	//
	// Returns:
	//   - error: An error if parsing fails, otherwise nil.
	ParseBody(data any) error

	// ParseQuery parses the query string into the provided struct.
	//
	// Parameters:
	//   - data (any): A pointer to a struct where the query data will be unmarshalled.
	//
	// Returns:
	//   - error: An error if parsing fails, otherwise nil.
	ParseQuery(data *Data) error

	// FormVal retrieves the value of a form key from a POST or PUT request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the form data.
	//
	// Returns:
	//   - []byte: The value associated with the provided key.
	FormVal(key string) []byte

	// FormInt retrieves the integer value of a form key from a POST or PUT request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the form data.
	//
	// Returns:
	//   - int: The integer value associated with the provided key.
	//   - error: An error if the value cannot be converted to an integer.
	FormInt(key string) (int, error)

	// FormBool retrieves the boolean value of a form key from a POST or PUT request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the form data.
	//
	// Returns:
	//   - bool: The boolean value associated with the provided key.
	//   - error: An error if the value cannot be converted to a boolean.
	FormBool(key string) (bool, error)

	// FormFloat retrieves the float value of a form key from a POST or PUT request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the form data.
	//
	// Returns:
	//   - float64: The float value associated with the provided key.
	//   - error: An error if the value cannot be converted to a float.
	FormFloat(key string) (float64, error)

	// FormUpload processes and retrieves uploaded files from a POST or PUT request.
	//
	// Parameters:
	//   - files (...string): Optional list of form field names to process. If none are specified, all uploaded files are processed.
	//
	// Returns:
	//   - []UploadedFile: A list of UploadedFile details for each processed file.
	//   - error: An error if the file upload processing fails.
	FormUpload(files ...string) ([]UploadedFile, error)

	// Queries retrieves all key-value pairs from the query string of the HTTP request.
	//
	// Returns:
	//   - map[string]string: A map of key-value pairs representing the query string data.
	Queries() map[string]string

	// QueryStr retrieves a string value from the query string of the HTTP request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the query string.
	//
	// Returns:
	//   - string: The value associated with the provided key.
	QueryStr(key string) string

	// QueryInt retrieves an integer value from the query string of the HTTP request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the query string.
	//
	// Returns:
	//   - int: The integer value associated with the provided key.
	//   - error: An error if the value cannot be converted to an integer.
	QueryInt(key string) (int, error)

	// QueryBool retrieves a boolean value from the query string of the HTTP request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the query string.
	//
	// Returns:
	//   - bool: The boolean value associated with the provided key.
	//   - error: An error if the value cannot be converted to a boolean.
	QueryBool(key string) (bool, error)

	// QueryFloat retrieves a float value from the query string of the HTTP request.
	//
	// Parameters:
	//   - key (string): The key to retrieve from the query string.
	//
	// Returns:
	//   - float64: The float value associated with the provided key.
	//   - error: An error if the value cannot be converted to a float.
	QueryFloat(key string) (float64, error)

	// PathVal retrieves a value from the path parameters of the HTTP request.
	//
	// Parameters:
	//   - key (string): The name of the path parameter to retrieve.
	//
	// Returns:
	//   - string: The value associated with the provided key.
	PathVal(key string) string

	// OriginalURL retrieves the original URL of the HTTP request.
	//
	// Returns:
	//   - string: The original URL requested by the client.
	OriginalURL() string
}

// ParseBody parses the HTTP request body into the provided struct.
//
// Parameters:
//   - data (any): A pointer to a struct where the body data will be unmarshalled.
//
// Returns:
//   - error: An error if parsing fails, otherwise nil.
func (c *Ctx) ParseBody(data any) error {
	jsonData := c.root.PostBody()

	err := json.Unmarshal(jsonData, data)
	if err != nil {
		return err
	}

	return nil
}

// ParseQuery parses the query string into the provided struct.
//
// Parameters:
//   - out (any): A pointer to a struct where the query data will be unmarshalled.
//
// Returns:
//   - error: Always returns nil (currently not implemented).
func (c *Ctx) ParseQuery(data *Data) error {
	c.root.QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := string(key)
		valStr := string(value)

		// Check if the key has the array suffix []
		if strings.HasSuffix(keyStr, "[]") {
			// Remove the [] suffix to get the base key name
			baseKey := strings.TrimSuffix(keyStr, "[]")

			// Get the existing array if it exists
			var values []string
			if existingVal := data.Get(baseKey); existingVal != nil {
				if existingArray, ok := existingVal.([]string); ok {
					values = existingArray
				}
			}

			// Add the new value to the array
			values = append(values, valStr)

			// Set the array back to the base key
			data.Set(baseKey, values)
		} else {
			// Regular key, set directly
			data.Set(keyStr, valStr)
		}
	})

	return nil
}

// FormVal retrieves the value of a form key from a POST or PUT request.
//
// Parameters:
//   - key (string): The key to retrieve from the form data.
//
// Returns:
//   - []byte: The value associated with the provided key.
func (c *Ctx) FormVal(key string) []byte {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return data
}

// FormInt retrieves the integer value of a form key from a POST or PUT request.
//
// Parameters:
//   - key (string): The key to retrieve from the form data.
//
// Returns:
//   - int: The integer value associated with the provided key.
//   - error: An error if the value cannot be converted to an integer.
func (c *Ctx) FormInt(key string) (int, error) {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.Atoi(string(data))
}

// FormBool retrieves the boolean value of a form key from a POST or PUT request.
//
// Parameters:
//   - key (string): The key to retrieve from the form data.
//
// Returns:
//   - bool: The boolean value associated with the provided key.
//   - error: An error if the value cannot be converted to a boolean.
func (c *Ctx) FormBool(key string) (bool, error) {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseBool(string(data))
}

// FormFloat retrieves the float value of a form key from a POST or PUT request.
//
// Parameters:
//   - key (string): The key to retrieve from the form data.
//
// Returns:
//   - float64: The float value associated with the provided key.
//   - error: An error if the value cannot be converted to a float.
func (c *Ctx) FormFloat(key string) (float64, error) {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseFloat(string(data), 64)
}

// FormUpload processes and retrieves uploaded files from a POST or PUT request.
//
// Parameters:
//   - files (...string): Optional list of form field names to process. If none are specified, all uploaded files are processed.
//
// Returns:
//   - []UploadedFile: A list of UploadedFile details for each processed file.
//   - error: An error if the file upload processing fails.
func (c *Ctx) FormUpload(files ...string) ([]UploadedFile, error) {
	var uploadedFiles []UploadedFile

	if len(files) > 0 {
		for _, file := range files {
			// Read file header
			header, err := c.root.FormFile(file)
			if err != nil {
				return uploadedFiles, err
			}

			// Create temporary file.
			tempName := fmt.Sprintf("%s.%s", utils.Token(), utils.FileExt(header.Filename))
			filePath := fmt.Sprintf("%s/%s", TempDir, tempName)

			// Save file
			err = fasthttp.SaveMultipartFile(header, filePath)
			if err != nil {
				return uploadedFiles, err
			}

			uploadedFiles = append(uploadedFiles, UploadedFile{
				Name:  header.Filename,
				Path:  filePath,
				Size:  header.Size,
				Field: file,
			})
		}
	} else {
		form, err := c.root.MultipartForm()
		if err != nil {
			return nil, err
		}

		for name, v := range form.File {
			for _, header := range v {
				// Create temporary file.
				tempName := fmt.Sprintf("%s.%s", utils.Token(), utils.FileExt(header.Filename))
				filePath := fmt.Sprintf("%s/%s", TempDir, tempName)

				err = fasthttp.SaveMultipartFile(header, filePath)
				if err != nil {
					return uploadedFiles, err
				}

				uploadedFiles = append(uploadedFiles, UploadedFile{
					Name:  header.Filename,
					Path:  filePath,
					Size:  header.Size,
					Field: name,
				})
			}
		}
	}

	return uploadedFiles, nil
}

// Queries retrieves all key-value pairs from the query string of the HTTP request.
//
// Returns:
//   - map[string]string: A map of key-value pairs representing the query string data.
func (c *Ctx) Queries() map[string]string {
	m := make(map[string]string, c.root.QueryArgs().Len())
	c.root.QueryArgs().VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})
	return m
}

// QueryStr retrieves a string value from the query string of the HTTP request.
//
// Parameters:
//   - key (string): The key to retrieve from the query string.
//
// Returns:
//   - string: The value associated with the provided key.
func (c *Ctx) QueryStr(key string) string {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return string(data)
}

// QueryInt retrieves an integer value from the query string of the HTTP request.
//
// Parameters:
//   - key (string): The key to retrieve from the query string.
//
// Returns:
//   - int: The integer value associated with the provided key.
//   - error: An error if the value cannot be converted to an integer.
func (c *Ctx) QueryInt(key string) (int, error) {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.Atoi(string(data))
}

// QueryBool retrieves a boolean value from the query string or form data of the HTTP request.
//
// Parameters:
//   - key (string): The key to retrieve the boolean value from.
//
// Returns:
//   - bool: The boolean value associated with the provided key.
//   - error: An error if the value cannot be converted to a boolean.
func (c *Ctx) QueryBool(key string) (bool, error) {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseBool(string(data))
}

// QueryFloat retrieves a float value from the query string or form data of the HTTP request.
//
// Parameters:
//   - key (string): The key to retrieve the float value from.
//
// Returns:
//   - float64: The float value associated with the provided key.
//   - error: An error if the value cannot be converted to a float.
func (c *Ctx) QueryFloat(key string) (float64, error) {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseFloat(string(data), 64)
}

// PathVal retrieves a value from the path parameters of the HTTP request.
//
// Parameters:
//   - key (string): The name of the path parameter to retrieve.
//
// Returns:
//   - string: The value associated with the provided key, or an empty string if the key is not found.
func (c *Ctx) PathVal(key string) string {
	val := c.root.UserValue(key)

	if val == nil {
		return ""
	}

	return val.(string)
}

// OriginalURL retrieves the original URL of the HTTP request.
//
// Returns:
//   - string: The original URL requested by the client.
func (c *Ctx) OriginalURL() string {
	return string(c.root.Request.Header.RequestURI())
}

// ====================================================================
//                          Ctx - Data
// ====================================================================

type IData interface {
	// SetData stores data in the request context Ctx.
	//
	// Parameters:
	//   - key (string): The key under which the data will be stored.
	//   - data (any): The data to be stored.
	SetData(key string, data any)

	// GetData retrieves data from the request context Ctx.
	//
	// Parameters:
	//   - key (string): The key associated with the data.
	//
	// Returns:
	//   - any: The data associated with the provided key.
	GetData(key string) any

	// SetSession stores data in the session context.
	//
	// Parameters:
	//   - key (string): The key under which the session data will be stored.
	//   - data (any): The session data to be stored.
	SetSession(key string, data any)

	// GetSession retrieves data from the session context.
	//
	// Parameters:
	//   - key (string): The key associated with the session data.
	//
	// Returns:
	//   - any: The session data associated with the provided key.
	GetSession(key string) any
}

// SetData stores data in the request context Ctx.
//
// Parameters:
//   - key (string): The key under which the data will be stored.
//   - data (any): The data to be stored.
func (c *Ctx) SetData(key string, data any) {
	c.data[key] = data
}

// GetData retrieves data from the request context Ctx.
//
// Parameters:
//   - key (string): The key associated with the data.
//
// Returns:
//   - any: The data associated with the provided key.
func (c *Ctx) GetData(key string) any {
	return c.data[key]
}

// SetSession stores data in the session context.
//
// Parameters:
//   - key (string): The key under which the session data will be stored.
//   - data (any): The session data to be stored.
func (c *Ctx) SetSession(key string, data any) {
	session.Set(c, key, data)
}

// GetSession retrieves data from the session context.
//
// Parameters:
//   - key (string): The key associated with the session data.
//
// Returns:
//   - any: The session data associated with the provided key.
func (c *Ctx) GetSession(key string) any {
	return session.Get(c, key)
}
