package apigateway

import (
	"github.com/aws/aws-lambda-go/events"
)

// ResponseSuccessful returns a 200 response for API Gateway that allows cors
func ResponseSuccessful(body string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type,Authorization,dfd-auth"
	resp.Body = body
	resp.StatusCode = 200
	return resp
}
