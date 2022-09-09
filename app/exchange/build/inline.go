package build

import (
	"mime/multipart"
	"net/textproto"

	"rest/app/input"
)

func buildInlinePart(field input.Field, multipartWriter *multipart.Writer) error {
	value, err := resolveFieldValue(field)
	if err != nil {
		return err
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", buildContentDisposition(field.Name, ""))
	w, err := multipartWriter.CreatePart(h)
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(value)); err != nil {
		return err
	}
	return nil
}
