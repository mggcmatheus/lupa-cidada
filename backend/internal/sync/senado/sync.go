package senado

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lupa-cidada/backend/internal/domain"
	"github.com/lupa-cidada/backend/internal/sync"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BaseURL = "https://legis.senado.leg.br/dadosabertos"
)

// SenadoSync sincroniza dados do Senado Federal
type SenadoSync struct {
	client *sync.HTTPClient
	db     *mongo.Database
}

// NewSenadoSync cria um novo sincronizador
func NewSenadoSync(db *mongo.Database) *SenadoSync {
	return &SenadoSync{
		client: sync.NewHTTPClient(3), // 3 requests por segundo (Senado é mais lento)
		db:     db,
	}
}

// SyncSenadores sincroniza todos os senadores em exercício
func (s *SenadoSync) SyncSenadores(ctx context.Context) error {
	log.Println("📥 Buscando senadores do Senado Federal...")

	url := fmt.Sprintf("%s/senador/lista/atual.json", BaseURL)

	var resp SenadoresResponse
	if err := s.client.Get(url, &resp); err != nil {
		return fmt.Errorf("erro ao buscar senadores: %w", err)
	}

	senadores := resp.ListaParlamentarEmExercicio.Parlamentares.Parlamentar
	log.Printf("📊 Total: %d senadores encontrados", len(senadores))

	for i, sen := range senadores {
		if err := s.syncSenador(ctx, sen); err != nil {
			log.Printf("⚠️  Erro ao sincronizar senador %s: %v",
				sen.IdentificacaoParlamentar.NomeParlamentar, err)
			continue
		}

		if (i+1)%20 == 0 {
			log.Printf("   Processados %d/%d senadores", i+1, len(senadores))
		}
	}

	log.Println("✅ Sincronização de senadores concluída!")
	return nil
}

// buscarPoliticoExistente busca um político existente por nome+data de nascimento
func (s *SenadoSync) buscarPoliticoExistente(ctx context.Context, nomeCivil string, dataNascimento time.Time) (*domain.Politico, error) {
	collection := s.db.Collection("politicos")

	if nomeCivil != "" && !dataNascimento.IsZero() {
		var politico domain.Politico
		err := collection.FindOne(ctx, bson.M{
			"nome_civil":      nomeCivil,
			"data_nascimento": dataNascimento,
		}).Decode(&politico)
		if err == nil {
			return &politico, nil
		}
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	return nil, nil // Não encontrado
}

// syncSenador sincroniza um senador específico
func (s *SenadoSync) syncSenador(ctx context.Context, sen Parlamentar) error {
	id := sen.IdentificacaoParlamentar

	// Buscar detalhes do senador
	url := fmt.Sprintf("%s/senador/%s.json", BaseURL, id.CodigoParlamentar)
	var detalhes SenadorDetalheResponse
	if err := s.client.Get(url, &detalhes); err != nil {
		// Se falhar, usar dados básicos
		log.Printf("   Usando dados básicos para %s", id.NomeParlamentar)
	}

	d := detalhes.DetalheParlamentar.Parlamentar
	dataNascimento := ParseDate(d.DadosBasicosParlamentar.DataNascimento)

	// Buscar se o político já existe no banco
	politicoExistente, err := s.buscarPoliticoExistente(ctx, id.NomeCompletoParlamentar, dataNascimento)
	if err != nil {
		return fmt.Errorf("erro ao buscar político existente: %w", err)
	}

	// Determinar data de início do mandato
	var dataInicio time.Time
	if sen.Mandato != nil && sen.Mandato.PrimeiraLegislatura != nil {
		dataInicio = ParseDate(sen.Mandato.PrimeiraLegislatura.DataInicio)
	}
	if dataInicio.IsZero() {
		dataInicio = time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
	}

	// Criar cargo do senador
	novoCargo := domain.CargoAtual{
		Tipo:        domain.CargoSenador,
		Esfera:      domain.EsferaFederal,
		Estado:      id.UfParlamentar,
		DataInicio:  dataInicio,
		EmExercicio: true, // Senado só retorna em exercício
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
			// Definir data de fim do cargo anterior
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

		// Atualizar cargo atual
		politico.CargoAtual = novoCargo

		// Atualizar outros dados se necessário
		if id.URLFotoParlamentar != "" {
			politico.FotoURL = id.URLFotoParlamentar
		}
		if id.EmailParlamentar != "" {
			politico.Contato.Email = id.EmailParlamentar
		}
		if d.DadosBasicosParlamentar.EnderecoParlamentar != "" {
			politico.Contato.Gabinete = d.DadosBasicosParlamentar.EnderecoParlamentar
		}
		politico.Partido = domain.Partido{
			Sigla: id.SiglaPartidoParlamentar,
			Nome:  "",
			Cor:   getPartidoCor(id.SiglaPartidoParlamentar),
		}
	} else {
		// Novo político - criar registro
		politico = domain.Politico{
			Nome:           id.NomeParlamentar,
			NomeCivil:      id.NomeCompletoParlamentar,
			FotoURL:        id.URLFotoParlamentar,
			DataNascimento: dataNascimento,
			Genero:         mapGenero(id.SexoParlamentar),
			Partido: domain.Partido{
				Sigla: id.SiglaPartidoParlamentar,
				Nome:  "",
				Cor:   getPartidoCor(id.SiglaPartidoParlamentar),
			},
			CargoAtual: novoCargo,
			Contato: domain.Contato{
				Email:    id.EmailParlamentar,
				Gabinete: d.DadosBasicosParlamentar.EnderecoParlamentar,
			},
			SalarioBruto:    41650.92, // Salário de senador
			SalarioLiquido:  30000.00,
			HistoricoCargos: []domain.CargoAtual{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		historicoCargos = []domain.CargoAtual{}
	}

	// Adicionar telefone se disponível
	if d.Telefones != nil && len(d.Telefones.Telefone) > 0 {
		politico.Contato.Telefone = d.Telefones.Telefone[0].NumeroTelefone
	}

	// Atualizar histórico
	politico.HistoricoCargos = historicoCargos
	politico.UpdatedAt = time.Now()

	// Upsert no MongoDB
	collection := s.db.Collection("politicos")
	var filter bson.M

	if politicoExistente != nil {
		// Atualizar existente
		filter = bson.M{"_id": politicoExistente.ID}
	} else {
		// Inserir novo
		filter = bson.M{
			"nome_civil":      id.NomeCompletoParlamentar,
			"data_nascimento": dataNascimento,
		}
	}

	update := bson.M{
		"$set": bson.M{
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
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// mapGenero converte o gênero da API para nosso modelo
func mapGenero(sexo string) domain.Genero {
	switch strings.ToUpper(sexo) {
	case "MASCULINO", "M":
		return domain.GeneroMasculino
	case "FEMININO", "F":
		return domain.GeneroFeminino
	default:
		return domain.GeneroOutro
	}
}

// getPartidoCor retorna a cor do partido
func getPartidoCor(sigla string) string {
	cores := map[string]string{
		"PT":            "#CC0000",
		"PL":            "#003366",
		"UNIÃO":         "#2E3092",
		"PP":            "#0066CC",
		"MDB":           "#00AA00",
		"PSD":           "#FF6600",
		"REPUBLICANOS":  "#0033CC",
		"PDT":           "#FF0000",
		"PSDB":          "#003399",
		"PSOL":          "#FFD700",
		"PSB":           "#FF6347",
		"PODE":          "#00CED1",
		"CIDADANIA":     "#9932CC",
		"AVANTE":        "#FF8C00",
		"SOLIDARIEDADE": "#FF4500",
		"PCdoB":         "#8B0000",
		"PV":            "#228B22",
		"NOVO":          "#FF6600",
		"REDE":          "#00AA66",
	}

	if cor, ok := cores[strings.ToUpper(sigla)]; ok {
		return cor
	}
	return "#666666"
}
