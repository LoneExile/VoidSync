FROM golang:1.20.4 as builder

WORKDIR /voidsync

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./cmd/voidsync

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /voidsync/

COPY --from=builder /voidsync/myapp /voidsync/

# Set the user to the non-root user
USER appuser

EXPOSE 8080

ENTRYPOINT ["/voidsync/myapp"]
