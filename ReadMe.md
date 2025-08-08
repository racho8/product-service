# Product Service – Cloud Run Deployment

This service provides CRUD APIs for managing product data using Go, Google Cloud Datastore, and is deployed on Cloud Run with automated CI/CD via GitHub Actions.

## Project Structure

- `main.go`: Entry point for the server
- `handlers/`: API handlers for CRUD operations
- `datastore/`: Datastore client wrapper
- `Dockerfile`: Container build configuration
- `cloudbuild.yaml`: Optional Cloud Build pipeline configuration
- `deploy.yaml`: GitHub Actions workflow for Cloud Run deployment

## Deployment Prerequisites

- Google Cloud Project with billing enabled
- Enabled APIs: Cloud Run, Artifact Registry, IAM
- GitHub repository (`product-service`) configured with GitHub Actions
- GitHub Secrets configured:
  - `GCP_SA_KEY`: Base64-encoded Google Service Account key

## Setup & Deploy Steps

### 1. Run Locally

#### Set Environment Variables
```bash
export GOOGLE_CLOUD_PROJECT=<YOUR_GCP_PROJECT_ID>
export PORT=8080
```

#### Set Up Google Cloud Authentication
- **Option 1: Use Application Default Credentials**
  ```bash
  gcloud auth application-default login
  ```
- **Option 2: Use a Service Account Key**
  ```bash
  export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-file.json"
  ```

#### Install Dependencies
```bash
go mod init product-service
go mod tidy
```

#### Run the Application
```bash
go run main.go
```

#### Test Local Endpoints
- **Health Check**
  ```bash
  curl http://localhost:8080/healthz
  ```
- **Create a Product**
  ```bash
  curl -X POST http://localhost:8080/products \
    -H "Content-Type: application/json" \
    -d '{"name": "Table", "category": "Office furniture", "segment": "Chair", "price": 149.99}'
  ```
- **List All Products**
  ```bash
  curl http://localhost:8080/products
  ```

### 2. Prepare Your Go App for Deployment

Ensure the app listens on the environment-provided port:

```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
log.Printf("Listening on port %s", port)
http.ListenAndServe(":"+port, router)
```

### 3. Service Account Configuration

<details>
<summary>Grant Required Roles</summary>

Run the following commands to create a service account and assign necessary roles:

```bash
# Create service account
gcloud iam service-accounts create github-deployer --display-name "GitHub Cloud Run Deployer"

# Assign roles
gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/datastore.admin"

gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/pubsub.publisher"

gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/logging.admin"

gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/viewer"

gcloud projects add-iam-policy-binding <YOUR_GCP_PROJECT_ID> \
  --member="serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.admin"

# Create JSON key
gcloud iam service-accounts keys create key.json \
  --iam-account=github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com
```

</details>

### 4. Useful Commands

<details>
<summary>View Commands</summary>

- **Check Service Account Roles**
  ```bash
  gcloud projects get-iam-policy <YOUR_GCP_PROJECT_ID> \
    --flatten="bindings[].members" \
    --format='table(bindings.role)' \
    --filter="bindings.members:serviceAccount:github-deployer@<YOUR_GCP_PROJECT_ID>.iam.gserviceaccount.com"
  ```

- **Trigger GitHub Deploy Action**
  ```bash
  git commit --allow-empty -m "trigger deploy" && git push
  ```

</details>

## API Endpoints

### Add a Product
```bash
curl -X POST https://product-service-256110662801.europe-west3.run.app/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Table", "category": "Office furniture", "segment": "Chair", "price": 149.99}'
```

### Get All Products
```bash
curl https://product-service-256110662801.europe-west3.run.app/products
```

## Common Issues & Fixes

| Issue | Fix |
|-------|-----|
| Container fails to start or port not listening | Ensure app listens on `os.Getenv("PORT")` |
| Missing environment variables (e.g., `GOOGLE_CLOUD_PROJECT`) | Configure Cloud Run environment variables or inject `.env` |
| Incorrect GitHub secret format | Verify `GCP_SA_KEY` is base64-encoded JSON |
| Timeout or cold start issues | Implement a `/healthz` endpoint for quick health checks |
| Docker push errors | Run `gcloud auth configure-docker` |
| Logs or errors not visible | View logs in [Cloud Run Logs](https://console.cloud.google.com/logs/viewer) |

