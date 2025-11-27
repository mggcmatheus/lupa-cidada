package presidente

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

// PresidenteSync sincroniza dados do Presidente da República
type PresidenteSync struct {
	db *mongo.Database
}

// NewPresidenteSync cria um novo sincronizador
func NewPresidenteSync(db *mongo.Database) *PresidenteSync {
	return &PresidenteSync{
		db: db,
	}
}

// SyncPresidente sincroniza o presidente atual
// Por enquanto, usa dados fixos que podem ser atualizados manualmente
// Futuramente: integrar com Portal da Transparência ou TSE
func (s *PresidenteSync) SyncPresidente(ctx context.Context) error {
	log.Println("📥 Sincronizando Presidente da República...")

	// Dados do presidente atual (atualizar conforme necessário)
	// Fonte: https://www.gov.br/planalto/
	presidenteData := PresidenteData{
		Nome:           "Luiz Inácio Lula da Silva",
		NomeCivil:      "Luiz Inácio Lula da Silva",
		DataNascimento: ParseDate("1945-10-27"),
		Genero:         "M",
		Partido:        "PT",
		Estado:         "SP",
		DataInicio:     ParseDate("2023-01-01"),
		EmExercicio:    true,
		FotoURL:        "https://www.gov.br/planalto/pt-br/acompanhe-o-planalto/fotos-do-presidente",
		Email:          "presidencia@planalto.gov.br",
	}

	return s.syncPresidente(ctx, presidenteData)
}

// syncPresidente sincroniza um presidente específico
func (s *PresidenteSync) syncPresidente(ctx context.Context, p PresidenteData) error {
	collection := s.db.Collection("politicos")

	// IMPORTANTE: Remover cargo de presidente de TODOS os outros políticos primeiro
	// (só pode haver um presidente em exercício por vez)
	filterOutrosPresidentes := bson.M{
		"cargo_atual.tipo": domain.CargoPresidente,
		"cargo_atual.em_exercicio": true,
	}
	
	// Buscar todos que têm cargo de presidente
	cursor, err := collection.Find(ctx, filterOutrosPresidentes)
	if err != nil {
		return fmt.Errorf("erro ao buscar outros presidentes: %w", err)
	}
	defer cursor.Close(ctx)

	var outrosPresidentes []domain.Politico
	if err := cursor.All(ctx, &outrosPresidentes); err != nil {
		return fmt.Errorf("erro ao ler outros presidentes: %w", err)
	}

	// Remover cargo de presidente de todos os outros (exceto se for a mesma pessoa)
	for _, outro := range outrosPresidentes {
		nomeOutro := strings.ToUpper(strings.TrimSpace(outro.NomeCivil))
		nomeNovo := strings.ToUpper(strings.TrimSpace(p.NomeCivil))
		
		// Se não é a mesma pessoa, remover cargo de presidente
		if !strings.Contains(nomeOutro, nomeNovo) && !strings.Contains(nomeNovo, nomeOutro) {
			log.Printf("⚠️  Removendo cargo de presidente de: %s", outro.Nome)
			
			cargoAnterior := outro.CargoAtual
			cargoAnterior.EmExercicio = false
			cargoAnterior.DataFim = time.Now()
			
			historicoAntigo := make([]domain.CargoAtual, len(outro.HistoricoCargos))
			copy(historicoAntigo, outro.HistoricoCargos)
			
			// Adicionar cargo de presidente ao histórico
			jaExiste := false
			for _, hc := range historicoAntigo {
				if hc.Tipo == cargoAnterior.Tipo && hc.DataInicio.Equal(cargoAnterior.DataInicio) {
					jaExiste = true
					break
				}
			}
			if !jaExiste {
				historicoAntigo = append(historicoAntigo, cargoAnterior)
			}
			
			// Restaurar último cargo do histórico que não seja presidente, ou deixar vazio
			novoCargoAntigo := domain.CargoAtual{}
			for i := len(historicoAntigo) - 1; i >= 0; i-- {
				if historicoAntigo[i].Tipo != domain.CargoPresidente {
					novoCargoAntigo = historicoAntigo[i]
					historicoAntigo = append(historicoAntigo[:i], historicoAntigo[i+1:]...)
					break
				}
			}
			
			collection.UpdateOne(ctx, bson.M{"_id": outro.ID}, bson.M{
				"$set": bson.M{
					"cargo_atual": novoCargoAntigo,
					"historico_cargos": historicoAntigo,
				},
			})
		}
	}

	// Agora buscar o político correto (pode ser o mesmo que já tinha cargo de presidente)
	var politicoExistente *domain.Politico
	filterPresidente := bson.M{
		"cargo_atual.tipo": domain.CargoPresidente,
		"cargo_atual.em_exercicio": true,
	}
	err = collection.FindOne(ctx, filterPresidente).Decode(&politicoExistente)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("erro ao buscar presidente existente: %w", err)
	}

	// Se encontrou presidente, verificar se é a mesma pessoa
	if politicoExistente != nil {
		nomeExistente := strings.ToUpper(strings.TrimSpace(politicoExistente.NomeCivil))
		nomeNovo := strings.ToUpper(strings.TrimSpace(p.NomeCivil))
		
		// Se os nomes são muito diferentes, buscar novamente
		if !strings.Contains(nomeExistente, nomeNovo) && !strings.Contains(nomeNovo, nomeExistente) {
			// Presidente diferente, buscar pelo nome correto
			politicoExistente = nil
		}
	}

	// Se não encontrou por cargo ou é pessoa diferente, buscar por CPF ou nome+data
	if politicoExistente == nil {
		filter := bson.M{}
		
		// Se tem CPF, buscar por CPF
		if p.CPF != "" {
			filter["cpf"] = p.CPF
		} else {
			// Buscar por nome civil exato (case insensitive) E data de nascimento exata
			// Usar regex para match exato do nome
			nomeRegex := "^" + strings.ReplaceAll(strings.ToUpper(p.NomeCivil), " ", "\\s+") + "$"
			filter["$and"] = []bson.M{
				{"nome_civil": bson.M{"$regex": nomeRegex, "$options": "i"}},
				{"data_nascimento": bson.M{
					"$gte": p.DataNascimento.AddDate(0, 0, -1),
					"$lte": p.DataNascimento.AddDate(0, 0, 1),
				}},
			}
		}

		err = collection.FindOne(ctx, filter).Decode(&politicoExistente)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("erro ao buscar político existente: %w", err)
		}
	}

	// Criar cargo do presidente
	novoCargo := domain.CargoAtual{
		Tipo:        domain.CargoPresidente,
		Esfera:      domain.EsferaFederal,
		Estado:      "BR", // Brasil
		DataInicio:  p.DataInicio,
		EmExercicio: p.EmExercicio,
	}

	if !p.EmExercicio && !p.DataFim.IsZero() {
		novoCargo.DataFim = p.DataFim
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
					hc.DataInicio.Equal(cargoAnterior.DataInicio) {
					jaExiste = true
					break
				}
			}
			if !jaExiste {
				historicoCargos = append(historicoCargos, cargoAnterior)
			}
		}

		// Atualizar cargo atual
		politico.CargoAtual = novoCargo

		// Atualizar TODOS os dados do presidente (não apenas se não vazio)
		politico.Nome = p.Nome
		politico.NomeCivil = p.NomeCivil
		if p.FotoURL != "" {
			politico.FotoURL = p.FotoURL
		}
		if p.Email != "" {
			politico.Contato.Email = p.Email
		}
		if p.Telefone != "" {
			politico.Contato.Telefone = p.Telefone
		}
		if p.CPF != "" {
			politico.CPF = p.CPF
		}
		politico.Partido = domain.Partido{
			Sigla: p.Partido,
			Nome:  "",
			Cor:   getPartidoCor(p.Partido),
		}
		politico.DataNascimento = p.DataNascimento
		politico.Genero = mapGenero(p.Genero)
	} else {
		// Novo político - criar registro
		politico = domain.Politico{
			CPF:             p.CPF,
			Nome:            p.Nome,
			NomeCivil:       p.NomeCivil,
			FotoURL:         p.FotoURL,
			DataNascimento:  p.DataNascimento,
			Genero:          mapGenero(p.Genero),
			Partido: domain.Partido{
				Sigla: p.Partido,
				Nome:  "",
				Cor:   getPartidoCor(p.Partido),
			},
			CargoAtual: novoCargo,
			Contato: domain.Contato{
				Email:    p.Email,
				Telefone: p.Telefone,
			},
			SalarioBruto:    41000.00, // Salário do presidente
			SalarioLiquido:  30000.00,
			HistoricoCargos: []domain.CargoAtual{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		historicoCargos = []domain.CargoAtual{}
	}

	politico.HistoricoCargos = historicoCargos
	politico.UpdatedAt = time.Now()

	// Upsert no banco
	// Se já existe um presidente, usar o ID dele, senão buscar por CPF ou nome+data
	var updateFilter bson.M
	if politicoExistente != nil {
		updateFilter = bson.M{"_id": politicoExistente.ID}
	} else {
		updateFilter = bson.M{}
		if p.CPF != "" {
			updateFilter["cpf"] = p.CPF
		} else {
			updateFilter["$and"] = []bson.M{
				{"nome_civil": bson.M{"$regex": "^" + strings.ToUpper(p.NomeCivil) + "$", "$options": "i"}},
				{"data_nascimento": bson.M{
					"$gte": p.DataNascimento.AddDate(0, 0, -1),
					"$lte": p.DataNascimento.AddDate(0, 0, 1),
				}},
			}
		}
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
		return fmt.Errorf("erro ao salvar presidente: %w", err)
	}

	log.Printf("✅ Presidente sincronizado: %s", p.Nome)
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

