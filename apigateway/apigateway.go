package apigateway

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// ResponseSuccessful returns a 200 response for API Gateway that allows cors
func ResponseSuccessful(v interface{}) events.APIGatewayProxyResponse {
	b, err := json.Marshal(v)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}
	}

	return ResponseSuccessfulString(string(b))
}

// ResponseSuccessfulString returns a 200 response for API Gateway that allows cors
func ResponseSuccessfulString(body string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type,Authorization,dfd-auth"
	resp.Body = body
	resp.StatusCode = 200
	return resp
}

// ResponseUnsuccessful return response for API Gateway that allows cors
func ResponseUnsuccessful(statusCode int) events.APIGatewayProxyResponse {
	return ResponseUnsuccessfulString(statusCode, "")
}

// ResponseUnsuccessfulString return response for API Gateway that allows cors
func ResponseUnsuccessfulString(statusCode int, body string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type,Authorization,dfd-auth"
	resp.Body = body
	resp.StatusCode = statusCode
	return resp
}
