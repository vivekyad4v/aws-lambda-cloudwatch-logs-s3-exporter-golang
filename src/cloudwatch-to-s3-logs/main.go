package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func mainHandler(ctx context.Context) (output string, err error) {
	fmt.Println("Done")
	return "Sucess", nil
}

func main() {
	lambda.Start(mainHandler)
}
