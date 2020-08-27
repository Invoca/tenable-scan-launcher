.PHONY: test build push-image

VERSION := `git fetch --tags && git tag | sort -V | tail -1`
PKG=github.com/invoca/tenable-scan-launcher

test:
	go get golang.org/x/lint/golint
	go fmt ./pkg/... ./cmd/...
	go vet ./pkg/... ./cmd/...
	golint ./pkg/... ./cmd/...
	go test ./pkg/... ./cmd/... --race $(PKG) -v

build:
	go fmt ./pkg/... ./cmd/...
	golint ./pkg/... ./cmd/...
	go vet ./pkg/... ./cmd/...
	go mod tidy
	go build -mod=readonly -o $(PWD)/tenable-scan-launcher $(PWD)/cmd/tenable-scan-launcher

push-image:
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	docker build -f Dockerfile -t invocaops/tenable-scan-launcher:latest -t invocaops/tenable-scan-launcher:$(VERSION) .
	echo "Pushing invocaops/tenable-scan-launcher:$(VERSION) and invocaops/tenable-scan-launcher:latest"
	docker push invocaops/tenable-scan-launcher
