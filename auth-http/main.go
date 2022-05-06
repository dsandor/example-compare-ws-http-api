package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strconv"
)

type dependencies struct {
	DelayMilliseconds int64
}

func main() {
	d := dependencies{}

	envDelayMilliseconds := os.Getenv("DELAY_MILLISECONDS")
	d.DelayMilliseconds = 0

	if envDelayMilliseconds != "" {
		parsedValue, err := strconv.ParseInt(envDelayMilliseconds, 10, 64)
		if err == nil {
			d.DelayMilliseconds = parsedValue
		}
	}

	fmt.Printf("Starting http auth handler. d.DelayMilliseconds=%d\n", d.DelayMilliseconds)
	lambda.Start(d.auth)
}
