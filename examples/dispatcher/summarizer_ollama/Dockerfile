# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23 AS build-stage

WORKDIR /app

COPY . .

WORKDIR /app/examples/dispatcher/dispatcher_gpt/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /summarizer

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /summarizer /summarizer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/summarizer"]