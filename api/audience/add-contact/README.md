# Add Contact Lambda Function

This Lambda function adds contacts to a Resend audience using the Resend API.

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Build the function:**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
   ```

3. **Create deployment package:**
   ```bash
   zip function.zip bootstrap
   ```

## Environment Variables

Set these environment variables in your Lambda function:

- `RESEND_API_KEY`: Your Resend API key

## Request Format

```json
{
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "audience": "newsletter"
}
```

## Response Format

**Success:**
```json
{
  "success": true,
  "message": "Contact added successfully",
  "id": "contact_id_from_resend"
}
```

**Error:**
```json
{
  "error": "Contact addition failed",
  "message": "Error details"
}
```

## Deployment

1. Upload the `function.zip` to AWS Lambda
2. Set the handler to `bootstrap`
3. Configure environment variables
4. Set appropriate IAM permissions for CloudWatch Logs 