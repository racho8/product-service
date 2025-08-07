# Product Service – Cloud Run Deployment

This service provides CRUD APIs for managing product data using **Go**, **Google Cloud Datastore**, and is deployed on **Cloud Run** with automated CI/CD via **GitHub Actions**.

---

## 🚀 Project Structure

- `main.go` – Entry point for the server
- `handlers/` – API handlers for CRUD operations
- `datastore/` – Datastore client wrapper
- `Dockerfile` – Container build config
- `cloudbuild.yaml` *(optional)* – For Cloud Build pipeline
- `deploy.yaml` – GitHub Actions workflow to deploy to Cloud Run

---

## ✅ Deployment Prerequisites

1. Google Cloud Project with billing enabled.
2. Cloud Run, Artifact Registry, and IAM APIs enabled.
3. `product-service` GitHub repo connected with GitHub Actions.
4. Secret manager configured in GitHub with:
    - `GCP_SA_KEY` – base64-encoded Google Service Account key

---

## 🛠️ Setup & Deploy Steps

### 1. Prepare Your Go App

- Accept `PORT` from environment:

  ```go
  port := os.Getenv("PORT")
  if port == "" {
      port = "8080"
  }
  log.Printf("Listening on port %s", port)
  http.ListenAndServe(":"+port, router)


## Useful commands


Roles to SA (without using terraform and doing it from local machine)

### Create service account
gcloud iam service-accounts create github-deployer --display-name "GitHub Cloud Run Deployer"

### Grant required roles
gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/run.admin"

gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/datastore.admin"

gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/pubsub.publisher"

gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/logging.admin"

gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/viewer"

gcloud projects add-iam-policy-binding ingka-find-racho8-dev \
--member="serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com" \
--role="roles/artifactregistry.admin"

### Create JSON key and save locally
gcloud iam service-accounts keys create key.json \
--iam-account=github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com

### TO CHECK THE ROLES OF A SA
gcloud projects get-iam-policy ingka-find-racho8-dev \
--flatten="bindings[].members" \
--format='table(bindings.role)' \
--filter="bindings.members:serviceAccount:github-deployer@ingka-find-racho8-dev.iam.gserviceaccount.com"


## To Make a dummy commit and trigger github deploy action
git commit --allow-empty -m "trigger deploy" && git push


# CURL COMMANDS TO THE END POINTS

### TO ADD A PRODUCT
curl -X POST https://product-service-256110662801.europe-west3.run.app/products \
-H "Content-Type: application/json" \
-d '{"name": "Table", "category": "Office furniture", "segment": "Chair", "price": 149.99}'

### TO GET ALL PRODUCTS
curl https://product-service-256110662801.europe-west3.run.app/products


# Common Issues (and how to avoid them)

| Issue | Fix |
|-------|-----|
| ❌ Container failed to start / port not listening | Ensure app listens on `os.Getenv("PORT")` |
| ❌ Missing environment variables like `GOOGLE_CLOUD_PROJECT` | Use Cloud Run env var config or .env injection |
| ❌ Incorrect GitHub secret format for service account | Make sure `GCP_SA_KEY` is base64-encoded JSON |
| ❌ Timeout or cold start issues | Add a fast `/healthz` endpoint |
| ❌ Docker push errors | Ensure `gcloud auth configure-docker` is run |
| ❌ Not seeing logs/errors | Check Cloud Run logs [here](https://console.cloud.google.com/logs/viewer) |