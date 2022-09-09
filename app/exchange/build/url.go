package build

import (
	"net/url"
	"rest/app/input"

	"github.com/pkg/errors"
)

func buildURL(in *input.Input) (*url.URL, error) {
	q, err := url.ParseQuery(in.URL.RawQuery)
	if err != nil {
		return nil, errors.Wrap(err, "parsing query string")
	}
	for _, field := range in.Parameters {
		value, err := resolveFieldValue(field)
		if err != nil {
			return nil, err
		}
		q.Add(field.Name, value)
	}

	u := *in.URL
	u.RawQuery = q.Encode()
	return &u, nil
}
