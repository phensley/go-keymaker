package keymaker

import (
	"fmt"
	"log"
)

// LogFail emits a fatal error message if error argument is non-nil
func LogFail(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalln(fmt.Sprintf(msg, args...), err)
	}
}
