package governadores

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GovernadoresSync sincroniza dados dos Governadores
type GovernadoresSync struct {
	db *mongo.Database
}

// NewGovernadoresSync cria um novo sincronizador
func NewGovernadoresSync(db *mongo.Database) *GovernadoresSync {
	return &GovernadoresSync{
		db: db,
	}
}

// SyncGovernadores sincroniza todos os governadores
// Por enquanto, usa dados fixos que podem ser atualizados manualmente
// Futuramente: integrar com Portal da Transparência ou TSE
func (s *GovernadoresSync) SyncGovernadores(ctx context.Context) error {
	log.Println("📥 Sincronizando Governadores dos Estados...")

	// Lista de governadores atuais (atualizar conforme necessário)
	// Fonte: https://www.gov.br/planalto/ e sites dos governos estaduais
	governadores := []GovernadorData{
		// Adicionar governadores aqui
		// Exemplo (atualizar com dados reais):
		// {Nome: "Nome do Governador", Estado: "SP", Partido: "PT", ...},
	}

	// Se a lista estiver vazia, logar aviso
	if len(governadores) == 0 {
		log.Println("⚠️  Lista de governadores vazia. Adicione os dados em sync.go")
		log.Println("   Fonte: https://www.gov.br/planalto/ ou sites dos governos estaduais")
		return nil
	}

	log.Printf("📊 Total: %d governadores para sincronizar", len(governadores))

	for i, gov := range governadores {
		if err := s.syncGovernador(ctx, gov); err != nil {
			log.Printf("⚠️  Erro ao sincronizar governador %s (%s): %v", gov.Nome, gov.Estado, err)
			continue
		}

		if (i+1)%10 == 0 {
			log.Printf("   Processados %d/%d governadores", i+1, len(governadores))
		}
	}

	log.Println("✅ Sincronização de governadores concluída!")
	return nil
}

// syncGovernador sincroniza um governador específico
func (s *GovernadoresSync) syncGovernador(ctx context.Context, g GovernadorData) error {
	collection := s.db.Collection("politicos")

	// Buscar se o político já existe
	filter := bson.M{
		"$or": []bson.M{
			{"cpf": g.CPF},
			{
				"nome_civil": strings.ToUpper(g.NomeCivil),
				"data_nascimento": bson.M{
					"$gte": g.DataNascimento.AddDate(0, 0, -1),
					"$lte": g.DataNascimento.AddDate(0, 0, 1),
				},
			},
		},
	}

	var politicoExistente *domain.Politico
	err := collection.FindOne(ctx, filter).Decode(&politicoExistente)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("erro ao buscar político existente: %w", err)
	}

	// Criar cargo do governador
	novoCargo := domain.CargoAtual{
		Tipo:        domain.CargoGovernador,
		Esfera:      domain.EsferaEstadual,
		Estado:      g.Estado,
		DataInicio:  g.DataInicio,
		EmExercicio: g.EmExercicio,
	}

	if !g.EmExercicio && !g.DataFim.IsZero() {
		novoCargo.DataFim = g.DataFim
	}

	var politico domain.Politico
	var historicoCargos []domain.CargoAtual

	if politicoExistente != nil {
		// Político já existe - atualizar dados e gerenciar histórico
		politico = *politicoExistente

		// Copiar histórico existente
		historicoCargos = make([]domain.CargoAtual, len(politico.HistoricoCargos))
		copy(historicoCargos, politico.HistoricoCargos)

		// Se o cargo atual mudou, mover para histórico
		cargoAtualMudou := politico.CargoAtual.Tipo != novoCargo.Tipo ||
			politico.CargoAtual.Estado != novoCargo.Estado ||
			!politico.CargoAtual.DataInicio.Equal(novoCargo.DataInicio)

		if cargoAtualMudou && politico.CargoAtual.Tipo != "" {
			cargoAnterior := politico.CargoAtual
			if cargoAnterior.DataFim.IsZero() {
				cargoAnterior.DataFim = time.Now()
			}
			cargoAnterior.EmExercicio = false

			// Adicionar ao histórico (evitar duplicatas)
			jaExiste := false
			for _, hc := range historicoCargos {
				if hc.Tipo == cargoAnterior.Tipo &&
					hc.Estado == cargoAnterior.Estado &&
					hc.DataInicio.Equal(cargoAnterior.DataInicio) {
					jaExiste = true
					break
				}
			}
			if !jaExiste {
				historicoCargos = append(historicoCargos, cargoAnterior)
			}
		}

		// Atualizar cargo atual apenas se estiver em exercício
		if novoCargo.EmExercicio {
			politico.CargoAtual = novoCargo
		} else {
			// Se não está em exercício, adicionar ao histórico
			jaExiste := false
			for _, hc := range historicoCargos {
				if hc.Tipo == novoCargo.Tipo &&
					hc.Estado == novoCargo.Estado &&
					hc.DataInicio.Equal(novoCargo.DataInicio) {
					jaExiste = true
					break
				}
			}
			if !jaExiste {
				historicoCargos = append(historicoCargos, novoCargo)
			}
		}

		// Atualizar outros dados
		if g.FotoURL != "" {
			politico.FotoURL = g.FotoURL
		}
		if g.Email != "" {
			politico.Contato.Email = g.Email
		}
		if g.Telefone != "" {
			politico.Contato.Telefone = g.Telefone
		}
		politico.Partido = domain.Partido{
			Sigla: g.Partido,
			Nome:  "",
			Cor:   getPartidoCor(g.Partido),
		}
	} else {
		// Novo político - criar registro
		politico = domain.Politico{
			CPF:             g.CPF,
			Nome:            g.Nome,
			NomeCivil:       g.NomeCivil,
			FotoURL:         g.FotoURL,
			DataNascimento:  g.DataNascimento,
			Genero:          mapGenero(g.Genero),
			Partido: domain.Partido{
				Sigla: g.Partido,
				Nome:  "",
				Cor:   getPartidoCor(g.Partido),
			},
			CargoAtual: novoCargo,
			Contato: domain.Contato{
				Email:    g.Email,
				Telefone: g.Telefone,
			},
			SalarioBruto:    35000.00, // Salário de governador (varia por estado)
			SalarioLiquido:  25000.00,
			HistoricoCargos: []domain.CargoAtual{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		historicoCargos = []domain.CargoAtual{}
	}

	politico.HistoricoCargos = historicoCargos
	politico.UpdatedAt = time.Now()

	// Upsert no banco
	updateFilter := bson.M{
		"$or": []bson.M{
			{"cpf": g.CPF},
			{
				"nome_civil": strings.ToUpper(g.NomeCivil),
				"data_nascimento": bson.M{
					"$gte": g.DataNascimento.AddDate(0, 0, -1),
					"$lte": g.DataNascimento.AddDate(0, 0, 1),
				},
			},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"cpf":              politico.CPF,
			"nome":             politico.Nome,
			"nome_civil":       politico.NomeCivil,
			"foto_url":         politico.FotoURL,
			"data_nascimento":  politico.DataNascimento,
			"genero":           politico.Genero,
			"partido":          politico.Partido,
			"cargo_atual":      politico.CargoAtual,
			"historico_cargos": politico.HistoricoCargos,
			"contato":          politico.Contato,
			"redes_sociais":    politico.RedesSociais,
			"salario_bruto":    politico.SalarioBruto,
			"salario_liquido":  politico.SalarioLiquido,
			"updated_at":       time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":        primitive.NewObjectID(),
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, updateFilter, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao salvar governador: %w", err)
	}

	log.Printf("✅ Governador sincronizado: %s (%s)", g.Nome, g.Estado)
	return nil
}

// mapGenero converte o gênero para nosso modelo
func mapGenero(sexo string) domain.Genero {
	switch strings.ToUpper(sexo) {
	case "M":
		return domain.GeneroMasculino
	case "F":
		return domain.GeneroFeminino
	default:
		return domain.GeneroOutro
	}
}

// getPartidoCor retorna a cor do partido (simplificado)
func getPartidoCor(sigla string) string {
	cores := map[string]string{
		"PT":   "#E31E24",
		"PSDB": "#003399",
		"PSB":  "#FF6600",
		"PMDB": "#0066CC",
		"DEM":  "#0066CC",
		"PP":   "#0066CC",
		"PSD":  "#FF6600",
		"PR":   "#0066CC",
		"PDT":  "#FF0000",
		"PTB":  "#0066CC",
		"PPS":  "#FF6600",
		"PV":   "#00AA00",
		"PCdoB": "#E31E24",
		"PSC":  "#0066CC",
		"PMN":  "#0066CC",
		"PRP":  "#0066CC",
		"PHS":  "#0066CC",
		"PRTB": "#0066CC",
		"PTC":  "#0066CC",
		"PSL":  "#0066CC",
		"REDE": "#00AA00",
		"MDB":  "#00AA00",
		"NOVO": "#FF6600",
		"PL":   "#003366",
		"UNIÃO": "#0066CC",
	}

	if cor, ok := cores[strings.ToUpper(sigla)]; ok {
		return cor
	}
	return "#666666"
}

