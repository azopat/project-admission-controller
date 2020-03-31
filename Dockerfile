FROM golang:1.13-alpine AS builder

RUN apk add --update --no-cache ca-certificates git
RUN mkdir -p /build
WORKDIR /build
COPY go.* /build/
RUN go mod download
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /project-admission-validation ./cmd/

FROM alpine:3.10

COPY --from=builder /project-admission-validation /project-admission-validation
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/project-admission-validation"]
