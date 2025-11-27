package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/lupa-cidada/backend/internal/sync/camara"
	"github.com/lupa-cidada/backend/internal/sync/senado"
	"github.com/lupa-cidada/backend/pkg/database"
)

func main() {
	// Flags
	mongoURI := flag.String("mongo", getEnv("MONGO_URI", "mongodb://lupa:lupa_secret_2024@localhost:27018/lupa_cidada?authSource=admin"), "MongoDB URI")
	syncCamara := flag.Bool("camara", false, "Sincronizar deputados da Câmara")
	syncSenado := flag.Bool("senado", false, "Sincronizar senadores do Senado")
	syncAll := flag.Bool("all", false, "Sincronizar tudo")
	flag.Parse()

	// Se nenhuma flag específica, sincronizar tudo
	if !*syncCamara && !*syncSenado {
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
