# # ==========================
# # STAGE 1: COMPILACIÓN
# # ==========================
# FROM golang:1.24-alpine AS builder

# # Variables de entorno para Go
# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Instalar dependencias necesarias
# RUN apk add --no-cache git

# # Crear carpeta de trabajo
# WORKDIR /app

# # Copiar los archivos de dependencias primero (para aprovechar la cache de Docker)
# COPY go.mod go.sum ./
# RUN go mod download

# # Copiar todo el código del proyecto
# COPY . .

# # Compilar la aplicación
# RUN go build -o noa-gestion ./cmd/api/main.go

# # ==========================
# # STAGE 2: IMAGEN FINAL
# # ==========================
# FROM alpine:3.19

# # Instalar certificados SSL para peticiones HTTPS (muy importante)
# RUN apk add --no-cache ca-certificates

# # Crear carpeta de la app
# WORKDIR /app

# FROM gcr.io/distroless/static-debian12

# # Copiar binario desde la etapa de compilación
# COPY --from=builder /app/noa-gestion .

# # Exponer el puerto
# EXPOSE 3000

# # Comando para ejecutar la app
# CMD ["./noa-gestion"]






# # ==========================
# # STAGE 1: COMPILACIÓN
# # ==========================
# FROM golang:1.24-alpine AS builder

# # Variables de entorno para Go
# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Instalar dependencias necesarias
# RUN apk add --no-cache git

# # Crear carpeta de trabajo
# WORKDIR /app

# # Copiar los archivos de dependencias primero (para aprovechar la cache de Docker)
# COPY go.mod go.sum ./
# RUN go mod download

# # Copiar todo el código del proyecto
# COPY . .

# # Compilar la aplicación con optimizaciones de tamaño
# RUN go build -ldflags="-s -w" -o noa-gestion ./cmd/api/main.go

# # ==========================
# # STAGE 2: IMAGEN FINAL
# # ==========================
# FROM alpine:3.19

# # Instalar certificados SSL y herramientas de MariaDB
# # Usar --no-cache y limpiar después para reducir tamaño
# RUN apk add --no-cache \
#     ca-certificates \
#     mariadb-client \
#     mariadb-connector-c \
#     tzdata \
#     && rm -rf /var/cache/apk/* \
#     && rm -rf /tmp/*

# # Crear usuario no-root para seguridad
# RUN addgroup -g 1000 appgroup && \
#     adduser -D -u 1000 -G appgroup appuser

# # Crear carpeta de la app y backups
# WORKDIR /app
# RUN mkdir -p /app/backups && chown -R appuser:appgroup /app

# # Copiar binario desde la etapa de compilación
# COPY --from=builder --chown=appuser:appgroup /app/noa-gestion .

# # Cambiar a usuario no-root
# USER appuser

# # Exponer el puerto
# EXPOSE 3000

# # Comando para ejecutar la app
# CMD ["./noa-gestion"]




# # ==========================
# # STAGE 1: COMPILACIÓN
# # ==========================
# FROM golang:1.24-alpine AS builder

# # Variables de entorno para Go
# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Instalar git para dependencias
# RUN apk add --no-cache git upx

# # Crear carpeta de trabajo
# WORKDIR /app

# # Copiar dependencias primero para cache
# COPY go.mod go.sum ./
# RUN go mod download

# # Copiar todo el código
# COPY . .

# # Compilar la aplicación con optimizaciones de tamaño
# RUN go build -ldflags="-s -w" -o noa-gestion ./cmd/api/main.go

# # Comprimir el binario con upx para reducir tamaño
# RUN upx --best --lzma /app/noa-gestion

# # ==========================
# # STAGE 2: MARIA DB MINIMAL
# # ==========================
# FROM alpine:3.19 AS mariadb-client
# RUN apk add --no-cache \
#     mariadb-client \
#     mariadb-connector-c \
#     ca-certificates \
#     tzdata \
#     && update-ca-certificates \
#     && rm -rf /var/cache/apk/* /tmp/* /var/lib/apk/* /usr/share/doc /usr/share/man

# # ==========================
# # STAGE 3: IMAGEN FINAL
# # ==========================
# FROM alpine:3.19

# # Crear usuario no-root
# RUN addgroup -g 1000 appgroup && \
#     adduser -D -u 1000 -G appgroup appuser

# # Crear carpeta de la app y backups
# WORKDIR /app
# RUN mkdir -p /app/backups && chown -R appuser:appgroup /app

# # Copiar binario Go desde builder
# COPY --from=builder --chown=appuser:appgroup /app/noa-gestion /app/noa-gestion

# # Copiar MariaDB client y certificados desde stage intermedio
# COPY --from=mariadb-client /usr/bin/mysql /usr/bin/mysql
# COPY --from=mariadb-client /usr/lib/libmariadb.so* /usr/lib/
# COPY --from=mariadb-client /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=mariadb-client /usr/share/zoneinfo /usr/share/zoneinfo

# # Cambiar a usuario no-root
# USER appuser

# # Exponer el puerto
# EXPOSE 3000

# # Comando para ejecutar la app
# CMD ["./noa-gestion"]


# # ==========================
# # STAGE 1: COMPILACIÓN
# # ==========================
# FROM golang:1.24-alpine AS builder

# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# RUN apk add --no-cache git upx

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -ldflags="-s -w -extldflags '-static'" -trimpath -o noa-gestion ./cmd/api/main.go
# RUN upx --best --lzma noa-gestion || echo "UPX failed, continuing..."

# # ==========================
# # STAGE 2: DESCARGAR MARIADB 11.8
# # ==========================
# FROM alpine:3.20 AS downloader

# RUN apk add --no-cache wget tar gzip

# # Descargar MariaDB 11.8.2 (última versión disponible de 11.8)
# RUN wget -O /tmp/mariadb.tar.gz https://archive.mariadb.org/mariadb-11.8.2/bintar-linux-systemd-x86_64/mariadb-11.8.2-linux-systemd-x86_64.tar.gz && \
#     cd /tmp && \
#     tar -xzf mariadb.tar.gz && \
#     mv mariadb-11.8.2-linux-systemd-x86_64 /mariadb

# # ==========================
# # STAGE FINAL
# # ==========================
# FROM alpine:3.20

# # Instalar dependencias necesarias para ejecutar binarios de MariaDB
# RUN apk add --no-cache \
#     ca-certificates \
#     tzdata \
#     bash \
#     libstdc++ \
#     libgcc \
#     ncurses-libs \
#     libaio \
#     gcompat \
#     && update-ca-certificates \
#     && rm -rf /var/cache/apk/*

# # Copiar binarios de MariaDB desde el downloader
# COPY --from=downloader /mariadb/bin/mariadb-binlog /usr/bin/mariadb-binlog
# COPY --from=downloader /mariadb/bin/mariadb /usr/bin/mariadb
# COPY --from=downloader /mariadb/bin/mariadb-dump /usr/bin/mariadb-dump

# # Copiar librerías necesarias
# COPY --from=downloader /mariadb/lib/libmariadb.so.3 /usr/lib/

# # Crear aliases
# RUN ln -sf /usr/bin/mariadb-binlog /usr/bin/mysqlbinlog && \
#     ln -sf /usr/bin/mariadb /usr/bin/mysql && \
#     ln -sf /usr/bin/mariadb-dump /usr/bin/mysqldump

# # Verificar versión (sin ejecutar para evitar errores en build)
# RUN echo "=== Binarios de MariaDB 11.8.2 instalados ===" && \
#     ls -lh /usr/bin/mariadb-binlog && \
#     ls -lh /usr/bin/mysqlbinlog && \
#     echo "✅ Instalación completa"

# # Crear usuario no-root
# RUN addgroup -g 1000 -S appgroup && \
#     adduser -D -u 1000 -S -G appgroup appuser && \
#     mkdir -p /app/backups && chown -R appuser:appgroup /app

# WORKDIR /app

# # Copiar binario Go
# COPY --from=builder --chown=appuser:appgroup /app/noa-gestion .

# USER appuser

# EXPOSE 3000

# CMD ["./noa-gestion"]

# ==========================
# STAGE 1: COMPILACIÓN
# ==========================
FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache git upx

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w -extldflags '-static'" -trimpath -o noa-gestion ./cmd/api/main.go
RUN upx --best --lzma noa-gestion || echo "UPX failed, continuing..."

# ==========================
# STAGE FINAL: Debian con MariaDB 11.8
# ==========================
FROM debian:bookworm-slim

# Instalar dependencias y MariaDB 11.x desde repositorio oficial
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    curl \
    wget \
    gnupg2 \
    lsb-release \
    && wget https://r.mariadb.com/downloads/mariadb_repo_setup && \
    chmod +x mariadb_repo_setup && \
    ./mariadb_repo_setup --mariadb-server-version="mariadb-11.8" && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    mariadb-client \
    && rm -f mariadb_repo_setup && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Crear aliases para compatibilidad
RUN ln -sf /usr/bin/mariadb-binlog /usr/bin/mysqlbinlog && \
    ln -sf /usr/bin/mariadb-dump /usr/bin/mysqldump && \
    ln -sf /usr/bin/mariadb /usr/bin/mysql

# Verificar instalación
RUN echo "=== Verificando binarios ===" && \
    mariadb-binlog --version && \
    mariadb-dump --version && \
    mysqldump --version && \
    mysqlbinlog --version && \
    echo "✅ MariaDB 11.8 instalado correctamente"

# Crear usuario no-root
RUN groupadd -g 1000 appgroup && \
    useradd -u 1000 -g appgroup -s /bin/bash -m appuser && \
    mkdir -p /app/backups && \
    chown -R appuser:appgroup /app

WORKDIR /app

# Copiar binario Go
COPY --from=builder --chown=appuser:appgroup /app/noa-gestion .

USER appuser

EXPOSE 3000

CMD ["./noa-gestion"]