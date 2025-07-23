package web

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

type errorResponse struct {
	Message string `json:"message"`
}

// Success returns a response with the given message and status 200
func Success(msg string) Response {
	response := map[string]string{"message": msg}
	return JsonResponse(response, http.StatusOK)
}

func ResponseMsg(msg string, statusCode int) Response {
	response := map[string]string{"message": msg}
	return JsonResponse(response, statusCode)
}

// JsonResponse returns a response with the given status code and marshalled body
// The body is marshalled to json
func JsonResponse(response any, statusCode int) Response {
	body, err := json.Marshal(response)
	if err != nil {
		return Error("error marshalling response", http.StatusInternalServerError)
	}
	return Response{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
			"Content-Type":                     "application/json",
		},
	}
}

func Error(msg string, code int) Response {
	body, err := json.Marshal(errorResponse{Message: msg})
	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       "error marshalling error response",
		}
	}
	return Response{
		StatusCode: code,
		Body:       string(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
			"Content-Type":                     "application/json",
		},
	}
}
