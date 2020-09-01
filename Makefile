.PHONY: test build build-and-push-image

PKG=github.com/invoca/tenable-scan-launcher

test:
	go fmt ./pkg/... ./cmd/...
	go vet ./pkg/... ./cmd/...
	go test ./pkg/... ./cmd/... --race $(PKG) -v

build:
	go fmt ./pkg/... ./cmd/...
	go vet ./pkg/... ./cmd/...
	go mod tidy
	go build -mod=readonly -o $(PWD)/tenable-scan-launcher $(PWD)/cmd/tenable-scan-launcher

build-and-push-image:
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin quay.io
	docker build -f Dockerfile -t quay.io/invoca/tenable-scan-launcher:latest -t quay.io/invoca/tenable-scan-launcher:$(VERSION) .
	echo "Pushing quay.io/invoca/tenable-scan-launcher:$(TRAVIS_BRANCH) and quay.io/invoca/tenable-scan-launcher:latest"
	docker push quay.io/invoca/tenable-scan-launcher
