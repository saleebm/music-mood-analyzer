package shared

import (
	"log"
	"runtime/debug"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("Error: \"%s\": \"%s\"\nStack trace:\n%s", msg, err, debug.Stack())
	}
}
