FROM golang:alpine as builder
WORKDIR /build
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY . .
RUN templ generate
RUN go build -o iphistory .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /build/iphistory .
EXPOSE 8080
CMD ["./iphistory"]
