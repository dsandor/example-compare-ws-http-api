package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

func (d *dependencies) list(ctx context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Default route handler hit.")
	fmt.Printf("%+v\n", websocketEvent)

	response := &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "Hello, from hello http api route.",
		IsBase64Encoded: false,
	}

	return *response, nil
}
