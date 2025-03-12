# Etapa de construcción
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Argumento para romper la cache
ARG CACHEBUST=1
RUN echo "Cache bust: $CACHEBUST"

# Copiar go.mod y go.sum para descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el ejecutable desde el directorio de main
RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd

# Imagen final
FROM alpine:3.17

WORKDIR /app

# Copiar el ejecutable desde la etapa de build
COPY --from=builder /app/server /app/server

# Exponer el puerto en el que corre la app
EXPOSE 8080

# Ejecutar el servicio
CMD ["/app/server"]