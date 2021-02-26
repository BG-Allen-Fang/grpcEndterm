FROM golang
COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build /app/greet/greet_client && chmod +x /app/greet/greet_client

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app .
CMD ["./greet_client"]
EXPOSE 50051