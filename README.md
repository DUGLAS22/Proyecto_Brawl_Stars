# Proyecto gRPC - Información de Brawlers de Brawl Stars 

Este proyecto implementa un servicio gRPC para gestionar información sobre los Brawlers de Brawl Stars. Utiliza Go, gRPC y Protocol Buffers para definir e implementar servicios que interactúan con una tabla de datos que incluye:
- **Nombre** del Brawler
- **Tipo** ( Legendario, Épico, Superespecial, etc.)
- **Categoría** ( Asesino, Control, Apoyo, etc.)
El proyecto está empaquetado en Docker y desplegado en Azure App Service. Además, utiliza Postman para probar los servicios y una base de datos creada en Azure, configurada y gestionada a través de Azure Data Studio.

## Estructura del Proyecto

### Definición del esquema
El archivo `.proto` contiene la definición de los servicios gRPC y los mensajes necesarios para interactuar con los datos de los Brawlers.

### Servicios implementados

-**Unary** Permite buscar un Brawler por su nombre.

-**Server Streaming** Devuelve todos los Brawlers disponibles en la base de datos.

-**Client Streaming** Permite agregar múltiples Brawlers de una sola vez.

-**Bidirectional Streaming** Permite intercambiar información sobre Brawlers en tiempo real.

# Configuración del proyecto

## Requisitos

**Lenguaje** Go (1.20+)

**gRPC** Implementado con Protocol Buffers

**Docker** Para empaquetar la aplicación

**Azure App Service** Para el despliegue

## Despliegue en Azure

 ### Empaquetado con Docker
 1 Crea la imagen de Docker:
 docker build -t grpc-brawl-server:latest .
 2 Ejecuta el contenedor localmente:
 docker run -d -p 50051:50051 -p 8080:8080 --name grpc-brawl-server grpc-brawl-server:latest

 # Despliegue en Azure App Service
 **Inicia sesión en Azure CLI**
 az login
 **Inicia sesion en el registro de contenedores**
 az acr login --name brawlcr

#  Pruebas con Postman

**Descarga e instala Postman**

**Crea una nueva colección para los servicios gRPC**

**Configura los endpoints de los servicios gRPC**

**Unary: Prueba con el nombre de un Brawler**

**Server Streaming: Realiza una solicitud para listar todos los Brawlers**

**Client Streaming: Envía múltiples registros de Brawlers**

**Bidirectional Streaming: Prueba el intercambio de información de Brawlers**

**Verifica las respuestas y asegura que cumplen con lo esperado**

# Documentación de los Servicios

**Unary**

Entrada: Nombre del Brawler

Salida: Información del Brawler

**Server Streaming**

Entrada: Ninguna

Salida: Flujo de todos los Brawlers

**Client Streaming**

Entrada: Flujo de nuevos Brawlers

Salida: Cantidad de Brawlers agregados

**Bidirectional Streaming**

Entrada/Salida: Flujo bidireccional de Brawlers


 
 
 
 
