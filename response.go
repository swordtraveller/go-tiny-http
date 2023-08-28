package http

type ResponseWriter interface {
	Write([]byte) (int, error)
}

type responseWriter struct {
	ResponseBody  string
	ContentLength int
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	rw.ResponseBody += string(p)
	rw.ContentLength = len(rw.ResponseBody)
	return len(p), nil
}
