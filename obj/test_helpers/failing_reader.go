package test_helpers

import "errors"

type FailingReader struct {
}

func (r *FailingReader) Read(p []byte) (n int, err error) {
	return 1, errors.New("Failed to read!")
}
