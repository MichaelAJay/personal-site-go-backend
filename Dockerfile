FROM golang:1.21.7-alpine3.19 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download
# At this point, you have a Go toolchain version 1.19.x and all your Go dependencies installed inside the image.

# Now, to compile your application, use the familiar RUN command:
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/server

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /docker-gs-ping /docker-gs-ping

USER nonroot:nonroot

ENTRYPOINT ["/docker-gs-ping"]
