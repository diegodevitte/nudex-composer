#!/bin/bash
# NUDEX Development Environment Setup Script

set -e

echo "🎬 NUDEX Development Environment Setup"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}❌ Error: Must run from nudex-composer root directory${NC}"
    exit 1
fi

echo -e "${YELLOW}📁 Creating project structure...${NC}"

# Create service directories if they don't exist
services=(
    "nudex-api-gateway"
    "nudex-microservice-catalog"
    "nudex-microservice-users"
    "nudex-microservice-library"
    "nudex-microservice-ingestion"
    "nudex-microservice-playback"
)

for service in "${services[@]}"; do
    if [ ! -d "services/$service" ]; then
        mkdir -p "services/$service"
        echo "  ✅ Created services/$service"
    fi
done

echo -e "${YELLOW}🐋 Setting up Docker environment...${NC}"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    cp .env.example .env
    echo "  ✅ Created .env from .env.example"
fi

# Make RabbitMQ init script executable
chmod +x scripts/rabbitmq-init.sh

echo -e "${YELLOW}🔧 Initializing Git submodules...${NC}"

# Initialize and update submodules
git submodule update --init --recursive

echo -e "${YELLOW}🚀 Starting NUDEX development environment...${NC}"

# Build and start all services
docker compose down --remove-orphans
docker compose up --build -d

echo -e "${GREEN}✅ NUDEX Development Environment Setup Complete!${NC}"
echo ""
echo "🌐 Access URLs:"
echo "   Frontend:           http://localhost:3000"
echo "   API Gateway:        http://localhost:8080"
echo "   RabbitMQ Management: http://localhost:15672"
echo "     (user: nudex_rabbit, pass: nudex_rabbit_pass)"
echo ""
echo "📊 Service Health:"
docker compose ps

echo ""
echo "📝 Useful commands:"
echo "   docker compose logs -f [service]  # Follow logs"
echo "   docker compose down               # Stop all services"
echo "   docker compose up -d              # Start all services"
echo "   docker compose exec [service] sh  # Shell into service"
echo ""
echo "🎬 Happy coding with NUDEX!"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}▶ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Verificar prerrequisitos
print_step "Verificando prerrequisitos..."

if ! command -v docker &> /dev/null; then
    print_error "Docker no está instalado. Instala Docker Desktop."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose no está disponible."
    exit 1
fi

if ! command -v git &> /dev/null; then
    print_error "Git no está instalado."
    exit 1
fi

print_success "Prerrequisitos verificados"

# Inicializar submódulos
print_step "Inicializando submódulos git..."

if [ -f ".gitmodules" ]; then
    git submodule update --init --recursive
    print_success "Submódulos inicializados"
else
    print_warning "No se encontró .gitmodules, configurar submódulos manualmente"
fi

# Verificar estructura
print_step "Verificando estructura del proyecto..."

if [ ! -d "apps" ]; then
    mkdir -p apps
    print_warning "Creada carpeta /apps"
fi

if [ ! -d "services" ]; then
    mkdir -p services
    print_warning "Creada carpeta /services"
fi

# Construir servicios
print_step "Construyendo servicios Docker..."

if docker-compose build; then
    print_success "Servicios construidos correctamente"
else
    print_error "Error al construir servicios"
    exit 1
fi

# Verificar red Docker
print_step "Verificando red Docker..."

if ! docker network ls | grep -q "nudex-dev"; then
    print_warning "Red nudex-dev será creada automáticamente"
fi

# Información final
echo ""
echo -e "${GREEN}🎉 Setup completado exitosamente!${NC}"
echo ""
echo "Próximos pasos:"
echo "1. Levantar servicios: ${BLUE}docker-compose up${NC}"
echo "2. Ver logs: ${BLUE}docker-compose logs -f${NC}"
echo "3. Acceder a frontend: ${BLUE}http://localhost:3000${NC}"
echo ""
echo "Para desarrollo:"
echo "• Hot reload habilitado"
echo "• Logs en tiempo real con: ${BLUE}docker-compose logs -f nudex-frontend${NC}"
echo "• Detener servicios: ${BLUE}docker-compose down${NC}"
echo ""

# Preguntar si quiere levantar servicios
read -p "¿Deseas levantar los servicios ahora? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_step "Levantando servicios..."
    docker-compose up -d
    
    echo ""
    print_success "Servicios iniciados en segundo plano"
    echo "Frontend disponible en: http://localhost:3000"
    echo ""
    echo "Ver logs con: ${BLUE}docker-compose logs -f${NC}"
fi

exit 0