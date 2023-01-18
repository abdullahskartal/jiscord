FROM golang:alpine3.17@sha256:8e0cab3057642c2e4cbbbc6d02dd5dbcd23eec5fdc5f106bce3a2869abdd2a47 as builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
RUN apk add build-base
RUN apk add make

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .
EXPOSE 8000
ENTRYPOINT ["/main"]
