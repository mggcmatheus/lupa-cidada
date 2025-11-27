package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lupa-cidada/backend/internal/config"
	"github.com/lupa-cidada/backend/internal/handlers"
	"github.com/lupa-cidada/backend/internal/repository"
	"github.com/lupa-cidada/backend/internal/services"
	"github.com/lupa-cidada/backend/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Carregar configuração
	cfg := config.Load()

	var db *mongo.Database
	var mongoClient *mongo.Client

	// Só conecta ao MongoDB se não estiver em modo debug
	if !cfg.Debug {
		var err error
		mongoClient, err = database.NewMongoClient(cfg.MongoURI)
		if err != nil {
			log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		}
		defer mongoClient.Disconnect(context.Background())
		db = mongoClient.Database("lupa_cidada")
		log.Println("📦 Conectado ao MongoDB")
	} else {
		log.Println("🔧 Modo DEBUG ativado - usando dados mockados")
	}

	// Inicializar repositórios (podem ser nil em modo debug)
	var politicoRepo *repository.PoliticoRepository
	var votacaoRepo *repository.VotacaoRepository
	var despesaRepo *repository.DespesaRepository
	var proposicaoRepo *repository.ProposicaoRepository

	if db != nil {
		politicoRepo = repository.NewPoliticoRepository(db)
		votacaoRepo = repository.NewVotacaoRepository(db)
		despesaRepo = repository.NewDespesaRepository(db)
		proposicaoRepo = repository.NewProposicaoRepository(db)
	}

	// Inicializar serviços (passa cfg.Debug para decidir fonte dos dados)
	politicoService := services.NewPoliticoService(cfg.Debug, politicoRepo, votacaoRepo, despesaRepo, proposicaoRepo)

	// Inicializar handlers
	politicoHandler := handlers.NewPoliticoHandler(politicoService)
	filtrosHandler := handlers.NewFiltrosHandler(db, cfg.Debug)
	estatisticasHandler := handlers.NewEstatisticasHandler(politicoService)

	// Configurar Echo
	e := echo.New()
	e.HideBanner = true

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"debug":  cfg.Debug,
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API v1
	api := e.Group("/api/v1")

	// Rotas de políticos
	politicos := api.Group("/politicos")
	politicos.GET("", politicoHandler.Listar)
	politicos.GET("/comparar", politicoHandler.Comparar) // Deve vir antes de /:id
	politicos.GET("/:id", politicoHandler.BuscarPorID)
	politicos.GET("/:id/estatisticas", politicoHandler.BuscarEstatisticas)
	politicos.GET("/:id/votacoes", politicoHandler.ListarVotacoes)
	politicos.GET("/:id/despesas", politicoHandler.ListarDespesas)
	politicos.GET("/:id/proposicoes", politicoHandler.ListarProposicoes)
	politicos.GET("/:id/presencas", politicoHandler.ListarPresencas)

	// Rotas de filtros
	filtros := api.Group("/filtros")
	filtros.GET("/partidos", filtrosHandler.ListarPartidos)
	filtros.GET("/estados", filtrosHandler.ListarEstados)
	filtros.GET("/cargos", filtrosHandler.ListarCargos)

	// Rotas de estatísticas
	estatisticas := api.Group("/estatisticas")
	estatisticas.GET("/geral", estatisticasHandler.Geral)
	estatisticas.GET("/ranking", estatisticasHandler.Ranking)

	// Rota de busca
	api.GET("/busca", politicoHandler.Buscar)

	// Iniciar servidor
	go func() {
		addr := ":" + cfg.Port
		log.Printf("🚀 Servidor iniciado em http://localhost%s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Encerrando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}

	log.Println("Servidor encerrado")
}
