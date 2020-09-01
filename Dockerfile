FROM golang:1.14.5 as build
ENV CGO_ENABLED=0
ENV GOPATH=/go

WORKDIR /go/src/invoca/tenable-scan-launcher

COPY . .

RUN go mod download

RUN go build -mod=readonly -o /tenable-scan-launcher $PWD/cmd/tenable-scan-launcher

FROM gcr.io/distroless/static:latest

COPY --from=build /tenable-scan-launcher /

USER 65534

ENTRYPOINT ["/tenable-scan-launcher"]
