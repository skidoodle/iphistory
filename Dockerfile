FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN go build -o iphistory .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /build/iphistory .
COPY --from=builder /build/web ./web
EXPOSE 8080
CMD ["./iphistory"]
