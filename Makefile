name		?= you-didnt-define-app-name

GO		?= go
DOCKER		?= docker
DOCKER_BUILDKIT ?= 1
VERSION		?= $(shell git log --pretty=format:%h -n 1)
BUILD_TIME	?= $(shell date)
# -s removes symbol table and -ldflags -w debugging symbols
LDFLAGS		?= -asmflags -trimpath -ldflags \
		   "-s -w -X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}'"
GOARCH		?=
GOOS		?=
# CGO_ENABLED=0 == static by default
CGO_ENABLED	?= 0

COMPOSE_FILE	?= docker-compose.yml


_build: dist/$(APP_NAME)

dist/$(APP_NAME):
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		$(GO) build $(LDFLAGS) \
		-o dist/$(APP_NAME) \
		main.go

.PHONY: clean
clean:
	rm -rf dist/

install-dependencies:
	@go get -d -v ./...

lint:
	@golangci-lint run ./...

vulncheck:
	@govulncheck ./...

escape-analysis:
	$(GO) build -gcflags="-m" 2>&1

# Easiest way to get proper profiler files:
# make -B LDFLAGS=-cover _build
launch-profiler:
	$(GO) tool pprof -http=: cpu.prof

test-coverage:
	go test -failfast -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-unit:
	go test -short -failfast -race ./...

# This runs all tests, including integration tests
test-integration: start-db
	-go test -failfast -race -tags=integration ./...
	docker compose down

