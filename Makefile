.PHONY: help dev dev-frontend dev-backend docker-up docker-down docker-logs install clean test

# Cores para output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m

help: ## Mostra esta mensagem de ajuda
	@echo "$(BLUE)🔍 Lupa Cidadã - Comandos Disponíveis$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# ==================== Desenvolvimento ====================

dev: ## Inicia backend e frontend em modo desenvolvimento
	@echo "$(YELLOW)🚀 Iniciando ambiente de desenvolvimento...$(NC)"
	@make -j2 dev-backend dev-frontend

dev-frontend: ## Inicia apenas o frontend
	@echo "$(BLUE)📱 Iniciando frontend...$(NC)"
	cd frontend && yarn dev

dev-backend: ## Inicia apenas o backend
	@echo "$(BLUE)⚙️  Iniciando backend...$(NC)"
	cd backend && go run cmd/api/main.go

stop: ## Para backend e frontend
	@echo "$(YELLOW)🛑 Parando serviços...$(NC)"
	@make stop-backend
	@make stop-frontend

stop-backend: ## Para o backend
	@echo "$(BLUE)🛑 Parando backend...$(NC)"
	@pkill -f "go-build.*main" 2>/dev/null || true
	@pkill -f "go run cmd/api" 2>/dev/null || true

stop-frontend: ## Para o frontend
	@echo "$(BLUE)🛑 Parando frontend...$(NC)"
	@pkill -f "vite" 2>/dev/null || true

# ==================== Docker ====================

docker-up: ## Sobe todos os containers
	@echo "$(YELLOW)🐳 Subindo containers...$(NC)"
	docker-compose up -d

docker-down: ## Para todos os containers
	@echo "$(YELLOW)🐳 Parando containers...$(NC)"
	docker-compose down

docker-logs: ## Mostra logs dos containers
	docker-compose logs -f

docker-build: ## Reconstrói as imagens
	docker-compose build --no-cache

# ==================== Instalação ====================

install: ## Instala dependências do frontend e backend
	@echo "$(YELLOW)📦 Instalando dependências...$(NC)"
	@make install-frontend
	@make install-backend

install-frontend: ## Instala dependências do frontend
	@echo "$(BLUE)📦 Instalando dependências do frontend...$(NC)"
	cd frontend && yarn install

install-backend: ## Instala dependências do backend
	@echo "$(BLUE)📦 Instalando dependências do backend...$(NC)"
	cd backend && go mod download

# ==================== Banco de Dados ====================

db-up: ## Inicia apenas MongoDB e Redis
	docker-compose up -d mongodb redis meilisearch

sync: ## Sincroniza dados das APIs públicas (Câmara + Senado + Presidente + Governadores)
	@echo "$(YELLOW)🔄 Sincronizando dados das APIs públicas...$(NC)"
	cd backend && go run cmd/sync/main.go -all

sync-all: ## Sincroniza TUDO: políticos + votações + proposições + despesas + presenças
	@echo "$(YELLOW)🔄 Sincronizando TODOS os dados...$(NC)"
	@echo "$(BLUE)1️⃣  Sincronizando políticos...$(NC)"
	cd backend && go run cmd/sync/main.go -all
	@echo "$(BLUE)2️⃣  Sincronizando votações, proposições, despesas e presenças...$(NC)"
	cd backend && go run cmd/sync/main.go -votacoes -proposicoes -despesas -presencas -ano $(shell date +%Y)

sync-camara: ## Sincroniza apenas deputados da Câmara
	@echo "$(YELLOW)🔄 Sincronizando deputados da Câmara...$(NC)"
	cd backend && go run cmd/sync/main.go -camara

sync-camara-completo: ## Sincroniza deputados + todos os dados da Câmara (votações, proposições, despesas, presenças)
	@echo "$(YELLOW)🔄 Sincronizando dados completos da Câmara...$(NC)"
	@echo "$(BLUE)1️⃣  Sincronizando deputados...$(NC)"
	cd backend && go run cmd/sync/main.go -camara
	@echo "$(BLUE)2️⃣  Sincronizando votações, proposições, despesas e presenças...$(NC)"
	cd backend && go run cmd/sync/main.go -votacoes -proposicoes -despesas -presencas -ano $(shell date +%Y)

sync-votacoes: ## Sincroniza votações da Câmara (ano atual)
	@echo "$(YELLOW)🔄 Sincronizando votações...$(NC)"
	cd backend && go run cmd/sync/main.go -votacoes -ano $(shell date +%Y)

sync-proposicoes: ## Sincroniza proposições da Câmara (ano atual)
	@echo "$(YELLOW)🔄 Sincronizando proposições...$(NC)"
	cd backend && go run cmd/sync/main.go -proposicoes -ano $(shell date +%Y)

sync-despesas: ## Sincroniza despesas da Câmara (ano atual)
	@echo "$(YELLOW)🔄 Sincronizando despesas...$(NC)"
	cd backend && go run cmd/sync/main.go -despesas -ano $(shell date +%Y)

sync-presencas: ## Sincroniza presenças em eventos da Câmara (ano atual)
	@echo "$(YELLOW)🔄 Sincronizando presenças...$(NC)"
	cd backend && go run cmd/sync/main.go -presencas -ano $(shell date +%Y)

sync-senado: ## Sincroniza apenas senadores do Senado
	@echo "$(YELLOW)🔄 Sincronizando senadores do Senado...$(NC)"
	cd backend && go run cmd/sync/main.go -senado

sync-presidente: ## Sincroniza apenas Presidente da República
	@echo "$(YELLOW)🔄 Sincronizando Presidente da República...$(NC)"
	cd backend && go run cmd/sync/main.go -presidente

sync-governadores: ## Sincroniza apenas Governadores
	@echo "$(YELLOW)🔄 Sincronizando Governadores...$(NC)"
	cd backend && go run cmd/sync/main.go -governadores

# ==================== Testes ====================

test: ## Roda todos os testes
	@make test-frontend
	@make test-backend

test-frontend: ## Roda testes do frontend
	cd frontend && yarn test

test-backend: ## Roda testes do backend
	cd backend && go test ./...

# ==================== Limpeza ====================

clean: ## Remove artefatos de build
	@echo "$(YELLOW)🧹 Limpando artefatos...$(NC)"
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/bin

# ==================== Produção ====================

build: ## Build para produção
	@echo "$(YELLOW)🏗️  Construindo para produção...$(NC)"
	cd frontend && yarn build
	cd backend && go build -o bin/api cmd/api/main.go

