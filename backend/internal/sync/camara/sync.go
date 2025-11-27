package camara

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
	BaseURL = "https://dadosabertos.camara.leg.br/api/v2"
)

// CamaraSync sincroniza dados da Câmara dos Deputados
type CamaraSync struct {
	client *sync.HTTPClient
	db     *mongo.Database
}

// NewCamaraSync cria um novo sincronizador
func NewCamaraSync(db *mongo.Database) *CamaraSync {
	return &CamaraSync{
		client: sync.NewHTTPClient(5), // 5 requests por segundo
		db:     db,
	}
}

// SyncDeputados sincroniza todos os deputados
func (s *CamaraSync) SyncDeputados(ctx context.Context) error {
	log.Println("📥 Buscando deputados da Câmara...")

	// Buscar deputados da legislatura atual (57)
	url := fmt.Sprintf("%s/deputados?idLegislatura=57&itens=100&ordem=ASC&ordenarPor=nome", BaseURL)

	var allDeputados []DeputadoResumo

	for url != "" {
		var resp DeputadoResponse
		if err := s.client.Get(url, &resp); err != nil {
			return fmt.Errorf("erro ao buscar deputados: %w", err)
		}

		allDeputados = append(allDeputados, resp.Dados...)
		log.Printf("   Carregados %d deputados...", len(allDeputados))

		// Próxima página
		url = ""
		for _, link := range resp.Links {
			if link.Rel == "next" {
				url = link.Href
				break
			}
		}
	}

	log.Printf("📊 Total: %d deputados encontrados", len(allDeputados))

	// Processar cada deputado
	for i, dep := range allDeputados {
		if err := s.syncDeputado(ctx, dep); err != nil {
			log.Printf("⚠️  Erro ao sincronizar deputado %s: %v", dep.Nome, err)
			continue
		}

		if (i+1)%50 == 0 {
			log.Printf("   Processados %d/%d deputados", i+1, len(allDeputados))
		}
	}

	log.Println("✅ Sincronização de deputados concluída!")
	return nil
}

// syncDeputado sincroniza um deputado específico
func (s *CamaraSync) syncDeputado(ctx context.Context, dep DeputadoResumo) error {
	// Buscar detalhes do deputado
	url := fmt.Sprintf("%s/deputados/%d", BaseURL, dep.ID)
	var detalhes DeputadoDetalheResponse
	if err := s.client.Get(url, &detalhes); err != nil {
		return fmt.Errorf("erro ao buscar detalhes: %w", err)
	}

	d := detalhes.Dados

	// Mapear para nosso modelo
	politico := domain.Politico{
		Nome:           d.UltimoStatus.Nome,
		NomeCivil:      d.NomeCivil,
		FotoURL:        d.UltimoStatus.URLFoto,
		DataNascimento: ParseDate(d.DataNascimento),
		Genero:         mapGenero(d.Sexo),
		Partido: domain.Partido{
			Sigla: d.UltimoStatus.SiglaPartido,
			Nome:  "", // Será preenchido depois
			Cor:   getPartidoCor(d.UltimoStatus.SiglaPartido),
		},
		CargoAtual: domain.CargoAtual{
			Tipo:        domain.CargoDeputadoFederal,
			Esfera:      domain.EsferaFederal,
			Estado:      d.UltimoStatus.SiglaUF,
			DataInicio:  ParseDate(d.UltimoStatus.Data),
			EmExercicio: d.UltimoStatus.Situacao == "Exercício",
		},
		Contato: domain.Contato{
			Email: d.UltimoStatus.Email,
		},
		RedesSociais:   mapRedesSociais(d.RedeSocial),
		SalarioBruto:   33763.00, // Salário padrão de deputado federal
		SalarioLiquido: 25000.00,
		UpdatedAt:      time.Now(),
	}

	// Preencher gabinete se disponível
	if d.UltimoStatus.Gabinete != nil {
		g := d.UltimoStatus.Gabinete
		politico.Contato.Gabinete = fmt.Sprintf("%s, %s, Sala %s", g.Predio, g.Andar, g.Sala)
		if g.Telefone != "" {
			politico.Contato.Telefone = g.Telefone
		}
	}

	// Upsert no MongoDB (atualiza se existir, insere se não)
	collection := s.db.Collection("politicos")
	filter := bson.M{
		"nome":               politico.Nome,
		"cargo_atual.tipo":   domain.CargoDeputadoFederal,
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

// SyncDespesas sincroniza despesas dos deputados
func (s *CamaraSync) SyncDespesas(ctx context.Context, ano int) error {
	log.Printf("📥 Buscando despesas do ano %d...", ano)

	// Buscar todos os deputados do banco
	collection := s.db.Collection("politicos")
	cursor, err := collection.Find(ctx, bson.M{"cargo_atual.tipo": domain.CargoDeputadoFederal})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var politicos []domain.Politico
	if err := cursor.All(ctx, &politicos); err != nil {
		return err
	}

	despesasCollection := s.db.Collection("despesas")

	// TODO: Implementar sincronização de despesas por deputado
	// Por enquanto, apenas log
	log.Printf("   %d políticos encontrados para sincronizar despesas", len(politicos))
	_ = despesasCollection

	log.Println("✅ Sincronização de despesas concluída!")
	return nil
}

// mapGenero converte o gênero da API para nosso modelo
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

// mapRedesSociais extrai redes sociais da lista
func mapRedesSociais(redes []string) domain.RedesSociais {
	rs := domain.RedesSociais{}
	for _, url := range redes {
		url = strings.ToLower(url)
		switch {
		case strings.Contains(url, "twitter.com") || strings.Contains(url, "x.com"):
			// Extrair handle do Twitter
			parts := strings.Split(url, "/")
			if len(parts) > 0 {
				rs.Twitter = "@" + parts[len(parts)-1]
			}
		case strings.Contains(url, "instagram.com"):
			parts := strings.Split(url, "/")
			if len(parts) > 0 {
				rs.Instagram = "@" + parts[len(parts)-1]
			}
		case strings.Contains(url, "facebook.com"):
			parts := strings.Split(url, "/")
			if len(parts) > 0 {
				rs.Facebook = parts[len(parts)-1]
			}
		}
	}
	return rs
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
		"PRD":           "#1E90FF",
		"AGIR":          "#4169E1",
	}

	if cor, ok := cores[strings.ToUpper(sigla)]; ok {
		return cor
	}
	return "#666666"
}
