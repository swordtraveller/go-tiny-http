package http

const (
	StatusOK                  = 200
	StatusAccepted            = 202
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

var MessageMap = map[int]string{
	StatusOK:                  "OK",
	StatusAccepted:            "Accepted",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
}
