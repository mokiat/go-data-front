package testutil

import "io"

func NewFailingReader(err error) io.Reader {
	return &failingReader{
		err: err,
	}
}

type failingReader struct {
	err error
}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 1, r.err
}
