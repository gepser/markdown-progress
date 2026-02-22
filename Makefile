GO ?= go

.PHONY: setup run test vet fmt tidy check smoke

setup:
	$(GO) mod download

run:
	$(GO) run ./cmd/main.go

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

check: test vet

smoke:
	@test -n "$(BASE_URL)" || (echo "BASE_URL is required"; exit 1)
	bash ./scripts/smoke.sh
