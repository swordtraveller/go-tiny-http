package http

type ResponseWriter interface {
	Write([]byte) (int, error)
	SetStatusCode(code int)
}

type responseWriter struct {
	ResponseBody  string
	ContentLength int
	StatusCode    int
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	rw.ResponseBody += string(p)
	rw.ContentLength = len(rw.ResponseBody)
	return len(p), nil
}

func (rw *responseWriter) SetStatusCode(code int) {
	rw.StatusCode = code
}
