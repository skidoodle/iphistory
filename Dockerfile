FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN go build -o iphistory .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/iphistory .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./iphistory"]
