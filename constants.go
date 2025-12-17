package core

// ====================================================================
//                          gFly's Information
// ====================================================================

// Version of the current gFly application
const (
	// Version indicates the version of the framework
	Version = "v1.17.12"
)

// ====================================================================
//                               Network
// ====================================================================

// Network
const (
	// GlobalIpv4Addr represents the global IPv4 address used across the application
	GlobalIpv4Addr = "0.0.0.0"

	// SchemaHTTP specifies the HTTP schema
	SchemaHTTP = "http"

	// SchemaHTTPS specifies the HTTPS schema
	SchemaHTTPS = "https"

	// NetworkTCP represents the TCP network protocol
	NetworkTCP = "tcp"

	// NetworkTCP4 represents the IPv4-based TCP network protocol
	NetworkTCP4 = "tcp4"

	// NetworkTCP6 represents the IPv6-based TCP network protocol
	NetworkTCP6 = "tcp6"
)

// Compression types
const (
	// StrGzip represents gzip compression type
	StrGzip = "gzip"

	// StrBr represents Brotli compression type
	StrBr = "br"

	// StrDeflate represents deflate compression type
	StrDeflate = "deflate"

	// StrBrotli represents Brotli compression type
	StrBrotli = "brotli"
)

// Cookie SameSite
const (
	// CookieSameSiteDisabled indicates that the SameSite attribute is disabled in cookies
	CookieSameSiteDisabled = "disabled"

	// CookieSameSiteLaxMode indicates SameSite attribute set to Lax, which prevents cross-site usage
	// except for top-level navigations with safe HTTP methods
	CookieSameSiteLaxMode = "lax"

	// CookieSameSiteStrictMode indicates SameSite attribute set to Strict, restricting cross-site usage entirely
	CookieSameSiteStrictMode = "strict"

	// CookieSameSiteNoneMode indicates SameSite attribute set to None, allowing cross-site usage with Secure flag
	CookieSameSiteNoneMode = "none"
)

// ====================================================================
//                             MINE types
// ====================================================================

// MIME types that are commonly used
const (
	// MIMETextXML represents the MIME type for XML text.
	MIMETextXML = "text/xml"
	// MIMETextHTML represents the MIME type for HTML text.
	MIMETextHTML = "text/html"
	// MIMETextPlain represents the MIME type for plain text.
	MIMETextPlain = "text/plain"
	// MIMETextJavaScript represents the MIME type for JavaScript text.
	MIMETextJavaScript = "text/javascript"
	// MIMEApplicationXML represents the MIME type for application XML.
	MIMEApplicationXML = "application/xml"
	// MIMEApplicationJSON represents the MIME type for application JSON.
	MIMEApplicationJSON = "application/json"
	// Deprecated: MIMEApplicationJavaScript use MIMETextJavaScript instead.
	MIMEApplicationJavaScript = "application/javascript"
	// MIMEApplicationForm represents the MIME type for form-urlencoded data.
	MIMEApplicationForm = "application/x-www-form-urlencoded"
	// MIMEOctetStream represents the MIME type for binary streams.
	MIMEOctetStream = "application/octet-stream"
	// MIMEMultipartForm represents the MIME type for multipart form data.
	MIMEMultipartForm = "multipart/form-data"

	// MIMETextXMLCharsetUTF8 represents the MIME type for XML text encoded in UTF-8.
	MIMETextXMLCharsetUTF8 = "text/xml; charset=utf-8"
	// MIMETextHTMLCharsetUTF8 represents the MIME type for HTML text encoded in UTF-8.
	MIMETextHTMLCharsetUTF8 = "text/html; charset=utf-8"
	// MIMETextPlainCharsetUTF8 represents the MIME type for plain text encoded in UTF-8.
	MIMETextPlainCharsetUTF8 = "text/plain; charset=utf-8"
	// MIMETextJavaScriptCharsetUTF8 represents the MIME type for JavaScript text encoded in UTF-8.
	MIMETextJavaScriptCharsetUTF8 = "text/javascript; charset=utf-8"
	// MIMEApplicationXMLCharsetUTF8 represents the MIME type for application XML encoded in UTF-8.
	MIMEApplicationXMLCharsetUTF8 = "application/xml; charset=utf-8"
	// MIMEApplicationJSONCharsetUTF8 represents the MIME type for application JSON encoded in UTF-8.
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"
	// Deprecated: MIMEApplicationJavaScriptCharsetUTF8 use MIMETextJavaScriptCharsetUTF8 instead.
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

// ====================================================================
//                             HTTP methods
// ====================================================================

// HTTP methods were copied from net/http.
const (
	// MethodGet represents the GET HTTP method used to retrieve data from the server.
	// It is a safe, idempotent, and cacheable request method. (RFC 7231, 4.3.1)
	MethodGet = "GET"

	// MethodHead represents the HEAD HTTP method used to retrieve headers
	// without the response body. It is often used for metadata such as length. (RFC 7231, 4.3.2)
	MethodHead = "HEAD"

	// MethodPost represents the POST HTTP method used to send data to the server,
	// typically for creating or updating resources. (RFC 7231, 4.3.3)
	MethodPost = "POST"

	// MethodPut represents the PUT HTTP method used to update or create a resource
	// at the specified URI. It is idempotent. (RFC 7231, 4.3.4)
	MethodPut = "PUT"

	// MethodPatch represents the PATCH HTTP method used to apply partial modifications
	// to a resource. It is not necessarily idempotent. (RFC 5789)
	MethodPatch = "PATCH"

	// MethodDelete represents the DELETE HTTP method used to remove a resource
	// at the specified URI. It is idempotent. (RFC 7231, 4.3.5)
	MethodDelete = "DELETE"

	// MethodConnect represents the CONNECT HTTP method, typically used for
	// creating a network connection tunnel, often for HTTPS. (RFC 7231, 4.3.6)
	MethodConnect = "CONNECT"

	// MethodOptions represents the OPTIONS HTTP method, used to describe the
	// communication options for the target resource or server. (RFC 7231, 4.3.7)
	MethodOptions = "OPTIONS"

	// MethodTrace represents the TRACE HTTP method, used for diagnostic purposes.
	// It performs a message loop-back along the path to the target resource. (RFC 7231, 4.3.8)
	MethodTrace = "TRACE"

	// MethodUse represents a custom USE HTTP method.
	// Its purpose and implementation may vary based on application requirements.
	MethodUse = "USE"
)

// ====================================================================
//                             HTTP status
// ====================================================================

// HTTP status codes were copied from net/http with the following updates:
// - Rename StatusNonAuthoritativeInfo to StatusNonAuthoritativeInformation
// - Add StatusSwitchProxy (306)
// NOTE: Keep this list in sync with statusMessage
const (
	// StatusContinue indicates that the client should continue with its request. (RFC 9110, 15.2.1)
	StatusContinue = 100

	// StatusSwitchingProtocols indicates the server agrees to switch protocols. (RFC 9110, 15.2.2)
	StatusSwitchingProtocols = 101

	// StatusProcessing indicates that the request is being processed and no response is available yet. (RFC 2518, 10.1)
	StatusProcessing = 102

	// StatusEarlyHints provides hints to help the client start preloading resources while waiting for the final response. (RFC 8297)
	StatusEarlyHints = 103

	// StatusOK represents that the request has succeeded. (RFC 9110, 15.3.1)
	StatusOK = 200

	// StatusCreated indicates that the request has resulted in one or more new resources being created. (RFC 9110, 15.3.2)
	StatusCreated = 201

	// StatusAccepted indicates that the request has been received but not yet acted upon. (RFC 9110, 15.3.3)
	StatusAccepted = 202

	// StatusNonAuthoritativeInformation indicates that the returned metadata is not from the origin server. (RFC 9110, 15.3.4)
	StatusNonAuthoritativeInformation = 203

	// StatusNoContent indicates that the server successfully fulfilled the request, but there is no content to send in response. (RFC 9110, 15.3.5)
	StatusNoContent = 204

	// StatusResetContent indicates that the server successfully fulfilled the request and expects the user agent to reset the document view. (RFC 9110, 15.3.6)
	StatusResetContent = 205

	// StatusPartialContent indicates that the server successfully fulfilled a range request for the target resource. (RFC 9110, 15.3.7)
	StatusPartialContent = 206

	// StatusMultiStatus indicates multiple status codes for a single request involving multiple operations. (RFC 4918, 11.1)
	StatusMultiStatus = 207

	// StatusAlreadyReported indicates that the members of a DAV binding have already been reported in a previous reply. (RFC 5842, 7.1)
	StatusAlreadyReported = 208

	// StatusIMUsed indicates that the server fulfilled a GET request including instance-manipulations. (RFC 3229, 10.4.1)
	StatusIMUsed = 226

	// StatusMultipleChoices indicates multiple options for the requested resource. (RFC 9110, 15.4.1)
	StatusMultipleChoices = 300

	// StatusMovedPermanently indicates that the target resource has been assigned a new permanent URI. (RFC 9110, 15.4.2)
	StatusMovedPermanently = 301

	// StatusFound indicates the target resource resides temporarily under a different URI and may change in the future. (RFC 9110, 15.4.3)
	StatusFound = 302

	// StatusSeeOther indicates that the server redirects the client to a different URI for the requested resource. (RFC 9110, 15.4.4)
	StatusSeeOther = 303

	// StatusNotModified indicates that the cached version of the requested resource is still valid. (RFC 9110, 15.4.5)
	StatusNotModified = 304

	// StatusUseProxy is deprecated and was reserved for proxy usage. (RFC 9110, 15.4.6)
	StatusUseProxy = 305

	// StatusSwitchProxy was used to indicate a new proxy, but it is unused. (RFC 9110, 15.4.7)
	StatusSwitchProxy = 306

	// StatusTemporaryRedirect indicates that the target resource resides temporarily under a different URI. (RFC 9110, 15.4.8)
	StatusTemporaryRedirect = 307

	// StatusPermanentRedirect indicates that the target resource has been assigned a permanent URI. (RFC 9110, 15.4.9)
	StatusPermanentRedirect = 308

	// StatusBadRequest indicates that the server could not understand the request due to invalid syntax. (RFC 9110, 15.5.1)
	StatusBadRequest = 400

	// StatusUnauthorized indicates that the request requires user authentication. (RFC 9110, 15.5.2)
	StatusUnauthorized = 401

	// StatusPaymentRequired is reserved for future use. (RFC 9110, 15.5.3)
	StatusPaymentRequired = 402

	// StatusForbidden indicates that the server understands the request but refuses to authorize it. (RFC 9110, 15.5.4)
	StatusForbidden = 403

	// StatusNotFound indicates that the server cannot find the requested resource. (RFC 9110, 15.5.5)
	StatusNotFound = 404

	// StatusMethodNotAllowed indicates that the request method is not supported for the target resource. (RFC 9110, 15.5.6)
	StatusMethodNotAllowed = 405

	// StatusNotAcceptable indicates that no resource is available matching the criteria specified by the request headers. (RFC 9110, 15.5.7)
	StatusNotAcceptable = 406

	// StatusProxyAuthRequired indicates that the request must be authenticated with a proxy. (RFC 9110, 15.5.8)
	StatusProxyAuthRequired = 407

	// StatusRequestTimeout indicates that the server timed out waiting for the request. (RFC 9110, 15.5.9)
	StatusRequestTimeout = 408

	// StatusConflict indicates that the request conflicts with the current state of the resource. (RFC 9110, 15.5.10)
	StatusConflict = 409

	// StatusGone indicates that the target resource is no longer available at the server. (RFC 9110, 15.5.11)
	StatusGone = 410

	// StatusLengthRequired indicates that the server requires a Content-Length header field in the request. (RFC 9110, 15.5.12)
	StatusLengthRequired = 411

	// StatusPreconditionFailed indicates that a precondition in the request header fields evaluated to false. (RFC 9110, 15.5.13)
	StatusPreconditionFailed = 412

	// StatusRequestEntityTooLarge indicates that the request entity is larger than the server is willing or able to process. (RFC 9110, 15.5.14)
	StatusRequestEntityTooLarge = 413

	// StatusRequestURITooLong indicates that the request URI is longer than the server is willing to interpret. (RFC 9110, 15.5.15)
	StatusRequestURITooLong = 414

	// StatusUnsupportedMediaType indicates that the request entity's media type is unsupported by the server. (RFC 9110, 15.5.16)
	StatusUnsupportedMediaType = 415

	// StatusRequestedRangeNotSatisfiable indicates that the requested range cannot be fulfilled. (RFC 9110, 15.5.17)
	StatusRequestedRangeNotSatisfiable = 416

	// StatusExpectationFailed indicates that the server cannot meet the requirements of the Expect request-header field. (RFC 9110, 15.5.18)
	StatusExpectationFailed = 417

	// StatusTeapot is a humorous response code. (RFC 9110, 15.5.19)
	StatusTeapot = 418

	// StatusMisdirectedRequest indicates that the server is not able to produce a response due to inappropriate routing. (RFC 9110, 15.5.20)
	StatusMisdirectedRequest = 421

	// StatusUnprocessableEntity indicates that the server understands the content type but cannot process the request. (RFC 9110, 15.5.21)
	StatusUnprocessableEntity = 422

	// StatusLocked indicates that the resource is currently locked. (RFC 4918, 11.3)
	StatusLocked = 423

	// StatusFailedDependency indicates that the request failed due to a dependency on another request. (RFC 4918, 11.4)
	StatusFailedDependency = 424

	// StatusTooEarly indicates that the server is unwilling to risk processing a replayed request. (RFC 8470, 5.2.)
	StatusTooEarly = 425

	// StatusUpgradeRequired indicates that the client should switch to a different protocol. (RFC 9110, 15.5.22)
	StatusUpgradeRequired = 426

	// StatusPreconditionRequired indicates that the server requires the request to be conditional. (RFC 6585, 3)
	StatusPreconditionRequired = 428

	// StatusTooManyRequests indicates that the user has sent too many requests in a given amount of time. (RFC 6585, 4)
	StatusTooManyRequests = 429

	// StatusRequestHeaderFieldsTooLarge indicates that the server refuses to process the request because the header fields are too large. (RFC 6585, 5)
	StatusRequestHeaderFieldsTooLarge = 431

	// StatusUnavailableForLegalReasons indicates access to the resource is denied for legal reasons. (RFC 7725, 3)
	StatusUnavailableForLegalReasons = 451

	// StatusInternalServerError indicates that the server encountered an unexpected condition that prevented it from fulfilling the request. (RFC 9110, 15.6.1)
	StatusInternalServerError = 500

	// StatusNotImplemented indicates that the server does not support the functionality required to fulfill the request. (RFC 9110, 15.6.2)
	StatusNotImplemented = 501

	// StatusBadGateway indicates that the server, acting as a gateway or proxy, received an invalid response from an upstream server. (RFC 9110, 15.6.3)
	StatusBadGateway = 502

	// StatusServiceUnavailable indicates that the server is temporarily unable to handle the request due to maintenance or overload. (RFC 9110, 15.6.4)
	StatusServiceUnavailable = 503

	// StatusGatewayTimeout indicates that the server, acting as a gateway or proxy, did not receive a timely response from an upstream server. (RFC 9110, 15.6.5)
	StatusGatewayTimeout = 504

	// StatusHTTPVersionNotSupported indicates that the server does not support the HTTP protocol version used in the request. (RFC 9110, 15.6.6)
	StatusHTTPVersionNotSupported = 505

	// StatusVariantAlsoNegotiates indicates that the server detects loop negotiation. (RFC 2295, 8.1)
	StatusVariantAlsoNegotiates = 506

	// StatusInsufficientStorage indicates that the server is unable to store the representation needed to complete the request. (RFC 4918, 11.5)
	StatusInsufficientStorage = 507

	// StatusLoopDetected indicates that the server detected an infinite loop while processing the request. (RFC 5842, 7.2)
	StatusLoopDetected = 508

	// StatusNotExtended indicates that further extensions to the request are required for the server to fulfill it. (RFC 2774, 7)
	StatusNotExtended = 510

	// StatusNetworkAuthenticationRequired indicates that the client needs to authenticate to gain network access. (RFC 6585, 6)
	StatusNetworkAuthenticationRequired = 511
)

// ====================================================================
//                             HTTP headers
// ====================================================================

// HTTP Headers were copied from net/http.
const (
	// HeaderAuthorization represents the HTTP Authorization header.
	HeaderAuthorization = "Authorization"
	// HeaderProxyAuthenticate represents the HTTP Proxy-Authenticate header.
	HeaderProxyAuthenticate = "Proxy-Authenticate"
	// HeaderProxyAuthorization represents the HTTP Proxy-Authorization header.
	HeaderProxyAuthorization = "Proxy-Authorization"
	// HeaderWWWAuthenticate represents the HTTP WWW-Authenticate header.
	HeaderWWWAuthenticate = "WWW-Authenticate"
	// HeaderAge represents the HTTP Age header.
	HeaderAge = "Age"
	// HeaderCacheControl represents the HTTP Cache-Control header.
	HeaderCacheControl = "Cache-Control"
	// HeaderClearSiteData represents the HTTP Clear-Site-Data header.
	HeaderClearSiteData = "Clear-Site-Data"
	// HeaderExpires represents the HTTP Expires header.
	HeaderExpires = "Expires"
	// HeaderPragma represents the HTTP Pragma header.
	HeaderPragma = "Pragma"
	// HeaderWarning represents the HTTP Warning header.
	HeaderWarning = "Warning"
	// HeaderAcceptCH represents the HTTP Accept-CH header.
	HeaderAcceptCH = "Accept-CH"
	// HeaderAcceptCHLifetime represents the HTTP Accept-CH-Lifetime header.
	HeaderAcceptCHLifetime = "Accept-CH-Lifetime"
	// HeaderContentDPR represents the HTTP Content-DPR header.
	HeaderContentDPR = "Content-DPR"
	// HeaderDPR represents the HTTP DPR header.
	HeaderDPR = "DPR"
	// HeaderEarlyData represents the HTTP Early-Data header.
	HeaderEarlyData = "Early-Data"
	// HeaderSaveData represents the HTTP Save-Data header.
	HeaderSaveData = "Save-Data"
	// HeaderViewportWidth represents the HTTP Viewport-Width header.
	HeaderViewportWidth = "Viewport-Width"
	// HeaderWidth represents the HTTP Width header.
	HeaderWidth = "Width"
	// HeaderETag represents the HTTP ETag header.
	HeaderETag = "ETag"
	// HeaderIfMatch represents the HTTP If-Match header.
	HeaderIfMatch = "If-Match"
	// HeaderIfModifiedSince represents the HTTP If-Modified-Since header.
	HeaderIfModifiedSince = "If-Modified-Since"
	// HeaderIfNoneMatch represents the HTTP If-None-Match header.
	HeaderIfNoneMatch = "If-None-Match"
	// HeaderIfUnmodifiedSince represents the HTTP If-Unmodified-Since header.
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	// HeaderLastModified represents the HTTP Last-Modified header.
	HeaderLastModified = "Last-Modified"
	// HeaderVary represents the HTTP Vary header.
	HeaderVary = "Vary"
	// HeaderConnection represents the HTTP Connection header.
	HeaderConnection = "Connection"
	// HeaderKeepAlive represents the HTTP Keep-Alive header.
	HeaderKeepAlive = "Keep-Alive"
	// HeaderAccept represents the HTTP Accept header.
	HeaderAccept = "Accept"
	// HeaderAcceptCharset represents the HTTP Accept-Charset header.
	HeaderAcceptCharset = "Accept-Charset"
	// HeaderAcceptEncoding represents the HTTP Accept-Encoding header.
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderAcceptLanguage represents the HTTP Accept-Language header.
	HeaderAcceptLanguage = "Accept-Language"
	// HeaderCookie represents the HTTP Cookie header.
	HeaderCookie = "Cookie"
	// HeaderExpect represents the HTTP Expect header.
	HeaderExpect = "Expect"
	// HeaderMaxForwards represents the HTTP Max-Forwards header.
	HeaderMaxForwards = "Max-Forwards"
	// HeaderSetCookie represents the HTTP Set-Cookie header.
	HeaderSetCookie = "Set-Cookie"
	// HeaderAccessControlAllowCredentials represents the HTTP Access-Control-Allow-Credentials header.
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// HeaderAccessControlAllowHeaders represents the HTTP Access-Control-Allow-Headers header.
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// HeaderAccessControlAllowMethods represents the HTTP Access-Control-Allow-Methods header.
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	// HeaderAccessControlAllowOrigin represents the HTTP Access-Control-Allow-Origin header.
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderAccessControlExposeHeaders represents the HTTP Access-Control-Expose-Headers header.
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// HeaderAccessControlMaxAge represents the HTTP Access-Control-Max-Age header.
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"
	// HeaderAccessControlRequestHeaders represents the HTTP Access-Control-Request-Headers header.
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
	// HeaderAccessControlRequestMethod represents the HTTP Access-Control-Request-Method header.
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"
	// HeaderOrigin represents the HTTP Origin header.
	HeaderOrigin = "Origin"
	// HeaderTimingAllowOrigin represents the HTTP Timing-Allow-Origin header.
	HeaderTimingAllowOrigin = "Timing-Allow-Origin"
	// HeaderXPermittedCrossDomainPolicies represents the HTTP X-Permitted-Cross-Domain-Policies header.
	HeaderXPermittedCrossDomainPolicies = "X-Permitted-Cross-Domain-Policies"
	// HeaderDNT represents the HTTP DNT (Do Not Track) header.
	HeaderDNT = "DNT"
	// HeaderTk represents the HTTP Tk (Tracking Status Value) header.
	HeaderTk = "Tk"
	// HeaderContentDisposition represents the HTTP Content-Disposition header.
	HeaderContentDisposition = "Content-Disposition"
	// HeaderContentEncoding represents the HTTP Content-Encoding header.
	HeaderContentEncoding = "Content-Encoding"
	// HeaderContentLanguage represents the HTTP Content-Language header.
	HeaderContentLanguage = "Content-Language"
	// HeaderContentLength represents the HTTP Content-Length header.
	HeaderContentLength = "Content-Length"
	// HeaderContentLocation represents the HTTP Content-Location header.
	HeaderContentLocation = "Content-Location"
	// HeaderContentType represents the HTTP Content-Type header.
	HeaderContentType = "Content-Type"
	// HeaderForwarded represents the HTTP Forwarded header.
	HeaderForwarded = "Forwarded"
	// HeaderVia represents the HTTP Via header.
	HeaderVia = "Via"
	// HeaderXForwardedFor represents the HTTP X-Forwarded-For header.
	HeaderXForwardedFor = "X-Forwarded-For"
	// HeaderXForwardedHost represents the HTTP X-Forwarded-Host header.
	HeaderXForwardedHost = "X-Forwarded-Host"
	// HeaderXForwardedProto represents the HTTP X-Forwarded-Proto header.
	HeaderXForwardedProto = "X-Forwarded-Proto"
	// HeaderXForwardedProtocol represents the HTTP X-Forwarded-Protocol header.
	HeaderXForwardedProtocol = "X-Forwarded-Protocol"
	// HeaderXForwardedSsl represents the HTTP X-Forwarded-Ssl header.
	HeaderXForwardedSsl = "X-Forwarded-Ssl"
	// HeaderXUrlScheme represents the HTTP X-Url-Scheme header.
	HeaderXUrlScheme = "X-Url-Scheme"
	// HeaderLocation represents the HTTP Location header.
	HeaderLocation = "Location"
	// HeaderFrom represents the HTTP From header.
	HeaderFrom = "From"
	// HeaderHost represents the HTTP Host header.
	HeaderHost = "Host"
	// HeaderReferer represents the HTTP Referer header.
	HeaderReferer = "Referer"
	// HeaderReferrerPolicy represents the HTTP Referrer-Policy header.
	HeaderReferrerPolicy = "Referrer-Policy"
	// HeaderUserAgent represents the HTTP User-Agent header.
	HeaderUserAgent = "User-Agent"
	// HeaderAllow represents the HTTP Allow header.
	HeaderAllow = "Allow"
	// HeaderServer represents the HTTP Server header.
	HeaderServer = "Server"
	// HeaderAcceptRanges represents the HTTP Accept-Ranges header.
	HeaderAcceptRanges = "Accept-Ranges"
	// HeaderContentRange represents the HTTP Content-Range header.
	HeaderContentRange = "Content-Range"
	// HeaderIfRange represents the HTTP If-Range header.
	HeaderIfRange = "If-Range"
	// HeaderRange represents the HTTP Range header.
	HeaderRange = "Range"
	// HeaderContentSecurityPolicy represents the HTTP Content-Security-Policy header.
	HeaderContentSecurityPolicy = "Content-Security-Policy"
	// HeaderContentSecurityPolicyReportOnly represents the HTTP Content-Security-Policy-Report-Only header.
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	// HeaderCrossOriginResourcePolicy represents the HTTP Cross-Origin-Resource-Policy header.
	HeaderCrossOriginResourcePolicy = "Cross-Origin-Resource-Policy"
	// HeaderExpectCT represents the HTTP Expect-CT header.
	HeaderExpectCT = "Expect-CT"
	// HeaderFeaturePolicy represents the deprecated HTTP Feature-Policy header.
	// Deprecated: use HeaderPermissionsPolicy instead.
	HeaderFeaturePolicy = "Feature-Policy"

	// HeaderPermissionsPolicy represents the HTTP Permissions-Policy header.
	HeaderPermissionsPolicy = "Permissions-Policy"

	// HeaderPublicKeyPins represents the HTTP Public-Key-Pins header.
	HeaderPublicKeyPins = "Public-Key-Pins"

	// HeaderPublicKeyPinsReportOnly represents the HTTP Public-Key-Pins-Report-Only header.
	HeaderPublicKeyPinsReportOnly = "Public-Key-Pins-Report-Only"

	// HeaderStrictTransportSecurity represents the HTTP Strict-Transport-Security header.
	HeaderStrictTransportSecurity = "Strict-Transport-Security"

	// HeaderUpgradeInsecureRequests represents the HTTP Upgrade-Insecure-Requests header.
	HeaderUpgradeInsecureRequests = "Upgrade-Insecure-Requests"

	// HeaderXContentTypeOptions represents the HTTP X-Content-Type-Options header.
	HeaderXContentTypeOptions = "X-Content-Type-Options"

	// HeaderXDownloadOptions represents the HTTP X-Download-Options header.
	HeaderXDownloadOptions = "X-Download-Options"

	// HeaderXFrameOptions represents the HTTP X-Frame-Options header.
	HeaderXFrameOptions = "X-Frame-Options"

	// HeaderXPoweredBy represents the HTTP X-Powered-By header.
	HeaderXPoweredBy = "X-Powered-By"

	// HeaderXXSSProtection represents the HTTP X-XSS-Protection header.
	HeaderXXSSProtection = "X-XSS-Protection"

	// HeaderLastEventID represents the HTTP Last-Event-ID header.
	HeaderLastEventID = "Last-Event-ID"

	// HeaderNEL represents the HTTP NEL (Network Error Logging) header.
	HeaderNEL = "NEL"

	// HeaderPingFrom represents the HTTP Ping-From header.
	HeaderPingFrom = "Ping-From"

	// HeaderPingTo represents the HTTP Ping-To header.
	HeaderPingTo = "Ping-To"

	// HeaderReportTo represents the HTTP Report-To header.
	HeaderReportTo = "Report-To"

	// HeaderTE represents the HTTP TE (Transfer Encoding) header.
	HeaderTE = "TE"

	// HeaderTrailer represents the HTTP Trailer header.
	HeaderTrailer = "Trailer"

	// HeaderTransferEncoding represents the HTTP Transfer-Encoding header.
	HeaderTransferEncoding = "Transfer-Encoding"

	// HeaderSecWebSocketAccept represents the HTTP Sec-WebSocket-Accept header.
	HeaderSecWebSocketAccept = "Sec-WebSocket-Accept"

	// HeaderSecWebSocketExtensions represents the HTTP Sec-WebSocket-Extensions header.
	HeaderSecWebSocketExtensions = "Sec-WebSocket-Extensions"

	// HeaderSecWebSocketKey represents the HTTP Sec-WebSocket-Key header.
	HeaderSecWebSocketKey = "Sec-WebSocket-Key"

	// HeaderSecWebSocketProtocol represents the HTTP Sec-WebSocket-Protocol header.
	HeaderSecWebSocketProtocol = "Sec-WebSocket-Protocol"

	// HeaderSecWebSocketVersion represents the HTTP Sec-WebSocket-Version header.
	HeaderSecWebSocketVersion = "Sec-WebSocket-Version"

	// HeaderAcceptPatch represents the HTTP Accept-Patch header.
	HeaderAcceptPatch = "Accept-Patch"

	// HeaderAcceptPushPolicy represents the HTTP Accept-Push-Policy header.
	HeaderAcceptPushPolicy = "Accept-Push-Policy"

	// HeaderAcceptSignature represents the HTTP Accept-Signature header.
	HeaderAcceptSignature = "Accept-Signature"

	// HeaderAltSvc represents the HTTP Alt-Svc header.
	HeaderAltSvc = "Alt-Svc"

	// HeaderDate represents the HTTP Date header.
	HeaderDate = "Date"

	// HeaderIndex represents the HTTP Index header.
	HeaderIndex = "Index"

	// HeaderLargeAllocation represents the HTTP Large-Allocation header.
	HeaderLargeAllocation = "Large-Allocation"

	// HeaderLink represents the HTTP Link header.
	HeaderLink = "Link"

	// HeaderPushPolicy represents the HTTP Push-Policy header.
	HeaderPushPolicy = "Push-Policy"

	// HeaderRetryAfter represents the HTTP Retry-After header.
	HeaderRetryAfter = "Retry-After"

	// HeaderServerTiming represents the HTTP Server-Timing header.
	HeaderServerTiming = "Server-Timing"

	// HeaderSignature represents the HTTP Signature header.
	HeaderSignature = "Signature"

	// HeaderSignedHeaders represents the HTTP Signed-Headers header.
	HeaderSignedHeaders = "Signed-Headers"

	// HeaderSourceMap represents the HTTP SourceMap header.
	HeaderSourceMap = "SourceMap"

	// HeaderUpgrade represents the HTTP Upgrade header.
	HeaderUpgrade = "Upgrade"

	// HeaderXDNSPrefetchControl represents the HTTP X-DNS-Prefetch-Control header.
	HeaderXDNSPrefetchControl = "X-DNS-Prefetch-Control"

	// HeaderXPingBack represents the HTTP X-Pingback header.
	HeaderXPingBack = "X-Pingback"

	// HeaderXRequestID represents the HTTP X-Request-ID header.
	HeaderXRequestID = "X-Request-ID"

	// HeaderXRequestedWith represents the HTTP X-Requested-With header.
	HeaderXRequestedWith = "X-Requested-With"

	// HeaderXRobotsTag represents the HTTP X-Robots-Tag header.
	HeaderXRobotsTag = "X-Robots-Tag"

	// HeaderXUACompatible represents the HTTP X-UA-Compatible header.
	HeaderXUACompatible = "X-UA-Compatible"
)

// ====================================================================
//                          Parameter Types
// ====================================================================

// Parameter types used in request parameter operations
const (
	// ParamTypePost represents the form parameter type
	ParamTypePost = "post"

	// ParamTypeQuery represents the query parameter type
	ParamTypeQuery = "query"

	// ParamTypeBody represents the JSON body parameter type
	ParamTypeBody = "body"
)
