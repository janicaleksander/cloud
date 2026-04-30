FROM golang:1.26.2-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /valuation-service ./cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /valuation-service /
COPY global-bundle.pem /global-bundle-pem
EXPOSE 8083
CMD ["/valuation-service"]