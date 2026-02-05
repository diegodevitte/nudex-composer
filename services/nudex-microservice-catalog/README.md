# NUDEX Catalog Service

Microservicio que maneja el catálogo de videos, productores, categorías y metadatos.

## 🚀 Stack

- **Go 1.22+** + **Gin Framework**
- **PostgreSQL** - Base de datos principal
- **Redis** - Cache de consultas
- **RabbitMQ** - Eventos

## 📊 Entidades

- **Videos**: ID, título, descripción, URL, duración, thumbnails
- **Producers**: ID, nombre, descripción, avatar, especialidades
- **Categories**: ID, slug, nombre, descripción
- **Tags**: Etiquetas para videos

## 📡 Endpoints

### Públicos

```
GET  /health                    # Health check
GET  /videos/:id                # Detalle de video
GET  /videos/search?q=term      # Búsqueda de videos
GET  /videos/category/:slug     # Videos por categoría
GET  /videos/producer/:slug     # Videos por productor
GET  /videos/random?limit=20    # Videos aleatorios
GET  /producers                 # Lista de productores
GET  /categories                # Lista de categorías
```

### Internos (API Key requerida)

```
POST /internal/videos/upsert    # Crear/actualizar video
POST /internal/producers/upsert # Crear/actualizar productor
```

## 🔧 Variables de Entorno

Ver `.env.example` para configuración completa.

## 🐋 Docker

```bash
# Desarrollo
docker build -f Dockerfile.dev -t nudex-catalog:dev .

# Producción
docker build -f Dockerfile -t nudex-catalog:prod .
```

## 📊 Features

- ✅ CRUD completo de videos
- ✅ Búsqueda por texto
- ✅ Filtros por categoría/productor
- ✅ Cache Redis para consultas frecuentes
- ✅ Migraciones automáticas
- ✅ Seed data con 20 videos
- ✅ Health checks
- ✅ Logging estructurado
- ✅ API Keys para endpoints internos
