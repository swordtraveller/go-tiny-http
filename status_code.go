package http

const (
	StatusOK                  = 200
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

var MessageMap = map[int]string{
	StatusOK:       "OK",
	StatusNotFound: "Not Found",
}
