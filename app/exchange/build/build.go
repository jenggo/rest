package build

import (
	"fmt"
	"io"
	"net/http"

	"rest/app/input"
	"rest/vars"
)

type bodyTuple struct {
	body          io.ReadCloser
	getBody       func() (io.ReadCloser, error)
	contentLength int64
	contentType   string
}

func HTTPRequest(in *input.Input, options *vars.ExchangeOptions) (*http.Request, error) {
	u, err := buildURL(in)
	if err != nil {
		return nil, err
	}

	header, err := buildHTTPHeader(in)
	if err != nil {
		return nil, err
	}

	bodyTuple, err := buildHTTPBody(in)
	if err != nil {
		return nil, err
	}

	if header.Get("Content-Type") == "" && bodyTuple.contentType != "" {
		header.Set("Content-Type", bodyTuple.contentType)
	}
	if header.Get("User-Agent") == "" {
		header.Set("User-Agent", fmt.Sprintf("rest/%s", vars.Current()))
	}

	r := http.Request{
		Method:        string(in.Method),
		URL:           u,
		Header:        header,
		Host:          header.Get("Host"),
		Body:          bodyTuple.body,
		GetBody:       bodyTuple.getBody,
		ContentLength: bodyTuple.contentLength,
	}

	if options.Auth.Enabled {
		r.SetBasicAuth(options.Auth.UserName, options.Auth.Password)
	}

	return &r, nil
}
