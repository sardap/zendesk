# BUILDER
FROM golang:1.16.2 as builder

WORKDIR /app
# Resolve downloads
COPY go.mod .
COPY go.sum .
RUN go mod download
# Copy source
COPY . .
RUN go build -o main .

# APP
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main /bin/main

ENTRYPOINT [ "/bin/main" ]
CMD [ "-h" ]

