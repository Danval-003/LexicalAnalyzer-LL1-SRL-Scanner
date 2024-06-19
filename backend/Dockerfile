# Usar una imagen base de Golang
FROM golang:1.20-alpine

# Setear el directorio de trabajo
WORKDIR /app

# Copiar los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el código fuente al contenedor
COPY . .

# Compilar el proyecto
RUN go build -o main .

# Copiar las variables de entorno de Render en un archivo .env
RUN printenv > .env

# Exponer el puerto de la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
