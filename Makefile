.PHONY: test build build-and-push-image

PKG=github.com/invoca/tenable-scan-launcher

test:
	go fmt ./pkg/... ./cmd/...
	go vet ./pkg/... ./cmd/...
	go test ./pkg/... ./cmd/... --race $(PKG) -v

build:
	go get github.com/mattn/goveralls
	go mod vendor
	go fmt ./pkg/... ./cmd/...
	go vet ./pkg/... ./cmd/...
	go mod tidy
	ls $(GOPATH)/bin/
	$(GOPATH)/bin/goveralls -service=travis-ci
	go build -mod=readonly -o $(PWD)/tenable-scan-launcher $(PWD)/cmd/tenable-scan-launcher

build-and-push-image:
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin quay.io
	docker build -f Dockerfile -t quay.io/invoca/tenable-scan-launcher:$(TAG) .
	echo "Pushing quay.io/invoca/tenable-scan-launcher:$(TAG)"
	docker push quay.io/invoca/tenable-scan-launcher
