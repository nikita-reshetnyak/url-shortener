FROM golang:1.23
WORKDIR /app
COPY . .
COPY ./config/local.yaml /app/config/local.yaml
RUN go build -o url-shortener ./cmd/url-shortener
EXPOSE 8082
CMD ["./url-shortener"]