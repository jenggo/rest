package build

import (
	"bytes"
	"fmt"
	"net/url"
)

func buildContentDisposition(name string, filename string) string {
	var buffer bytes.Buffer
	buffer.WriteString("form-data")

	if name != "" {
		if needEscape(name) {
			fmt.Fprintf(&buffer, `; name*=utf-8''%s`, url.PathEscape(name))
		} else {
			fmt.Fprintf(&buffer, `; name="%s"`, name)
		}
	}

	if filename != "" {
		if needEscape(filename) {
			fmt.Fprintf(&buffer, `; filename*=utf-8''%s`, url.PathEscape(filename))
		} else {
			fmt.Fprintf(&buffer, `; filename="%s"`, filename)
		}
	}

	return buffer.String()
}
