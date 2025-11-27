# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências de build
RUN apk add --no-cache git

# Copiar go.mod e go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates tzdata

# Copiar binário
COPY --from=builder /app/api .

# Expor porta
EXPOSE 8080

# Executar
CMD ["./api"]

