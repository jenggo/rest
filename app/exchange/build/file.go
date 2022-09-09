package build

import (
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path"
	"rest/app/input"
	"strings"

	"github.com/pkg/errors"
)

func buildFilePart(field input.Field, multipartWriter *multipart.Writer) error {
	h := make(textproto.MIMEHeader)

	var filename string
	if field.IsFile {
		filename = path.Base(field.Value)
	}
	h.Set("Content-Disposition", buildContentDisposition(field.Name, filename))

	w, err := multipartWriter.CreatePart(h)
	if err != nil {
		return err
	}

	if field.IsFile {
		file, err := os.Open(field.Value)
		if err != nil {
			return errors.Wrapf(err, "failed to open '%s'", field.Value)
		}
		defer file.Close()

		if _, err := io.Copy(w, file); err != nil {
			return errors.Wrapf(err, "failed to read from '%s'", field.Value)
		}
	} else {
		if _, err := io.Copy(w, strings.NewReader(field.Value)); err != nil {
			return errors.Wrap(err, "failed to write to multipart writer")
		}
	}
	return nil
}
