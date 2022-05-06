package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"time"
)

func (d *dependencies) auth(ctx context.Context, websocketEvent events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	fmt.Println("Auth handler hit.")
	fmt.Printf("%+v\n", websocketEvent)

	userContext := make(map[string]interface{})
	userContext["username"] = "some_user"
	userContext["permissions"] = []string{"Admin", "Write"}

	response := events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context:      userContext,
	}

	// allow simulating processing. For example, query an oauth provider for a jwks document, resolving a user's
	// additional claims or entitlements, or simply additional security checks.
	log.Printf("Sleeping for %dms..\n", d.DelayMilliseconds)
	time.Sleep(time.Duration(d.DelayMilliseconds) * time.Millisecond)
	log.Println("Slept.")

	return response, nil
}
