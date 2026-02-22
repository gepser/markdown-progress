# Markdown Progress ![](https://geps.dev/progress/100)

Progress bars for markdown.

Have you ever wanted to track some progress in your markdown documents?
Well, I do, and I used `progressed.io` before but it was shut down.

So I decided to recreate it.

## Usage

```md
![](https://geps.dev/progress/10)
```

> **Note**
> I'll try to keep this domain name up as much as possible, so wish me a long life 🙂

## API Contract

### Endpoint

- `GET /progress/{percentage}`
- `HEAD /progress/{percentage}`

`percentage` must be an integer and is clamped to `0..100`.

### Query params

- `dangerColor`
- `warningColor`
- `successColor`

All color values must be 6-character hex values without `#` (example: `ff9900`).

### Response behavior

- `200 OK`: valid request, returns SVG.
- `400 Bad Request`: invalid percentage or color format.
- `405 Method Not Allowed`: any method different from `GET` or `HEAD`.
- `500 Internal Server Error`: template parse/render failure.

Headers for successful responses:

- `Content-Type: image/svg+xml`
- `Cache-Control: public, max-age=300`

## Examples

![](https://geps.dev/progress/10)

![](https://geps.dev/progress/50)

![](https://geps.dev/progress/75)

### Custom colors

You can customize colors through query params:

- `dangerColor`
- `warningColor`
- `successColor`

```md
![](https://geps.dev/progress/32?dangerColor=800000&warningColor=ff9900&successColor=006600)
```

Rendered examples:

![](https://geps.dev/progress/10?dangerColor=800000&warningColor=ff9900&successColor=006600)

![](https://geps.dev/progress/50?dangerColor=800000&warningColor=ff9900&successColor=006600)

![](https://geps.dev/progress/75?dangerColor=800000&warningColor=ff9900&successColor=006600)

## Local Development

### Prerequisites

- `mise` installed

### Setup

```bash
mise install
mise exec -- make setup
```

Run locally:

```bash
mise exec -- make run
```

Try it in a browser:

```text
http://localhost:8080/progress/76
```

### Quality checks

```bash
mise exec -- make test
mise exec -- make vet
mise exec -- make check
```

### Smoke tests against deployed URL

```bash
BASE_URL=https://geps.dev mise exec -- make smoke
```

The smoke test validates status codes, headers, and basic content contract.

If your `BASE_URL` already includes the function path (for example
`https://REGION-PROJECT.cloudfunctions.net/progress`), set:

```bash
BASE_URL=https://REGION-PROJECT.cloudfunctions.net/progress PROGRESS_PATH="" mise exec -- make smoke
```

## Deploy (Google Cloud)

Set your project first:

```bash
gcloud auth login
gcloud config set project THE_PROJECT_NAME
```

Deploy as an HTTP function with `Progress` as entrypoint:

```bash
gcloud functions deploy progress --gen2 --runtime go124 --entry-point Progress --trigger-http --allow-unauthenticated --region us-central1
```

After deploy, run smoke tests:

```bash
BASE_URL=https://YOUR_DOMAIN_OR_FUNCTION_URL mise exec -- make smoke
```

## Automated Deploy (GitHub Actions -> GCP)

This repo includes `.github/workflows/deploy.yml` to deploy automatically on
push to `master` (and manually via `workflow_dispatch`).

### 1) Configure GitHub repository variables

Required:

- `GCP_PROJECT_ID` (example: `progress-markdown`)
- `GCP_WORKLOAD_IDENTITY_PROVIDER` (full resource name)
- `GCP_SERVICE_ACCOUNT` (deployer service account email)

Optional (defaults are already set in workflow):

- `GCP_REGION` (`us-central1`)
- `GCP_FUNCTION_NAME` (`progress`)
- `GCP_RUNTIME` (`go124`)
- `GCP_ENTRY_POINT` (`Progress`)

### 2) One-time GCP setup (OIDC/WIF, no JSON keys)

Create deployer service account:

```bash
gcloud iam service-accounts create github-deployer \
  --display-name "GitHub deployer for markdown-progress"
```

Grant deploy permissions on project:

```bash
PROJECT_ID=progress-markdown
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')
SA="github-deployer@${PROJECT_ID}.iam.gserviceaccount.com"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${SA}" \
  --role="roles/cloudfunctions.developer"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${SA}" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${SA}" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${SA}" \
  --role="roles/cloudbuild.builds.editor"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${SA}" \
  --role="roles/iam.serviceAccountUser"
```

Create Workload Identity Pool + Provider:

```bash
PROJECT_ID=progress-markdown
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')
POOL_ID=github
PROVIDER_ID=github-oidc

gcloud iam workload-identity-pools create "$POOL_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --display-name="GitHub Actions Pool"

gcloud iam workload-identity-pools providers create-oidc "$PROVIDER_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="$POOL_ID" \
  --display-name="GitHub Actions Provider" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository,attribute.ref=assertion.ref"
```

Allow only this repository to impersonate the deployer service account:

```bash
PROJECT_ID=progress-markdown
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')
POOL_ID=github
REPO="gepser/markdown-progress"
SA="github-deployer@${PROJECT_ID}.iam.gserviceaccount.com"

gcloud iam service-accounts add-iam-policy-binding "$SA" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}/attribute.repository/${REPO}"
```

Provider resource name to set in GitHub variable:

```text
projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/providers/PROVIDER_ID
```

## CI

- `CI` workflow runs `go test` and `go vet` on pushes and PRs.
- `Smoke Tests` workflow can be run manually (`workflow_dispatch`) with a `base_url` input.
- `Deploy` workflow deploys to GCP on `master` using OIDC/WIF.

## Contributing

See `CONTRIBUTING.md`.
