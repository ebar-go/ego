package date

import (
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {
	fmt.Println(GetMicroTimeStampStr())
	fmt.Println(GetTimeStr())
	fmt.Println(GetTime())
	fmt.Println(GetDateStr())
	fmt.Println(GetTimeStamp())
	fmt.Println(GetLocalTimeZone())
	fmt.Println(GetDateTime("2020-02-02"))
}
