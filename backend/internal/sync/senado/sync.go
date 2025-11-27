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

	// Determinar data de início do mandato
	var dataInicio time.Time
	if sen.Mandato != nil && sen.Mandato.PrimeiraLegislatura != nil {
		dataInicio = ParseDate(sen.Mandato.PrimeiraLegislatura.DataInicio)
	}
	if dataInicio.IsZero() {
		dataInicio = time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
	}

	// Mapear para nosso modelo
	politico := domain.Politico{
		Nome:           id.NomeParlamentar,
		NomeCivil:      id.NomeCompletoParlamentar,
		FotoURL:        id.URLFotoParlamentar,
		DataNascimento: ParseDate(d.DadosBasicosParlamentar.DataNascimento),
		Genero:         mapGenero(id.SexoParlamentar),
		Partido: domain.Partido{
			Sigla: id.SiglaPartidoParlamentar,
			Nome:  "",
			Cor:   getPartidoCor(id.SiglaPartidoParlamentar),
		},
		CargoAtual: domain.CargoAtual{
			Tipo:        domain.CargoSenador,
			Esfera:      domain.EsferaFederal,
			Estado:      id.UfParlamentar,
			DataInicio:  dataInicio,
			EmExercicio: true,
		},
		Contato: domain.Contato{
			Email:    id.EmailParlamentar,
			Gabinete: d.DadosBasicosParlamentar.EnderecoParlamentar,
		},
		SalarioBruto:   41650.92, // Salário de senador
		SalarioLiquido: 30000.00,
		UpdatedAt:      time.Now(),
	}

	// Adicionar telefone se disponível
	if d.Telefones != nil && len(d.Telefones.Telefone) > 0 {
		politico.Contato.Telefone = d.Telefones.Telefone[0].NumeroTelefone
	}

	// Upsert no MongoDB
	collection := s.db.Collection("politicos")
	filter := bson.M{
		"nome":               politico.Nome,
		"cargo_atual.tipo":   domain.CargoSenador,
		"cargo_atual.estado": politico.CargoAtual.Estado,
	}

	// Usar $set para campos que sempre atualizam
	update := bson.M{
		"$set": bson.M{
			"nome":            politico.Nome,
			"nome_civil":      politico.NomeCivil,
			"foto_url":        politico.FotoURL,
			"data_nascimento": politico.DataNascimento,
			"genero":          politico.Genero,
			"partido":         politico.Partido,
			"cargo_atual":     politico.CargoAtual,
			"contato":         politico.Contato,
			"redes_sociais":   politico.RedesSociais,
			"salario_bruto":   politico.SalarioBruto,
			"salario_liquido": politico.SalarioLiquido,
			"updated_at":      time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":        primitive.NewObjectID(),
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
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
