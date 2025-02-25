FROM golang:1.23
WORKDIR /app
COPY . .
RUN go build -o url-shortener ./cmd/url-shortener
EXPOSE 8082
CMD ["./url-shortener"]