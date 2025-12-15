# Construcción (Builder)
FROM golang:1.23-alpine AS build

# Instalar dependencias del sistema
RUN apk add --no-cache git tzdata

WORKDIR /app

# Copiar SOLO definicion de dependencias (para aprovechar caché de Docker)
# ¡OJO: Quitamos el .env de aquí por seguridad!
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar
# CGO_ENABLED=0: Asegura binario estático (sin depencias de C)
# -ldflags="-w -s": Reduce el peso del binario quitando info de debug
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o farmasanti_backend "./cmd"


# Producción
FROM alpine:latest

ENV TZ=America/La_Paz

# Instalar dependencias de Runtime
# Incluimos 'postgresql-client' para que funcione pg_dump (Backups)
RUN apk add --no-cache ca-certificates tzdata postgresql-client

# Configurar Zona Horaria
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# SEGURIDAD: Crear un usuario no-root llamado 'ron86'
RUN addgroup -S ron86 && adduser -S ron86 -G ron86

WORKDIR /app

# Copiar el binario y asignar permisos al usuario
COPY --from=build /app/farmasanti_backend .

# Dar permisos al usuario (para que pg_dump pueda escribir)
RUN chown -R ron86:ron86 /app

# Cambiar al usuario seguro
USER ron86

EXPOSE 8890

CMD ["/app/farmasanti_backend"]