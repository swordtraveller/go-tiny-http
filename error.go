package http

func Error(w ResponseWriter, error string, code int) {
	w.SetStatusCode(code)
	w.Write([]byte(error + "\n"))
}
