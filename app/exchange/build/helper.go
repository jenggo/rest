package build

import (
	"io/ioutil"
	"rest/app/input"

	"github.com/pkg/errors"
)

func needEscape(s string) bool {
	for _, c := range s {
		if c > 127 {
			return true
		}
		if c < 32 && c != '\t' {
			return true
		}
		if c == '"' || c == '\\' {
			return true
		}
	}
	return false
}

func resolveFieldValue(field input.Field) (string, error) {
	if field.IsFile {
		data, err := ioutil.ReadFile(field.Value)
		if err != nil {
			return "", errors.Wrapf(err, "reading field value of '%s'", field.Name)
		}
		return string(data), nil
	} else {
		return field.Value, nil
	}
}
