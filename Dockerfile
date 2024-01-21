# syntax=docker/dockerfile:1

ARG ALPINE_VERSION="latest"
ARG GOLANG_VERSION="latest"

FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION} AS alpine
FROM --platform=$BUILDPLATFORM golang:${GOLANG_VERSION} AS golang


FROM golang as builder
WORKDIR /build
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o fail-pod-cleaner cmd/fail-pod-cleaner.go

FROM alpine
COPY --from=builder /build/fail-pod-cleaner /fail-pod-cleaner
ENTRYPOINT ["/fail-pod-cleaner"]