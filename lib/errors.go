package lib

const (
	UnsupportedType     = iota
	ErrorInvalidMessage = iota
	ErrorConnInUse      = iota
	ErrorConnNotInUse   = iota
	ErrorAuth           = iota
	ErrorInternal       = iota
)

var errorCodes = map[int]string{
	UnsupportedType:     "unsupported_type",
	ErrorInvalidMessage: "invalid_message",
	ErrorConnInUse:      "conn_in_use",
	ErrorConnNotInUse:   "conn_not_in_use",
	ErrorAuth:           "invalid_key",
	ErrorInternal:       "internal_error",
}

//func GetError(code int) string {
//	return errorCodes[code]
//}
