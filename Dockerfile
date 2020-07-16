FROM golang:1.14.5 as build
ENV CGO_ENABLED=0
ENV GOPATH=/go

WORKDIR /go/src/invoca/scan-launcher

COPY . .

RUN go mod download

RUN go build -mod=readonly

FROM gcr.io/distroless/static:latest

COPY --from=build /go/src/invoca/scan-launcher/scan-launcher /

USER 65534

ENTRYPOINT ["/scan-launcher"]
