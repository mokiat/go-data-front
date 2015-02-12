package common_test_help

import "io"

type failingReader struct {
	err error
}

func NewFailingReader(err error) io.Reader {
	return &failingReader{
		err: err,
	}
}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 1, r.err
}
