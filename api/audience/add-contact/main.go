package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/esvarez/lealty-landing/internal/web"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/resend/resend-go/v2"
)

type Handler struct {
	resendClient   *resend.Client
	audienceId     string
	allowedDomains []string
}

func newHandler(resendClient *resend.Client, audienceId string, allowedDomains []string) *Handler {
	return &Handler{
		resendClient:   resendClient,
		audienceId:     audienceId,
		allowedDomains: allowedDomains,
	}
}

func main() {
	secretName := os.Getenv("SECRET_NAME")
	region := os.Getenv("REGION")
	if secretName == "" || region == "" {
		panic("ERROR: SECRET_NAME or REGION environment variable not set")
	}

	domains := os.Getenv("ALLOWED_DOMAINS")
	allowedDomains := strings.Split(domains, ",")

	secret := getSecret(secretName, region)

	resendClient := resend.NewClient(secret.ResendAPIKey)

	handler := newHandler(resendClient, secret.AudienceId, allowedDomains)

	lambda.Start(handler.handleRequest)
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

	// Validate allowed domains
	if len(h.allowedDomains) > 0 && !slices.Contains(h.allowedDomains, req.Headers["referer"]) {
		return web.Error("unauthorized: "+req.Headers["referer"], http.StatusUnauthorized), nil
	}

	// Validate required fields
	request := &Request{}

	if err := json.Unmarshal([]byte(req.Body), request); err != nil {
		return web.Error(fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest), nil
	}

	if request.Email == "" {
		return web.Error("email is required", http.StatusBadRequest), nil
	}

	if h.audienceId == "" {
		return web.Error("audience is required", http.StatusBadRequest), nil
	}

	params := &resend.CreateContactRequest{
		Email:        request.Email,
		Unsubscribed: false,
		AudienceId:   h.audienceId,
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
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
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
