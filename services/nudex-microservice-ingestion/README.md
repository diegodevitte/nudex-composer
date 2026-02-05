# NUDEX Ingestion Service

Worker Python que consume jobs de RabbitMQ para sincronización de datos.

## 🚀 Stack

- **Python 3.11** + **asyncio**
- **RabbitMQ** - Consumer de jobs
- **PostgreSQL + MongoDB** - Multi DB
- **httpx** - HTTP client async

## 📡 Jobs que procesa

```
ingestion.sync.videos      # Sincronizar videos desde fuentes externas
ingestion.sync.producers   # Sincronizar productores
ingestion.process.upload   # Procesar videos subidos
ingestion.generate.thumbnails # Generar thumbnails
```

## 🔧 Features

- ✅ Consumer RabbitMQ async
- ✅ Generación de datos fake
- ✅ Upsert via Catalog API
- ✅ Retry logic con exponential backoff
- ✅ Health monitoring
- ✅ Structured logging
