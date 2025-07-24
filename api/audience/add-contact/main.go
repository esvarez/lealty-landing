package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/esvarez/lealty-landing/internal/web"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/resend/resend-go/v2"
)

type Handler struct {
	resendClient *resend.Client
}

func newHandler(resendClient *resend.Client) *Handler {
	return &Handler{
		resendClient: resendClient,
	}
}

func main() {
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	if resendAPIKey == "" {
		panic("ERROR: RESEND_API_KEY environment variable not set")
	}

	resendClient := resend.NewClient(resendAPIKey)

	handler := newHandler(resendClient)

	lambda.Start(handler)
}

// Request represents the incoming Lambda request
type Request struct {
	Email string `json:"email"`
}

// Response represents the Lambda response
type Response struct {
	Status string `json:"status"`
}

func (h *Handler) handleRequest(_ context.Context, req web.Request) (web.Response, error) {
	// Validate required fields
	request := &Request{}

	if err := json.Unmarshal([]byte(req.Body), request); err != nil {
		return web.Error(fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest), nil
	}

	if request.Email == "" {
		return web.Error("email is required", http.StatusBadRequest), nil
	}

	// Set default audience if not provided
	audience := os.Getenv("AUDIENCE_ID")
	if audience == "" {
		return web.Error("audience is required", http.StatusBadRequest), nil
	}

	params := &resend.CreateContactRequest{
		Email:        request.Email,
		Unsubscribed: false,
		AudienceId:   audience,
	}

	_, err := h.resendClient.Contacts.Create(params)
	if err != nil {
		return web.Error(fmt.Sprintf("failed to create contact: %v", err), http.StatusInternalServerError), nil
	}

	return web.JsonResponse(&Response{Status: "success"}, http.StatusCreated), nil
}
