# NUDEX Users Service

Microservicio de autenticación y gestión de usuarios.

## 🚀 Stack

- **NestJS** + **TypeScript**
- **PostgreSQL** - Base de datos de usuarios
- **JWT** - Autenticación
- **bcrypt** - Hash de passwords
- **RabbitMQ** - Eventos

## 📊 Entidades

- **Users**: ID, email, password, name, avatar, roles
- **Sessions**: Tokens JWT activos

## 📡 Endpoints

```
GET  /health                # Health check
POST /auth/register         # Registro de usuario
POST /auth/login            # Login y JWT token
POST /auth/refresh          # Refresh token
GET  /me                    # Perfil del usuario
PATCH /me                   # Actualizar perfil
```

## 🔧 Features

- ✅ Registro y login seguro
- ✅ JWT tokens con refresh
- ✅ Hash bcrypt para passwords
- ✅ Validación de datos
- ✅ Rate limiting
- ✅ Eventos RabbitMQ
- ✅ Health checks
