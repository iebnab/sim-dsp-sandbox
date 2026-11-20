# Synthetic HTTP/job fixture with no outbound dependencies.
# The deliberately old immutable runtime base exercises vulnerability fix age.
FROM golang:1.26.2-alpine AS build
WORKDIR /src
COPY main.go main_test.go docker_test.go Dockerfile ./
RUN go test main.go main_test.go docker_test.go && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/fixture-service main.go

FROM alpine:3.18.0@sha256:02bb6f428431fbc2809c5d1b41eab5a68350194fb508869a33cb1af4444c9b11
# This sentinel is deliberately fake and not an issued service-account credential.
RUN printf '%s\n' '{"type": "service_account", "project_id": "synthetic-f11-sentinel", "private_key_id": "NOT-A-REAL-KEY"}' \
    > /synthetic-credential.json
COPY --from=build /out/fixture-service /fixture-service
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/fixture-service"]
