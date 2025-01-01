FROM golang:1.23 AS builder

ARG TARGETARCH=amd64
ARG TARGETOS=linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o cluster-info .

FROM alpine:latest
WORKDIR /app/
COPY templates/ templates/
COPY --from=builder /app/cluster-info .
ENTRYPOINT ["./cluster-info"]