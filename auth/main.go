package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

type dependencies struct {
}

func main() {
	d := dependencies{}

	fmt.Printf("%+v\n", os.Environ())
	fmt.Println("Starting auth handler.")
	lambda.Start(d.auth)
}
