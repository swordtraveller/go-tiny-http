package http

import (
	"io"
)

type Request struct {
	Method  string
	URL     string
	Proto   string
	Headers map[string]string
	Body    ReadCloser
}

type ReadCloser struct {
	value    []byte
	progress int
}

func (rc ReadCloser) Read(p []byte) (n int, err error) {
	curProgress := 0
	for index, _ := range p {
		if rc.progress == len(rc.value) {
			return curProgress, io.EOF
		} else {
			p[index] = rc.value[rc.progress]
			rc.progress++
			curProgress++
		}
	}
	return curProgress, nil
}
