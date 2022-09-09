package build

import (
	"net/http"
	"rest/app/input"
)

func buildHTTPHeader(in *input.Input) (http.Header, error) {
	header := make(http.Header)
	for _, field := range in.Header.Fields {
		value, err := resolveFieldValue(field)
		if err != nil {
			return nil, err
		}
		header.Add(field.Name, value)
	}
	return header, nil
}
