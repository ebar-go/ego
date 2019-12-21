package helper

import "log"

// CheckError
func CheckError(msg string, err error) {
	if err != nil {
		log.Printf("%s Error: %v\n", msg, err)
	}
}

// FatalError program will exit when err not nil
func FatalError(msg string , err error)  {
	if err != nil {
		log.Fatalf("%s Error: %v\n", msg, err)
	}
}
