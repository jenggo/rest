package vars

import (
	"net/http"
	"time"
)

type Options struct {
	// Transport is applied to the underlying HTTP client. Use to mock or
	// intercept network traffic.  If nil, http.DefaultTransport will be cloned.
	Transport http.RoundTripper
}

type ExchangeOptions struct {
	Timeout         time.Duration
	FollowRedirects bool
	Auth            AuthOptions
	SkipVerify      bool
	ForceHTTP1      bool
	CheckStatus     bool
	Transport       http.RoundTripper
}

type InputOptions struct {
	JSON      bool
	Form      bool
	ReadStdin bool
}

type OutputOptions struct {
	PrintRequestHeader  bool
	PrintRequestBody    bool
	PrintResponseHeader bool
	PrintResponseBody   bool

	EnableFormat bool
	EnableColor  bool

	Download   bool
	OutputFile string
	Overwrite  bool
}

type AuthOptions struct {
	Enabled  bool
	UserName string
	Password string
}
