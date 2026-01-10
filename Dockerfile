FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

FROM golang:alpine AS build-stage
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=generate-stage /app /app

RUN CGO_ENABLED=0 go build \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o iphistory .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=build-stage /app/iphistory .

EXPOSE 8080
CMD ["./iphistory"]
