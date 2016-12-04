package keymaker

import "github.com/valyala/gorpc"

// KeyResponse is the response for a single RPC request
type KeyResponse struct {
	Key     []byte
	Error   int
	Message string
}

const (
	// ErrOK - no error occurred
	ErrOK = iota
	// ErrBadRequest - request was invalid
	ErrBadRequest
	// ErrKeyGen - key generation failed
	ErrKeyGen
)

func init() {
	gorpc.RegisterType(&KeyResponse{})
}
