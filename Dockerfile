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
# # STAGE FINAL: Debian con MariaDB 11.8
# # ==========================
# FROM debian:bookworm-slim

# # Instalar dependencias y MariaDB 11.x desde repositorio oficial
# RUN apt-get update && \
#     apt-get install -y --no-install-recommends \
#     ca-certificates \
#     tzdata \
#     curl \
#     wget \
#     gnupg2 \
#     lsb-release \
#     && wget https://r.mariadb.com/downloads/mariadb_repo_setup && \
#     chmod +x mariadb_repo_setup && \
#     ./mariadb_repo_setup --mariadb-server-version="mariadb-11.8" --skip-maxscale && \
#     apt-get update && \
#     apt-get install -y --no-install-recommends \
#     mariadb-client \
#     && rm -f mariadb_repo_setup && \
#     update-ca-certificates && \
#     rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# # Crear aliases para compatibilidad
# RUN ln -sf /usr/bin/mariadb-binlog /usr/bin/mysqlbinlog && \
#     ln -sf /usr/bin/mariadb-dump /usr/bin/mysqldump && \
#     ln -sf /usr/bin/mariadb /usr/bin/mysql

# # Verificar instalación
# RUN echo "=== Verificando binarios ===" && \
#     mariadb-binlog --version && \
#     mariadb-dump --version && \
#     mysqldump --version && \
#     mysqlbinlog --version && \
#     echo "✅ MariaDB 11.8 instalado correctamente"

# # Crear usuario no-root
# RUN groupadd -g 1000 appgroup && \
#     useradd -u 1000 -g appgroup -s /bin/bash -m appuser && \
#     mkdir -p /app/backups && \
#     chown -R appuser:appgroup /app

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
# STAGE 2: IMAGEN FINAL
# ==========================
FROM alpine:3.20

# Solo necesitamos mariadb-dump para backups full
# Ya NO necesitamos: mariadb-binlog, mariadb/mysql client
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    mariadb-client \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

# Crear usuario no-root
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    mkdir -p /app/backups && \
    chown -R appuser:appgroup /app

WORKDIR /app

# Copiar binario Go
COPY --from=builder --chown=appuser:appgroup /app/noa-gestion .

USER appuser

EXPOSE 3000

CMD ["./noa-gestion"]