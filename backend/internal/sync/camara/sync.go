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

// buscarPoliticoExistente busca um político existente por CPF ou nome+data de nascimento
func (s *CamaraSync) buscarPoliticoExistente(ctx context.Context, cpf string, nomeCivil string, dataNascimento time.Time) (*domain.Politico, error) {
	collection := s.db.Collection("politicos")

	// Tentar buscar por CPF primeiro (mais confiável)
	if cpf != "" {
		var politico domain.Politico
		err := collection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&politico)
		if err == nil {
			return &politico, nil
		}
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	// Se não encontrou por CPF, tentar por nome civil + data de nascimento
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

// syncDeputado sincroniza um deputado específico
func (s *CamaraSync) syncDeputado(ctx context.Context, dep DeputadoResumo) error {
	// Buscar detalhes do deputado
	url := fmt.Sprintf("%s/deputados/%d", BaseURL, dep.ID)
	var detalhes DeputadoDetalheResponse
	if err := s.client.Get(url, &detalhes); err != nil {
		return fmt.Errorf("erro ao buscar detalhes: %w", err)
	}

	d := detalhes.Dados
	dataNascimento := ParseDate(d.DataNascimento)

	// Buscar se o político já existe no banco
	politicoExistente, err := s.buscarPoliticoExistente(ctx, d.CPF, d.NomeCivil, dataNascimento)
	if err != nil {
		return fmt.Errorf("erro ao buscar político existente: %w", err)
	}

	// Criar cargo do deputado
	novoCargo := domain.CargoAtual{
		Tipo:        domain.CargoDeputadoFederal,
		Esfera:      domain.EsferaFederal,
		Estado:      d.UltimoStatus.SiglaUF,
		DataInicio:  ParseDate(d.UltimoStatus.Data),
		EmExercicio: d.UltimoStatus.Situacao == "Exercício",
	}

	// Se não está em exercício, definir data de fim
	if !novoCargo.EmExercicio {
		novoCargo.DataFim = time.Now()
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

		// Atualizar outros dados se necessário
		if d.UltimoStatus.URLFoto != "" {
			politico.FotoURL = d.UltimoStatus.URLFoto
		}
		if d.UltimoStatus.Email != "" {
			politico.Contato.Email = d.UltimoStatus.Email
		}
		politico.RedesSociais = mapRedesSociais(d.RedeSocial)
		politico.Partido = domain.Partido{
			Sigla: d.UltimoStatus.SiglaPartido,
			Nome:  "",
			Cor:   getPartidoCor(d.UltimoStatus.SiglaPartido),
		}
	} else {
		// Novo político - criar registro
		politico = domain.Politico{
			CPF:            d.CPF,
			Nome:           d.UltimoStatus.Nome,
			NomeCivil:      d.NomeCivil,
			FotoURL:        d.UltimoStatus.URLFoto,
			DataNascimento: dataNascimento,
			Genero:         mapGenero(d.Sexo),
			Partido: domain.Partido{
				Sigla: d.UltimoStatus.SiglaPartido,
				Nome:  "",
				Cor:   getPartidoCor(d.UltimoStatus.SiglaPartido),
			},
			Contato: domain.Contato{
				Email: d.UltimoStatus.Email,
			},
			RedesSociais:   mapRedesSociais(d.RedeSocial),
			SalarioBruto:   33763.00,
			SalarioLiquido: 25000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Se está em exercício, definir como cargo atual
		// Se não está, adicionar ao histórico
		if novoCargo.EmExercicio {
			politico.CargoAtual = novoCargo
			historicoCargos = []domain.CargoAtual{}
		} else {
			politico.CargoAtual = domain.CargoAtual{} // Vazio até ter outro cargo
			historicoCargos = []domain.CargoAtual{novoCargo}
		}
	}

	// Preencher gabinete se disponível e estiver em exercício
	if novoCargo.EmExercicio && d.UltimoStatus.Gabinete != nil {
		g := d.UltimoStatus.Gabinete
		politico.Contato.Gabinete = fmt.Sprintf("%s, %s, Sala %s", g.Predio, g.Andar, g.Sala)
		if g.Telefone != "" {
			politico.Contato.Telefone = g.Telefone
		}
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
			"$or": []bson.M{
				{"cpf": d.CPF},
				{
					"nome_civil":      d.NomeCivil,
					"data_nascimento": dataNascimento,
				},
			},
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
	_, err = collection.UpdateOne(ctx, filter, update, opts)
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
