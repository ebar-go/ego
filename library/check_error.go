package library

import "log"

// CheckErr Check the err
// You can determine whether to exit the program by setting the exit parameter
func CheckErr(msg string, err error, exit bool) {
	if err == nil {
		log.Printf("%s Finish.\n", msg)
		return
	}

	if exit {
		log.Fatalf("%s Error: %v\n", msg, err)
	} else {
		log.Printf("%s Error: %v\n", msg, err)
	}
}