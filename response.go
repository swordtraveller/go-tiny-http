package http

type ResponseWriter interface {
	Write([]byte) (int, error)
	SetStatusCode(code int)
	Header() Header
	WriteHeader(statusCode int)
}

type responseWriter struct {
	ResponseBody  string
	ContentLength int
	StatusCode    int
	handlerHeader Header
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	rw.ResponseBody += string(p)
	rw.ContentLength = len(rw.ResponseBody)
	return len(p), nil
}

func (rw *responseWriter) SetStatusCode(code int) {
	rw.StatusCode = code
}

func (rw *responseWriter) Header() Header {
	if rw.handlerHeader == nil {
		rw.handlerHeader = Header{}
	}
	return rw.handlerHeader
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.StatusCode = code
}

type Header map[string][]string

func (h Header) Set(key, value string) {
	tmp := make([]string, 0)
	tmp = append(tmp, value)
	h[key] = tmp
}
