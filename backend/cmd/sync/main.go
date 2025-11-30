package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/lupa-cidada/backend/internal/sync/camara"
	"github.com/lupa-cidada/backend/internal/sync/governadores"
	"github.com/lupa-cidada/backend/internal/sync/presidente"
	"github.com/lupa-cidada/backend/internal/sync/senado"
	"github.com/lupa-cidada/backend/pkg/database"
)

func main() {
	// Flags
	mongoURI := flag.String("mongo", getEnv("MONGO_URI", "mongodb://lupa:lupa_secret_2024@localhost:27018/lupa_cidada?authSource=admin"), "MongoDB URI")
	syncCamara := flag.Bool("camara", false, "Sincronizar deputados da Câmara")
	syncSenado := flag.Bool("senado", false, "Sincronizar senadores do Senado")
	syncPresidente := flag.Bool("presidente", false, "Sincronizar Presidente da República")
	syncGovernadores := flag.Bool("governadores", false, "Sincronizar Governadores")
	syncVotacoes := flag.Bool("votacoes", false, "Sincronizar votações da Câmara")
	syncProposicoes := flag.Bool("proposicoes", false, "Sincronizar proposições da Câmara")
	syncDespesas := flag.Bool("despesas", false, "Sincronizar despesas da Câmara")
	syncPresencas := flag.Bool("presencas", false, "Sincronizar presenças em eventos da Câmara")
	ano := flag.Int("ano", time.Now().Year(), "Ano para sincronização de votações, proposições, despesas e presenças")
	syncAll := flag.Bool("all", false, "Sincronizar tudo")
	flag.Parse()

	// Se nenhuma flag específica, sincronizar tudo
	if !*syncCamara && !*syncSenado && !*syncPresidente && !*syncGovernadores {
		*syncAll = true
	}

	log.Println("🔍 Lupa Cidadã - Sincronização de Dados")
	log.Println("========================================")

	// Conectar ao MongoDB
	log.Println("📦 Conectando ao MongoDB...")
	client, err := database.NewMongoClient(*mongoURI)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("lupa_cidada")
	log.Println("✅ Conectado ao MongoDB!")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	start := time.Now()

	// Sincronizar Câmara
	if *syncAll || *syncCamara {
		log.Println("")
		log.Println("🏛️  CÂMARA DOS DEPUTADOS")
		log.Println("------------------------")

		camaraSync := camara.NewCamaraSync(db)
		if err := camaraSync.SyncDeputados(ctx); err != nil {
			log.Printf("❌ Erro na sincronização da Câmara: %v", err)
		}
	}

	// Sincronizar dados adicionais da Câmara
	if *syncAll || *syncVotacoes || *syncProposicoes || *syncDespesas || *syncPresencas {
		camaraSync := camara.NewCamaraSync(db)

		if *syncAll || *syncVotacoes {
			log.Println("")
			log.Println("📊 VOTAÇÕES DA CÂMARA")
			log.Println("---------------------")
			if err := camaraSync.SyncVotacoes(ctx, *ano); err != nil {
				log.Printf("❌ Erro na sincronização de votações: %v", err)
			}
		}

		if *syncAll || *syncProposicoes {
			log.Println("")
			log.Println("📄 PROPOSIÇÕES DA CÂMARA")
			log.Println("------------------------")
			if err := camaraSync.SyncProposicoes(ctx, *ano); err != nil {
				log.Printf("❌ Erro na sincronização de proposições: %v", err)
			}
		}

		if *syncAll || *syncDespesas {
			log.Println("")
			log.Println("💰 DESPESAS DA CÂMARA")
			log.Println("---------------------")
			if err := camaraSync.SyncDespesas(ctx, *ano); err != nil {
				log.Printf("❌ Erro na sincronização de despesas: %v", err)
			}
		}

		if *syncAll || *syncPresencas {
			log.Println("")
			log.Println("✅ PRESENÇAS EM EVENTOS DA CÂMARA")
			log.Println("----------------------------------")
			if err := camaraSync.SyncPresencas(ctx, *ano); err != nil {
				log.Printf("❌ Erro na sincronização de presenças: %v", err)
			}
		}
	}

	// Sincronizar Senado
	if *syncAll || *syncSenado {
		log.Println("")
		log.Println("🏛️  SENADO FEDERAL")
		log.Println("------------------")

		senadoSync := senado.NewSenadoSync(db)
		if err := senadoSync.SyncSenadores(ctx); err != nil {
			log.Printf("❌ Erro na sincronização do Senado: %v", err)
		}
	}

	// Sincronizar Presidente
	if *syncAll || *syncPresidente {
		log.Println("")
		log.Println("🇧🇷 PRESIDÊNCIA DA REPÚBLICA")
		log.Println("----------------------------")

		presidenteSync := presidente.NewPresidenteSync(db)
		if err := presidenteSync.SyncPresidente(ctx); err != nil {
			log.Printf("❌ Erro na sincronização do Presidente: %v", err)
		}
	}

	// Sincronizar Governadores
	if *syncAll || *syncGovernadores {
		log.Println("")
		log.Println("🏛️  GOVERNADORES DOS ESTADOS")
		log.Println("----------------------------")

		governadoresSync := governadores.NewGovernadoresSync(db)
		if err := governadoresSync.SyncGovernadores(ctx); err != nil {
			log.Printf("❌ Erro na sincronização dos Governadores: %v", err)
		}
	}

	// Estatísticas finais
	log.Println("")
	log.Println("========================================")
	log.Printf("⏱️  Tempo total: %s", time.Since(start).Round(time.Second))

	// Contar registros
	countPoliticos, _ := db.Collection("politicos").CountDocuments(ctx, map[string]interface{}{})
	log.Printf("📊 Total de políticos no banco: %d", countPoliticos)

	log.Println("✅ Sincronização concluída!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
