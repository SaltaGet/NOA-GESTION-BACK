# ==========================
# STAGE 1: COMPILACIÓN
# ==========================
FROM golang:1.24-alpine AS builder

# Variables de entorno para Go
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Crear carpeta de trabajo
WORKDIR /app

# Copiar los archivos de dependencias primero (para aprovechar la cache de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código del proyecto
COPY . .

# Compilar la aplicación
RUN go build -o noa-gestion ./cmd/api/main.go

# ==========================
# STAGE 2: IMAGEN FINAL
# ==========================
FROM alpine:3.19

# Instalar certificados SSL para peticiones HTTPS (muy importante)
RUN apk add --no-cache ca-certificates

# Crear carpeta de la app
WORKDIR /app

FROM gcr.io/distroless/static-debian12

# Copiar binario desde la etapa de compilación
COPY --from=builder /app/noa-gestion .

# Exponer el puerto
EXPOSE 3000

# Comando para ejecutar la app
CMD ["./noa-gestion"]