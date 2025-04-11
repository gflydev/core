package core

import "time"

type Config struct {
	// AppName specifies the application name.
	//
	// Default: "Laravel inspired web framework written in Go".
	AppName string

	// Name specifies the server name for sending in response headers.
	//
	// Default: "gFly".
	Name string

	// Concurrency defines the maximum number of concurrent connections the server may serve.
	// The maximum number of concurrent connections the server may serve.
	//
	// DefaultConcurrency is used if not set.
	//
	// Concurrency only works if you either call Serve once, or only ServeConn multiple times.
	// It works with ListenAndServe as well.
	// Default: 256 * 1024.
	Concurrency int

	// ReadTimeout specifies the maximum duration allowed to read the full request, including the body.
	// The connection's read deadline is reset when a connection opens or for keep-alive connections,
	// after the first byte has been read.
	//
	// Default: 60 minutes.
	ReadTimeout time.Duration

	// WriteTimeout specifies the maximum duration before timing out writes of the response.
	// It is reset after the request handler has returned.
	//
	// Default: 60 minutes.
	WriteTimeout time.Duration

	// IdleTimeout defines the maximum duration to wait for the next request when keep-alive is enabled.
	// If IdleTimeout is zero, the value of ReadTimeout is used.
	//
	// Default: 60 minutes.
	IdleTimeout time.Duration

	// ReadBufferSize specifies the buffer size for request reading, which also limits the maximum header size.
	// Increase this buffer if clients send large RequestURIs or headers.
	//
	// Default: 40,960 bytes.
	ReadBufferSize int

	// WriteBufferSize specifies the buffer size for response writing.
	//
	// Default: 40,960 bytes.
	WriteBufferSize int

	// NoDefaultDate determines whether the default "Date" header should be excluded from the response.
	//
	// Default: false.
	NoDefaultDate bool

	// NoDefaultContentType determines whether the default "Content-Type" header should be excluded from the response.
	//
	// Default: false.
	NoDefaultContentType bool

	// DisableHeaderNamesNormalizing determines whether header names should be passed as-is without normalization.
	// By default, header names are normalized.
	//
	// Default: false.
	DisableHeaderNamesNormalizing bool

	// DisableKeepalive indicates whether to disable keep-alive connections.
	//
	// Default: false.
	DisableKeepalive bool

	// MaxRequestBodySize specifies the maximum size of the request body.
	// A zero value means the default value will be honored.
	//
	// Default: 4 MB.
	MaxRequestBodySize int

	// NoDefaultServerHeader determines whether the default "Server" header should be excluded from the response.
	//
	// Default: false.
	NoDefaultServerHeader bool

	// GetOnly rejects all non-GET requests if set to true. Useful as anti-DoS protection for servers
	// that only accept GET and HEAD requests. Request size is limited by ReadBufferSize when enabled.
	//
	// Default: false.
	GetOnly bool

	// ReduceMemoryUsage aggressively reduces memory usage at the cost of higher CPU usage.
	// Only enable if serving mostly idle keep-alive connections.
	//
	// Default: false.
	ReduceMemoryUsage bool

	// StreamRequestBody enables request body streaming and calls the handler sooner
	// when the body exceeds the current limit.
	//
	// Default: true.
	StreamRequestBody bool

	// DisablePreParseMultipartForm determines whether to avoid pre-parsing multipart form data.
	// Useful for treating multipart form data as binary or controlling when data is parsed.
	//
	// Default: true.
	DisablePreParseMultipartForm bool

	// DisableStartupMessage determines whether to suppress the printing of "gFly" ASCII art and listening address
	// during startup.
	//
	// Default: false.
	DisableStartupMessage bool

	// Network specifies the network type ("tcp", "tcp4", "tcp6").
	// When prefork is true, only "tcp4" and "tcp6" are allowed.
	//
	// Default: "tcp4".
	Network string

	// Prefork enables preforking with multiple Go processes listening on the same port.
	//
	// Default: false.
	Prefork bool

	// CompressedFileSuffix specifies a suffix to the file name when saving the resulting compressed file.
	//
	// Default: ".gz".
	CompressedFileSuffix string
}

// DefaultConfig holds the default configuration settings for the gFly framework.
//
// Fields:
//   - AppName (string): Application name. Default: "Laravel inspired web framework written in Go".
//   - Name (string): Server name for sending in response headers. Default: "gFly".
//   - Concurrency (int): Maximum number of concurrent connections. Default: 256 * 1024.
//   - ReadTimeout (time.Duration): Request read timeout. Default: 60 minutes.
//   - WriteTimeout (time.Duration): Response write timeout. Default: 60 minutes.
//   - IdleTimeout (time.Duration): Maximum time to wait for the next request when keep-alive is enabled. Default: 60 minutes.
//   - ReadBufferSize (int): Buffer size for request reading. Default: 40,960 bytes.
//   - WriteBufferSize (int): Buffer size for response writing. Default: 40,960 bytes.
//   - NoDefaultDate (bool): Exclude default "Date" header from the response. Default: false.
//   - NoDefaultContentType (bool): Exclude default "Content-Type" header from the response. Default: false.
//   - DisableHeaderNamesNormalizing (bool): Disable normalization of header names. Default: false.
//   - DisableKeepalive (bool): Disable keep-alive connections. Default: false.
//   - MaxRequestBodySize (int): Maximum request body size. Default: 4 MB.
//   - NoDefaultServerHeader (bool): Exclude default "Server" header from the response. Default: false.
//   - GetOnly (bool): Reject all non-GET requests. Default: false.
//   - ReduceMemoryUsage (bool): Reduce memory usage aggressively. Default: false.
//   - StreamRequestBody (bool): Enable request body streaming for large requests. Default: true.
//   - DisablePreParseMultipartForm (bool): Do not pre-parse multipart form data. Default: true.
//   - DisableStartupMessage (bool): Do not print startup messages. Default: false.
//   - Network (string): Network type ("tcp", "tcp4", "tcp6"). Default: "tcp4".
//   - Prefork (bool): Enable prefork with multiple Go processes. Default: false.
//   - CompressedFileSuffix (string): Suffix for compressed file names. Default: ".gz".
var DefaultConfig = Config{
	AppName:                       "Laravel inspired web framework written in Go",
	Name:                          "gFly",
	Concurrency:                   256 * 1024,
	ReadTimeout:                   60 * 60 * time.Second, // unlimited
	WriteTimeout:                  60 * 60 * time.Second, // unlimited
	IdleTimeout:                   60 * 60 * time.Second, // unlimited
	ReadBufferSize:                4096 * 10,
	WriteBufferSize:               4096 * 10,
	NoDefaultDate:                 false,
	NoDefaultContentType:          false,
	DisableHeaderNamesNormalizing: false,
	DisableKeepalive:              false,
	MaxRequestBodySize:            4 * 1024 * 1024,
	NoDefaultServerHeader:         false, // True when `Name` Empty
	GetOnly:                       false,
	ReduceMemoryUsage:             false,
	StreamRequestBody:             true,
	DisablePreParseMultipartForm:  true,
	DisableStartupMessage:         false,
	Network:                       NetworkTCP4,
	Prefork:                       false,
	CompressedFileSuffix:          ".gz",
}
