package echohandler

import (
	"fmt"
	"time"
)

func mustParseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

var errFake = fmt.Errorf("fake error")
