package serv

import (
	"io"
)

//Response ...
type Response struct {
	Headers  map[string]string
	Code     int
	CodeMsg  string
	BodySize int64
	Body     io.Reader
}
