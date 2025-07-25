package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/esvarez/lealty-landing/internal/web"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/resend/resend-go/v2"
)

type Handler struct {
	resendClient *resend.Client
	audienceId   string
}

func newHandler(resendClient *resend.Client, audienceId string) *Handler {
	return &Handler{
		resendClient: resendClient,
		audienceId:   audienceId,
	}
}

func main() {
	secretName := os.Getenv("SECRET_NAME")
	region := os.Getenv("REGION")
	if secretName == "" || region == "" {
		panic("ERROR: SECRET_NAME or REGION environment variable not set")
	}

	secret := getSecret(secretName, region)

	resendClient := resend.NewClient(secret.ResendAPIKey)

	handler := newHandler(resendClient, secret.AudienceId)

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

// Secret represents the secret stored in AWS Secrets Manager
type Secret struct {
	ResendAPIKey string `json:"RESEND_API_KEY"`
	AudienceId   string `json:"AUDIENCE_ID"`
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

func getSecret(secretName string, region string) *Secret {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("ERROR: Failed to load AWS config")
	}

	svc := secretsmanager.NewFromConfig(config)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		panic("ERROR: Failed to get secret")
	}

	secret := &Secret{}
	err = json.Unmarshal([]byte(*result.SecretString), secret)
	if err != nil {
		panic("ERROR: Failed to unmarshal secret")
	}

	return secret
}
