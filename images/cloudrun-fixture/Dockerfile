FROM golang:1.27.4-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS build
WORKDIR /src
COPY main.go main_test.go docker_test.go Dockerfile ./
RUN go test main.go main_test.go docker_test.go && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/fixture-service main.go

FROM scratch
COPY --from=build /out/fixture-service /fixture-service
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/fixture-service"]
