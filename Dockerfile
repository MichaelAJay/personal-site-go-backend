FROM golang:1.21.7-alpine3.19

WORKDIR /app

ENV ENV=local
ENV GCP_PROJECT=ezman-386111

COPY . .

RUN go mod download
# At this point, you have a Go toolchain version 1.19.x and all your Go dependencies installed inside the image.

# Now, to compile your application, use the familiar RUN command:
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/server

CMD ["/docker-gs-ping"]
