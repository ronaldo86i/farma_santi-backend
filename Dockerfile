# Etapa de construcción
FROM golang:1.24.2-alpine AS build

# Instalar dependencias necesarias
RUN apk add --no-cache git tzdata

# Establecer la variable de entorno TZ
ENV TZ=America/La_Paz

# Configurar el huso horario
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos de configuración de Go
COPY .env go.mod go.sum ./

RUN go mod download

# Copiar el código fuente al contenedor
COPY . .

# Construir el ejecutable desde el archivo main.go en la carpeta cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o farmasanti_backend ./cmd

# Etapa final
FROM alpine:latest

# Instalar certificados para HTTPS y tzdata
RUN apk add --no-cache ca-certificates tzdata

# Establecer la variable de entorno TZ
ENV TZ=America/La_Paz

# Configurar el huso horario
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el .env desde la etapa de construcción
COPY --from=build /app/.env .env

# Copiar el binario desde la etapa de construcción
COPY --from=build /app/farmasanti_backend .

# Exponer el puerto del servidor
EXPOSE 8890

# Comando para ejecutar el binario
CMD ["/app/farmasanti_backend"]
