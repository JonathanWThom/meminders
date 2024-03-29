FROM golang:alpine
RUN apk add build-base
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o meminders .

EXPOSE 8080

CMD ["./meminders"]
