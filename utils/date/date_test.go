package date

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(GetMicroTimeStampStr())
	fmt.Println(GetTimeStr())
	fmt.Println(GetTime())
	fmt.Println(GetDateStr())
	fmt.Println(GetTimeStamp())
}
