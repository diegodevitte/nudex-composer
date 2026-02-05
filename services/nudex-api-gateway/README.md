# NUDEX API Gateway

Gateway único (BFF - Backend For Frontend) que actúa como proxy entre el frontend y todos los microservicios.

## 🚀 Stack

- **Node.js 18+** + **TypeScript**
- **Fastify** - Framework web rápido
- **Redis** - Cache de respuestas
- **RabbitMQ** - Publicación de eventos

## 📡 Endpoints

### Públicos

```
GET  /health                    # Health check
GET  /api/feed/home             # Feed principal de videos
GET  /api/videos/:id            # Detalle de video específico
GET  /api/search?q=term         # Búsqueda de videos
POST /api/auth/login            # Autenticación (mock)
GET  /api/favorites             # Favoritos del usuario
POST /api/favorites/:videoId    # Agregar/quitar favorito
GET  /api/playlists             # Playlists del usuario
```

### Proxying

- **Catalog Service**: Videos, productores, categorías
- **Users Service**: Autenticación, perfil
- **Library Service**: Favoritos, historial, playlists
- **Playback Service**: Tokens de reproducción

## 🔧 Variables de Entorno

Ver `.env.example` para configuración completa.

## 🐋 Docker

```bash
# Desarrollo
docker build -f Dockerfile.dev -t nudex-api-gateway:dev .

# Producción
docker build -f Dockerfile -t nudex-api-gateway:prod .
```

## 📊 Features

- ✅ Cache Redis con TTL configurable
- ✅ Rate limiting por IP
- ✅ CORS configurado para frontend
- ✅ Logging estructurado
- ✅ Health checks
- ✅ Error handling centralizado
- ✅ Request tracing (traceId)
- ✅ Metrics básicas
