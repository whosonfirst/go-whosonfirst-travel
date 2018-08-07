package writer

import (
	"github.com/whosonfirst/go-whosonfirst-readwrite/utils"
	"io"
	"io/ioutil"
)

type MultiWriter struct {
	Writer
	writers []Writer
}

func NewMultiWriter(writers ...Writer) (Writer, error) {

	w := MultiWriter{
		writers: writers,
	}

	return &w, nil
}

func (w *MultiWriter) Write(path string, fh io.ReadCloser) error {

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	// please make this concurrent with a cancel context

	for _, wr := range w.writers {

		reader, err := utils.ReadCloserFromBytes(body)

		if err != nil {
			return err
		}

		err = wr.Write(path, reader)

		if err != nil {
			return err
		}
	}

	return nil
}
