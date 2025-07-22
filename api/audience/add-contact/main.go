package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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
		log.Println("ERROR: RESEND_API_KEY environment variable not set")
		return
	}

	resendClient := resend.NewClient(resendAPIKey)

	handler := newHandler(resendClient)

	lambda.Start(handler)
}

// Request represents the incoming Lambda request
type Request struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Audience  string `json:"audience,omitempty"`
}

// Response represents the Lambda response
type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

// ResendContact represents the contact data to send to Resend
type ResendContact struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Audience     string `json:"audience,omitempty"`
	Unsubscribed bool   `json:"unsubscribed"`
	CreatedAt    string `json:"createdAt"`
}

// ResendResponse represents the response from Resend API
type ResendResponse struct {
	ID string `json:"id"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (h *Handler) handleRequest(ctx context.Context, request Request) (Response, error) {
	// Validate required fields
	if request.Email == "" {
		log.Println("ERROR: Email is required")
		return createErrorResponse(http.StatusBadRequest, "Email is required"), nil
	}

	// Set default audience if not provided
	audience := request.Audience
	if audience == "" {
		audience = "default"
	}

	params := &resend.CreateContactRequest{
		Email:        "steve.wozniak@gmail.com",
		FirstName:    "Steve",
		LastName:     "Wozniak",
		Unsubscribed: false,
		AudienceId:   "78261eea-8f8b-4381-83c6-79fa7120f1cf",
	}

	contact, err := h.resendClient.Contacts.Create(params)
	if err != nil {
		log.Printf("ERROR: Failed to create contact: %v", err)
		return createErrorResponse(http.StatusInternalServerError, "Failed to create contact"), nil
	}

	return createErrorResponse(resp.StatusCode, fmt.Sprintf("Failed to add contact: %s", string(body))), nil
}
