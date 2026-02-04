# ETAPA 1: Construcción (Builder)
FROM golang:1.25-alpine AS builder
WORKDIR /app

# 1. Instalar dependencias necesarias para compilar en Alpine
RUN apk add --no-cache git

# 2. Cachear dependencias (esto ahorra mucho tiempo)
COPY go.mod go.sum ./
RUN go mod download

# 3. Copiar el código fuente
COPY . .

# 4. Compilar el binario de debug
# CGO_ENABLED=0 crea un binario estático (no necesita librerías externas)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/debug/main.go

# ETAPA 2: Producción (Imagen final ultra ligera)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 5. Copiar el binario desde el builder
COPY --from=builder /app/main .

# 6. COMENTADO: Actívalo cuando crees tu carpeta de migraciones
# COPY --from=builder /app/cmd/migrate/migrations ./migrations

# 7. Configuración de ejecución
# Puerto 8888 para servidor TCP debug
EXPOSE 8888

CMD ["./main"]
