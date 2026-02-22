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

## CI

- `CI` workflow runs `go test` and `go vet` on pushes and PRs.
- `Smoke Tests` workflow can be run manually (`workflow_dispatch`) with a `base_url` input.

## Contributing

See `CONTRIBUTING.md`.
