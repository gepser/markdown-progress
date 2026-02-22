# Contributing

## Prerequisites

- `mise`
- `git`
- `make`

## Setup

```bash
mise trust .mise.toml
mise install
mise exec -- make setup
```

## Local workflow

Run app:

```bash
mise exec -- make run
```

Run checks before opening a PR:

```bash
mise exec -- make fmt
mise exec -- make check
```

## Smoke tests for deployed endpoint

```bash
BASE_URL=https://YOUR_DOMAIN_OR_FUNCTION_URL mise exec -- make smoke
```

## Pull requests

- Keep changes focused and small.
- Include tests for behavioral changes.
- Update `README.md` when API behavior changes.
