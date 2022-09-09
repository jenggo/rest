package build

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"strings"

	"rest/app/input"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func buildHTTPBody(in *input.Input) (bodyTuple, error) {
	switch in.Body.BodyType {
	case input.EmptyBody:
		return bodyTuple{}, nil
	case input.JSONBody:
		return buildJSONBody(in)
	case input.FormBody:
		return buildFormBody(in)
	case input.RawBody:
		return buildRawBody(in)
	default:
		return bodyTuple{}, errors.Errorf("unknown body type: %v", in.Body.BodyType)
	}
}

func buildJSONBody(in *input.Input) (bodyTuple, error) {
	obj := map[string]interface{}{}
	for _, field := range in.Body.Fields {
		value, err := resolveFieldValue(field)
		if err != nil {
			return bodyTuple{}, err
		}
		obj[field.Name] = value
	}
	for _, field := range in.Body.RawJSONFields {
		value, err := resolveFieldValue(field)
		if err != nil {
			return bodyTuple{}, err
		}
		var v interface{}
		if err := json.Unmarshal([]byte(value), &v); err != nil {
			return bodyTuple{}, errors.Wrapf(err, "parsing JSON value of '%s'", field.Name)
		}
		obj[field.Name] = v
	}
	body, err := json.Marshal(obj)
	if err != nil {
		return bodyTuple{}, errors.Wrap(err, "marshaling JSON of HTTP body")
	}
	return bodyTuple{
		body: ioutil.NopCloser(bytes.NewReader(body)),
		getBody: func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(body)), nil
		},
		contentLength: int64(len(body)),
		contentType:   "application/json",
	}, nil
}

func buildFormBody(in *input.Input) (bodyTuple, error) {
	if len(in.Body.Files) > 0 {
		return buildMultipartBody(in)
	} else {
		return buildURLEncodedBody(in)
	}
}

func buildURLEncodedBody(in *input.Input) (bodyTuple, error) {
	form := url.Values{}
	for _, field := range in.Body.Fields {
		value, err := resolveFieldValue(field)
		if err != nil {
			return bodyTuple{}, err
		}
		form.Add(field.Name, value)
	}
	body := form.Encode()
	return bodyTuple{
		body: ioutil.NopCloser(strings.NewReader(body)),
		getBody: func() (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(body)), nil
		},
		contentLength: int64(len(body)),
		contentType:   "application/x-www-form-urlencoded; charset=utf-8",
	}, nil
}

func buildMultipartBody(in *input.Input) (bodyTuple, error) {
	var buffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&buffer)

	for _, field := range in.Body.Fields {
		if err := buildInlinePart(field, multipartWriter); err != nil {
			return bodyTuple{}, err
		}
	}
	for _, field := range in.Body.Files {
		if err := buildFilePart(field, multipartWriter); err != nil {
			return bodyTuple{}, err
		}
	}

	multipartWriter.Close()

	body := buffer.Bytes()
	return bodyTuple{
		body: ioutil.NopCloser(bytes.NewReader(body)),
		getBody: func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(body)), nil
		},
		contentLength: int64(len(body)),
		contentType:   multipartWriter.FormDataContentType(),
	}, nil
}

func buildRawBody(in *input.Input) (bodyTuple, error) {
	return bodyTuple{
		body: ioutil.NopCloser(bytes.NewReader(in.Body.Raw)),
		getBody: func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(in.Body.Raw)), nil
		},
		contentLength: int64(len(in.Body.Raw)),
		contentType:   "application/json",
	}, nil
}
