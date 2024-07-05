BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

precommit: ensure format generate test check
	@echo "ready to commit"

ensure:
	go mod tidy
	go mod verify
	go mod vendor

format:
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec go run -mod=vendor github.com/incu6us/goimports-reviser -project-name github.com/bborbe/backup -file-path "{}" \;

generate:
	rm -rf mocks avro
	go generate -mod=vendor ./...

test:
	go test -mod=vendor -p=$${GO_TEST_PARALLEL:-1} -cover -race $(shell go list -mod=vendor ./... | grep -v /vendor/)

check: lint vet errcheck vulncheck

vet:
	go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

lint:
	go run -mod=vendor golang.org/x/lint/golint -min_confidence 1 $(shell go list -mod=vendor ./... | grep -v /vendor/)

errcheck:
	go run -mod=vendor github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=vendor ./... | grep -v /vendor/)

vulncheck:
	go run -mod=vendor golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=vendor ./... | grep -v /vendor/)

deps:
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-config-parser@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-file@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-url@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-username@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-password@latest
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

run:
	@go mod vendor && go run -mod=vendor main.go \
	-listen="localhost:8080" \
	-sentry-dsn="$(shell teamvault-url --teamvault-config ~/.teamvault.json --teamvault-key=NqAM7q)" \
	-google-client-id="$(shell teamvault-username --teamvault-config ~/.teamvault.json --teamvault-key=gOpDMO)" \
	-google-client-secret="$(shell teamvault-password --teamvault-config ~/.teamvault.json --teamvault-key=gOpDMO)" \
	-google-hosted-domain="gmail.com" \
	-google-redirect-url="http://localhost:8080/callback" \
	-jwt-signing-key="pgTziFJhK6d5dMy7DPABcaaS8jv1liSHUQ3CzhPLBVc=" \
	-v=2

# openssl rand -base64 32
