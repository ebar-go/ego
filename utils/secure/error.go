package secure

import (
	"fmt"
	"log"
)

// FatalError
func FatalError(msg string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s:%v", msg, err))
	}

	log.Printf("%s Success\n", msg)
}
