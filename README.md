# 🎬 NUDEX - Video Platform Microservices

## Descripción

NUDEX es una plataforma de videos profesional construida con arquitectura de microservicios. Este repositorio es el orquestador principal que maneja la infraestructura de desarrollo usando Docker Compose.

## 🏗️ Arquitectura

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway    │    │  Microservicios │
│   (Next.js)     │◄──►│   (Node/Fastify) │◄──►│   Backend       │
│   Port: 3000    │    │   Port: 8080     │    │   Internal      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                    ┌─────────────────────┐
                    │   Infrastructure    │
                    │  PostgreSQL, Mongo  │
                    │  Redis, RabbitMQ    │
                    └─────────────────────┘
```

## 🚀 Stack Tecnológico

### Frontend

- **Next.js 14** (App Router) + **TypeScript**
- **Tailwind CSS** + **Zustand** (state)
- Tema NUDEX: rojo (#E1062C) y negro (#0B0B0B)

### Microservicios

| Servicio        | Stack               | Puerto | Base de Datos |
| --------------- | ------------------- | ------ | ------------- |
| **API Gateway** | Node.js + Fastify   | 8080   | Redis (cache) |
| **Catalog**     | Go + Gin            | 8081   | PostgreSQL    |
| **Users**       | NestJS + TypeScript | 8082   | PostgreSQL    |
| **Library**     | FastAPI + Python    | 8083   | MongoDB       |
| **Ingestion**   | Python Worker       | -      | Multi DB      |
| **Playback**    | Go + Gin            | 8085   | Redis         |

### Infraestructura

- **PostgreSQL 15** - catalog y users schemas
- **MongoDB 7** - library y metadata de videos
- **Redis 7** - cache y tokens efímeros
- **RabbitMQ 3.12** - mensajería asíncrona
- **Docker** + **Docker Compose** - orquestación DEV

## 📡 Mensajería (RabbitMQ)

### Exchanges

- `nudex.events` (topic) - Eventos del dominio
- `nudex.jobs` (topic) - Jobs asíncronos

### Eventos Principales

```bash
catalog.video.upserted     # Video creado/actualizado
user.created               # Usuario registrado
library.favorited          # Video agregado a favoritos
library.playlist.updated   # Playlist modificada
playback.started           # Reproducción iniciada
ingestion.sync.videos      # Job de sincronización
```

### Formato de Eventos

```json
{
  "eventId": "uuid-v4",
  "eventType": "catalog.video.upserted",
  "timestamp": "2026-02-04T10:30:00Z",
  "traceId": "uuid-v4",
  "payload": {
    "videoId": "123",
    "title": "Video Title"
  }
}
```

## 🛠️ Configuración de Desarrollo

### Pre-requisitos

- Docker Desktop 4.0+
- Git 2.30+
- Node.js 18+ (para desarrollo local)

### Instalación Rápida

```bash
# 1. Clonar repositorio principal
git clone <nudex-composer-repo-url>
cd nudex-composer

# 2. Configurar entorno
cp .env.example .env
chmod +x scripts/setup.sh

# 3. Ejecutar setup automático
./scripts/setup.sh
```

### Instalación Manual

```bash
# 1. Inicializar submódulos
git submodule update --init --recursive

# 2. Configurar variables de entorno
cp .env.example .env

# 3. Levantar infraestructura
docker compose up --build -d

# 4. Verificar servicios
docker compose ps
```

## 🌐 URLs de Acceso (DEV)

| Servicio                | URL                    | Credenciales                     |
| ----------------------- | ---------------------- | -------------------------------- |
| **Frontend**            | http://localhost:3000  | -                                |
| **API Gateway**         | http://localhost:8080  | -                                |
| **RabbitMQ Management** | http://localhost:15672 | nudex_rabbit / nudex_rabbit_pass |

## 📊 Comandos Útiles

```bash
# Ver logs de todos los servicios
docker compose logs -f

# Ver logs de un servicio específico
docker compose logs -f api-gateway

# Reiniciar un servicio
docker compose restart users-service

# Acceder a shell de un servicio
docker compose exec catalog-service sh

# Ver estado de servicios
docker compose ps

# Parar todo
docker compose down

# Limpiar volúmenes (⚠️ BORRA DATOS)
docker compose down -v
```

## 🏗️ Estructura de Archivos

```
nudex-composer/
├── apps/
│   └── nudex-frontend/              # Submódulo Frontend
├── services/
│   ├── nudex-api-gateway/           # BFF único
│   ├── nudex-microservice-catalog/  # Videos & Metadata
│   ├── nudex-microservice-users/    # Autenticación
│   ├── nudex-microservice-library/  # Favoritos/Playlists
│   ├── nudex-microservice-ingestion/# Pipeline de datos
│   └── nudex-microservice-playback/ # Streaming tokens
├── scripts/
│   ├── setup.sh                     # Setup automático
│   ├── init-postgres.sql            # Init PostgreSQL
│   └── rabbitmq-init.sh            # Init RabbitMQ
├── docker-compose.yml               # Orquestación principal
├── .env.example                     # Variables de entorno
└── README.md                        # Este archivo
```

## 🔧 APIs y Endpoints

### API Gateway (Puerto 8080)

```
GET  /health                    # Health check
GET  /api/feed/home             # Feed principal
GET  /api/videos/:id            # Detalle de video
GET  /api/search?q=term         # Búsqueda
POST /api/auth/login            # Login (mock)
GET  /api/favorites             # Favoritos del usuario
POST /api/favorites/:videoId    # Agregar favorito
```

### Microservicios (Red Interna)

- **Catalog**: `/videos`, `/producers`, `/categories`
- **Users**: `/auth/login`, `/auth/register`, `/me`
- **Library**: `/favorites`, `/history`, `/playlists`
- **Playback**: `/playback/token`, `/playback/resolve`

## 🔐 Seguridad

### Red Docker

- Red interna `nudex-dev` para comunicación entre servicios
- Solo Frontend, API Gateway y RabbitMQ Management expuestos
- Microservicios **NO** accesibles desde el host

### Autenticación

- JWT tokens para usuarios
- API Keys internas entre servicios
- Redis para cache de sesiones
- Passwords hasheados con bcrypt

## 🐛 Troubleshooting

### Error: Puerto ocupado

```bash
# Verificar procesos usando puertos
lsof -ti:3000,8080,15672
kill -9 <PID>
```

### Error: Base de datos no disponible

```bash
# Recrear volúmenes
docker compose down -v
docker compose up --build -d
```

### Error: Submódulo no inicializado

```bash
git submodule update --init --recursive
```

## 📈 Monitoring y Logs

### Health Checks

Todos los servicios implementan `/health`:

- ✅ Green: Servicio saludable
- ❌ Red: Servicio con problemas

### Structured Logging

Los servicios usan logging estructurado (JSON):

```json
{
  "timestamp": "2026-02-04T10:30:00Z",
  "level": "info",
  "service": "api-gateway",
  "traceId": "uuid-v4",
  "message": "Request processed",
  "metadata": { "userId": "123", "duration": "150ms" }
}
```

## 🚢 Despliegue

### Desarrollo

```bash
docker compose up --build -d
```

### Producción

Cada servicio incluye `Dockerfile` optimizado para producción.

## 🤝 Contribución

1. Fork del repositorio
2. Crear branch feature: `git checkout -b feature/nueva-funcionalidad`
3. Commit changes: `git commit -m 'Add nueva funcionalidad'`
4. Push branch: `git push origin feature/nueva-funcionalidad`
5. Crear Pull Request

## 📄 Licencia

Este proyecto está bajo licencia MIT. Ver [LICENSE](LICENSE) para más detalles.

## 🏷️ Versioning

Usamos [SemVer](http://semver.org/) para versionado.

## 👥 Team

- **Backend Lead** - Microservicios y APIs
- **Frontend Lead** - Next.js y UI/UX
- **DevOps Lead** - Docker y CI/CD
- **Product Owner** - Requirements y roadmap

---

**🎬 NUDEX** - _Professional Video Platform_
Built with ❤️ by the NUDEX Team

## Arquitectura del Proyecto

```
nudex-composer/
├── apps/
│   └── nudex-frontend/          # Submódulo: Frontend Next.js
├── services/                    # Preparado para microservicios futuros
├── scripts/
│   └── setup.sh                # Script de configuración inicial
├── docker-compose.yml          # Orquestación de servicios
└── README.md
```

## Requisitos Previos

- Docker y Docker Compose
- Git
- Node.js 18+ (para desarrollo local opcional)

## Instalación y Configuración

### 1. Clonar el repositorio principal

```bash
git clone <URL_REPO_PRINCIPAL> nudex-composer
cd nudex-composer
```

### 2. Inicializar submódulos

```bash
# Clonar todos los submódulos
git submodule update --init --recursive

# O usar el script automatizado
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### 3. Levantar el entorno de desarrollo

```bash
# Construir e iniciar todos los servicios
docker-compose up --build

# En segundo plano
docker-compose up -d --build
```

El frontend estará disponible en: http://localhost:3000

## Comandos Útiles

### Gestión de Docker

```bash
# Ver logs en tiempo real
docker-compose logs -f

# Ver logs solo del frontend
docker-compose logs -f nudex-frontend

# Detener servicios
docker-compose down

# Reconstruir sin caché
docker-compose build --no-cache

# Limpiar volúmenes
docker-compose down -v
```

### Gestión de Submódulos

```bash
# Actualizar todos los submódulos a la última versión
git submodule update --remote

# Actualizar submódulo específico
git submodule update --remote apps/nudex-frontend

# Clonar con submódulos incluidos
git clone --recursive <URL_REPO_PRINCIPAL>
```

## Estructura de Servicios

### Frontend (nudex-frontend)

- **Tecnología**: Next.js 14 + TypeScript + Tailwind CSS
- **Puerto**: 3000
- **Ubicación**: `/apps/nudex-frontend`
- **Repositorio**: Submódulo git independiente

### Microservicios Futuros

Los microservicios se agregarán en `/services/` como submódulos independientes.

## Desarrollo

### Agregar Nuevo Microservicio

1. Crear el repositorio del servicio independientemente
2. Agregarlo como submódulo:
   ```bash
   git submodule add <URL_SERVICIO> services/nombre-servicio
   ```
3. Agregar configuración en `docker-compose.yml`
4. Actualizar este README

### Trabajar en un Submódulo

```bash
# Navegar al submódulo
cd apps/nudex-frontend

# Crear rama y trabajar normalmente
git checkout -b feature/nueva-funcionalidad
# ... hacer cambios ...
git add . && git commit -m "feat: nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# Volver al repo principal y actualizar referencia
cd ../..
git add apps/nudex-frontend
git commit -m "chore: update frontend submodule"
```

## Configuración de Desarrollo

### Variables de Entorno

Las variables específicas de cada servicio se configuran en sus respectivos directorios:

- Frontend: `/apps/nudex-frontend/.env.local`
- Microservicios: `/services/[servicio]/.env`

### Hot Reload

El docker-compose está configurado para hot reload en modo desarrollo:

- Los cambios en el código se reflejan automáticamente
- Los `node_modules` se mantienen en volúmenes para mejor rendimiento
- WATCHPACK_POLLING habilitado para compatibilidad cross-platform

## Solución de Problemas

### El frontend no actualiza los cambios

```bash
# Verificar volúmenes
docker-compose down -v
docker-compose up --build

# O reiniciar solo el servicio
docker-compose restart nudex-frontend
```

### Problemas con submódulos

```bash
# Resetear submódulos
git submodule deinit -f .
git submodule update --init --recursive
```

### Puerto 3000 ocupado

Cambiar el puerto en `docker-compose.yml`:

```yaml
ports:
  - "3001:3000" # Cambiar primer puerto
```

## Comandos de Desarrollo Rápido

```bash
# Setup completo (primera vez)
git clone --recursive <URL> && cd nudex-composer && docker-compose up --build

# Desarrollo diario
docker-compose up -d && docker-compose logs -f

# Actualizar y desplegar
git submodule update --remote && docker-compose up --build -d
```

## Contribución

1. Cada servicio tiene su propio repositorio y flujo de contribución
2. Los cambios se integran mediante actualización de submódulos
3. El repositorio principal solo orquesta la infraestructura

## Identidad Visual NUDEX

- **Fondo**: #0B0B0B
- **Superficies**: #121212, #1E1E1E
- **Acento**: #E1062C (rojo NUDEX)
- **Texto**: #F5F5F5
- **Tipografía**: Inter + Poppins

---

**NUDEX** - Platform de contenido premium
