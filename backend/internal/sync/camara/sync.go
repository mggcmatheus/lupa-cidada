package camara

import (
	"context"
	"fmt"
	"log"
	"strings"
	syncpkg "sync"
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
		client: sync.NewHTTPClient(15), // 15 requests por segundo (aumentado para acelerar)
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

	// Processar deputados em paralelo com worker pool
	const numWorkers = 10 // Número de goroutines simultâneas
	deputadosChan := make(chan DeputadoResumo, numWorkers)
	var wg syncpkg.WaitGroup
	var mu syncpkg.Mutex
	processed := 0
	errors := 0

	// Iniciar workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for dep := range deputadosChan {
				if err := s.syncDeputado(ctx, dep); err != nil {
					mu.Lock()
					errors++
					mu.Unlock()
					log.Printf("⚠️  Erro ao sincronizar deputado %s: %v", dep.Nome, err)
					continue
				}

				mu.Lock()
				processed++
				if processed%50 == 0 {
					log.Printf("   Processados %d/%d deputados", processed, len(allDeputados))
				}
				mu.Unlock()
			}
		}()
	}

	// Enviar deputados para processamento
	for _, dep := range allDeputados {
		deputadosChan <- dep
	}
	close(deputadosChan)

	// Aguardar conclusão
	wg.Wait()

	if errors > 0 {
		log.Printf("⚠️  %d erros durante a sincronização", errors)
	}

	log.Printf("✅ Sincronização de deputados concluída! (%d processados)", processed)
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
		if d.UltimoStatus.NomeEleitoral != "" {
			politico.NomeEleitoral = d.UltimoStatus.NomeEleitoral
		}
		if d.URLWebsite != "" {
			politico.Website = d.URLWebsite
		}
		if d.Escolaridade != "" {
			politico.Escolaridade = d.Escolaridade
		}
		if d.MunicipioNascimento != "" {
			politico.MunicipioNascimento = d.MunicipioNascimento
		}
		if d.UfNascimento != "" {
			politico.UFNascimento = d.UfNascimento
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
			NomeEleitoral:  d.UltimoStatus.NomeEleitoral,
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
			RedesSociais:        mapRedesSociais(d.RedeSocial),
			SalarioBruto:        33763.00,
			SalarioLiquido:      25000.00,
			Escolaridade:        d.Escolaridade,
			MunicipioNascimento: d.MunicipioNascimento,
			UFNascimento:        d.UfNascimento,
			Website:             d.URLWebsite,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
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
			"cpf":                  politico.CPF,
			"nome":                 politico.Nome,
			"nome_civil":           politico.NomeCivil,
			"nome_eleitoral":       politico.NomeEleitoral,
			"foto_url":             politico.FotoURL,
			"data_nascimento":      politico.DataNascimento,
			"genero":               politico.Genero,
			"partido":              politico.Partido,
			"cargo_atual":          politico.CargoAtual,
			"historico_cargos":     politico.HistoricoCargos,
			"contato":              politico.Contato,
			"redes_sociais":        politico.RedesSociais,
			"salario_bruto":        politico.SalarioBruto,
			"salario_liquido":      politico.SalarioLiquido,
			"escolaridade":         politico.Escolaridade,
			"municipio_nascimento": politico.MunicipioNascimento,
			"uf_nascimento":        politico.UFNascimento,
			"website":              politico.Website,
			"id_externo_camara":    d.ID,
			"updated_at":           time.Now(),
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

	log.Printf("   %d políticos encontrados para sincronizar despesas", len(politicos))

	// Processar despesas em paralelo
	const numWorkers = 5
	politicosChan := make(chan domain.Politico, numWorkers)
	var wg syncpkg.WaitGroup
	var mu syncpkg.Mutex
	processed := 0

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for politico := range politicosChan {
				if politico.IDExternoCamara == 0 {
					continue
				}
				if err := s.SyncDespesasPorDeputado(ctx, politico.IDExternoCamara, ano); err != nil {
					log.Printf("⚠️  Erro ao sincronizar despesas do deputado %s: %v", politico.Nome, err)
				}
				mu.Lock()
				processed++
				if processed%10 == 0 {
					log.Printf("   Processando despesas: %d/%d", processed, len(politicos))
				}
				mu.Unlock()
			}
		}()
	}

	for _, politico := range politicos {
		if politico.IDExternoCamara != 0 {
			politicosChan <- politico
		}
	}
	close(politicosChan)
	wg.Wait()

	log.Println("✅ Sincronização de despesas concluída!")
	return nil
}

// SyncDespesasPorDeputado sincroniza despesas de um deputado específico
func (s *CamaraSync) SyncDespesasPorDeputado(ctx context.Context, deputadoID int, ano int) error {
	url := fmt.Sprintf("%s/deputados/%d/despesas?ano=%d&itens=100&ordem=ASC&ordenarPor=mes", BaseURL, deputadoID, ano)
	despesasCollection := s.db.Collection("despesas")

	var allDespesas []Despesa
	for url != "" {
		var resp DespesasResponse
		if err := s.client.Get(url, &resp); err != nil {
			return fmt.Errorf("erro ao buscar despesas: %w", err)
		}

		allDespesas = append(allDespesas, resp.Dados...)

		// Próxima página
		url = ""
		for _, link := range resp.Links {
			if link.Rel == "next" {
				url = link.Href
				break
			}
		}
	}

	// Buscar político pelo ID externo da API
	var politico domain.Politico
	err := s.db.Collection("politicos").FindOne(ctx, bson.M{"id_externo_camara": deputadoID}).Decode(&politico)
	if err != nil {
		return fmt.Errorf("político não encontrado com ID externo %d: %w", deputadoID, err)
	}

	for _, despesa := range allDespesas {
		dataDoc := ParseDate(despesa.DataDocumento)
		if dataDoc.IsZero() {
			dataDoc = time.Date(despesa.Ano, time.Month(despesa.Mes), 1, 0, 0, 0, 0, time.UTC)
		}

		despesaDoc := domain.Despesa{
			PoliticoID:     politico.ID,
			Tipo:           despesa.TipoDespesa,
			Descricao:      fmt.Sprintf("%s - %s", despesa.TipoDespesa, despesa.NomeFornecedor),
			Fornecedor:     despesa.NomeFornecedor,
			CNPJFornecedor: despesa.CNPJCPFFornecedor,
			Valor:          despesa.ValorLiquido,
			Data:           dataDoc,
			MesReferencia:  despesa.Mes,
			AnoReferencia:  despesa.Ano,
			DocumentoURL:   despesa.URLDocumento,
		}

		// Upsert despesa
		filter := bson.M{
			"politico_id":     politico.ID,
			"ano_referencia":  despesa.Ano,
			"mes_referencia":  despesa.Mes,
			"tipo":            despesa.TipoDespesa,
			"cnpj_fornecedor": despesa.CNPJCPFFornecedor,
			"valor":           despesa.ValorLiquido,
		}

		update := bson.M{
			"$set": bson.M{
				"descricao":     despesaDoc.Descricao,
				"fornecedor":    despesaDoc.Fornecedor,
				"data":          despesaDoc.Data,
				"documento_url": despesaDoc.DocumentoURL,
				"updated_at":    time.Now(),
			},
			"$setOnInsert": bson.M{
				"_id":        primitive.NewObjectID(),
				"created_at": time.Now(),
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err := despesasCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("⚠️  Erro ao salvar despesa: %v", err)
			continue
		}
	}

	return nil
}

// SyncVotacoes sincroniza votações dos deputados
func (s *CamaraSync) SyncVotacoes(ctx context.Context, ano int) error {
	log.Printf("📥 Buscando votações do ano %d...", ano)

	// Buscar todas as votações do ano
	url := fmt.Sprintf("%s/votacoes?ano=%d&itens=100&ordem=ASC&ordenarPor=data", BaseURL, ano)
	votacoesCollection := s.db.Collection("votacoes")
	proposicoesCollection := s.db.Collection("proposicoes")

	var allVotacoes []Votacao
	for url != "" {
		var resp VotacoesResponse
		if err := s.client.Get(url, &resp); err != nil {
			return fmt.Errorf("erro ao buscar votações: %w", err)
		}

		allVotacoes = append(allVotacoes, resp.Dados...)
		log.Printf("   Carregadas %d votações...", len(allVotacoes))

		// Próxima página
		url = ""
		for _, link := range resp.Links {
			if link.Rel == "next" {
				url = link.Href
				break
			}
		}
	}

	log.Printf("📊 Total: %d votações encontradas", len(allVotacoes))

	// Processar votações em paralelo
	const numWorkers = 5
	votacoesChan := make(chan Votacao, numWorkers)
	var wg syncpkg.WaitGroup
	var mu syncpkg.Mutex
	processed := 0

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for votacao := range votacoesChan {
				s.processarVotacao(ctx, votacao, votacoesCollection, proposicoesCollection)
				mu.Lock()
				processed++
				if processed%50 == 0 {
					log.Printf("   Processando votações: %d/%d", processed, len(allVotacoes))
				}
				mu.Unlock()
			}
		}()
	}

	for _, votacao := range allVotacoes {
		votacoesChan <- votacao
	}
	close(votacoesChan)
	wg.Wait()

	log.Println("✅ Sincronização de votações concluída!")
	return nil
}

// SyncProposicoes sincroniza proposições dos deputados
func (s *CamaraSync) SyncProposicoes(ctx context.Context, ano int) error {
	log.Printf("📥 Buscando proposições do ano %d...", ano)

	// Buscar todos os deputados do banco para buscar suas proposições
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

	proposicoesCollection := s.db.Collection("proposicoes")
	log.Printf("   %d políticos encontrados para sincronizar proposições", len(politicos))

	// Para cada deputado, buscar suas proposições
	// A API da Câmara permite buscar proposições por autor
	// Mas precisamos do ID externo do deputado
	// Por enquanto, vamos buscar todas as proposições do ano e filtrar depois

	// Primeiro, descobrir quantas páginas existem
	url := fmt.Sprintf("%s/proposicoes?ano=%d&itens=100&ordem=ASC&ordenarPor=id", BaseURL, ano)
	var firstResp ProposicoesResponse
	if err := s.client.Get(url, &firstResp); err != nil {
		return fmt.Errorf("erro ao buscar primeira página: %w", err)
	}

	// Descobrir total de páginas pelos links
	totalPaginas := 1
	foundLast := false
	for _, link := range firstResp.Links {
		if link.Rel == "last" {
			// Extrair número da página do link "last"
			// Formato: ...&pagina=447 ou ...?pagina=447
			if strings.Contains(link.Href, "pagina=") {
				parts := strings.Split(link.Href, "pagina=")
				if len(parts) > 1 {
					pageParts := strings.Split(parts[1], "&")
					if len(pageParts) > 0 {
						pageParts = strings.Split(pageParts[0], "?")
						var pageNum int
						if _, err := fmt.Sscanf(pageParts[0], "%d", &pageNum); err == nil && pageNum > 0 {
							totalPaginas = pageNum
							foundLast = true
						}
					}
				}
			}
			break
		}
	}

	// Se não encontrou link "last", usar uma estimativa conservadora
	// e buscar páginas até encontrar uma vazia (em paralelo)
	if !foundLast {
		log.Printf("   Link 'last' não encontrado, usando busca adaptativa...")
		// Começar com uma estimativa e ajustar conforme necessário
		totalPaginas = 500 // Estimativa conservadora
	}

	log.Printf("📄 Total de páginas: %d (buscar todas em paralelo...)", totalPaginas)

	// Adicionar primeira página já buscada
	allProposicoes := firstResp.Dados

	// Buscar todas as outras páginas em paralelo
	paginasChan := make(chan int, totalPaginas)
	var wgPaginas syncpkg.WaitGroup
	var muPaginas syncpkg.Mutex
	errosPaginas := 0

	// Workers para buscar páginas (muitos workers para buscar páginas rapidamente)
	const numPageWorkers = 100 // Aumentado para buscar muitas páginas em paralelo
	for i := 0; i < numPageWorkers; i++ {
		wgPaginas.Add(1)
		go func() {
			defer wgPaginas.Done()
			for pagina := range paginasChan {
				pageURL := fmt.Sprintf("%s/proposicoes?ano=%d&itens=100&ordem=ASC&ordenarPor=id&pagina=%d", BaseURL, ano, pagina)
				var resp ProposicoesResponse
				if err := s.client.Get(pageURL, &resp); err != nil {
					// Página não existe ou erro - ignorar (pode ser além do limite)
					muPaginas.Lock()
					errosPaginas++
					muPaginas.Unlock()
					continue
				}

				// Se página vazia, pode ser que chegamos ao fim
				if len(resp.Dados) == 0 {
					continue
				}

				muPaginas.Lock()
				allProposicoes = append(allProposicoes, resp.Dados...)
				if len(allProposicoes)%5000 == 0 {
					log.Printf("   Carregadas %d proposições de páginas...", len(allProposicoes))
				}
				muPaginas.Unlock()
			}
		}()
	}

	// Enviar todas as páginas para processamento (página 1 já foi buscada)
	for i := 2; i <= totalPaginas; i++ {
		paginasChan <- i
	}
	close(paginasChan)
	wgPaginas.Wait()

	if errosPaginas > 0 {
		log.Printf("⚠️  %d erros ao buscar páginas", errosPaginas)
	}

	log.Printf("📊 Total: %d proposições encontradas", len(allProposicoes))

	// Processar proposições em paralelo (aumentado para acelerar)
	const numWorkers = 30
	proposicoesChan := make(chan Proposicao, numWorkers)
	var wg syncpkg.WaitGroup
	var mu syncpkg.Mutex
	processed := 0

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for prop := range proposicoesChan {
				s.processarProposicao(ctx, prop, proposicoesCollection)
				mu.Lock()
				processed++
				if processed%500 == 0 {
					log.Printf("   Processando proposições: %d/%d (%.1f%%)", processed, len(allProposicoes), float64(processed)/float64(len(allProposicoes))*100)
				}
				mu.Unlock()
			}
		}()
	}

	for _, prop := range allProposicoes {
		proposicoesChan <- prop
	}
	close(proposicoesChan)
	wg.Wait()

	log.Println("✅ Sincronização de proposições concluída!")
	return nil
}

// processarVotacao processa uma votação individual (usado em goroutines)
func (s *CamaraSync) processarVotacao(ctx context.Context, votacao Votacao, votacoesCollection, proposicoesCollection *mongo.Collection) {
	// Buscar proposição se existir
	var proposicaoID primitive.ObjectID
	if votacao.Proposicao != nil {
		proposicao := domain.Proposicao{
			Tipo:   votacao.Proposicao.Siglum,
			Numero: fmt.Sprintf("%d", votacao.Proposicao.Numero),
			Ano:    votacao.Proposicao.Ano,
			Ementa: votacao.Proposicao.Ementa,
		}

		filter := bson.M{
			"tipo":   proposicao.Tipo,
			"numero": proposicao.Numero,
			"ano":    proposicao.Ano,
		}

		update := bson.M{
			"$set": bson.M{
				"ementa":     proposicao.Ementa,
				"updated_at": time.Now(),
			},
			"$setOnInsert": bson.M{
				"_id":        primitive.NewObjectID(),
				"situacao":   domain.SituacaoEmTramitacao,
				"created_at": time.Now(),
			},
		}

		var result struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		err := proposicoesCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
		if err == nil {
			proposicaoID = result.ID
		}
	}

	// Buscar votos dos deputados nesta votação
	votosURL := fmt.Sprintf("%s/votacoes/%s/votos", BaseURL, votacao.ID)
	var votosResp VotoDeputadoResponse
	if err := s.client.Get(votosURL, &votosResp); err != nil {
		log.Printf("⚠️  Erro ao buscar votos da votação %s: %v", votacao.ID, err)
		return
	}

	dataVotacao := ParseDate(votacao.Data)
	if dataVotacao.IsZero() {
		dataVotacao = time.Now()
	}

	// Salvar cada voto
	for _, voto := range votosResp.Dados {
		var politico domain.Politico
		politicoFilter := bson.M{"id_externo_camara": voto.Deputado.ID}
		err := s.db.Collection("politicos").FindOne(ctx, politicoFilter).Decode(&politico)
		if err != nil {
			continue
		}

		tipoVoto := mapTipoVoto(voto.TipoVoto)
		filter := bson.M{
			"politico_id": politico.ID,
			"data":        dataVotacao,
			"sessao":      votacao.SiglaOrgao,
		}

		update := bson.M{
			"$set": bson.M{
				"voto":          tipoVoto,
				"proposicao_id": proposicaoID,
				"data":          dataVotacao,
				"sessao":        votacao.SiglaOrgao,
			},
			"$setOnInsert": bson.M{
				"_id": primitive.NewObjectID(),
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err = votacoesCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("⚠️  Erro ao salvar voto: %v", err)
		}
	}
}

// processarProposicao processa uma proposição individual (usado em goroutines)
func (s *CamaraSync) processarProposicao(ctx context.Context, prop Proposicao, proposicoesCollection *mongo.Collection) {
	// Buscar detalhes da proposição
	url := fmt.Sprintf("%s/proposicoes/%d", BaseURL, prop.ID)
	var detalhes ProposicaoDetalheResponse
	if err := s.client.Get(url, &detalhes); err != nil {
		// Erro não crítico, continuar
		return
	}

	d := detalhes.Dados

	// Buscar autores da proposição (essencial)
	autoresURL := fmt.Sprintf("%s/proposicoes/%d/autores", BaseURL, prop.ID)
	var autoresResp AutoresResponse
	var autorID primitive.ObjectID
	var coautoresIDs []primitive.ObjectID

	if err := s.client.Get(autoresURL, &autoresResp); err == nil {
		for _, autor := range autoresResp.Dados {
			var politico domain.Politico
			politicoFilter := bson.M{"id_externo_camara": autor.ID}
			err := s.db.Collection("politicos").FindOne(ctx, politicoFilter).Decode(&politico)
			if err == nil {
				if autor.CodTipo == 1 || autorID.IsZero() {
					autorID = politico.ID
				} else {
					coautoresIDs = append(coautoresIDs, politico.ID)
				}
			}
		}
	}

	// Buscar tramitações e temas em paralelo (opcional, não bloqueia)
	var tramitacoes []domain.TramitacaoItem
	var temas []string
	var wg syncpkg.WaitGroup
	var mu syncpkg.Mutex

	// Tramitações (pode ser carregado depois se necessário)
	wg.Add(1)
	go func() {
		defer wg.Done()
		tramitacoesURL := fmt.Sprintf("%s/proposicoes/%d/tramitacoes?itens=100", BaseURL, prop.ID)
		var tramitacoesResp TramitacoesResponse
		if err := s.client.Get(tramitacoesURL, &tramitacoesResp); err == nil {
			var tempTramitacoes []domain.TramitacaoItem
			for _, tram := range tramitacoesResp.Dados {
				dataTram := ParseDate(tram.DataHora)
				if !dataTram.IsZero() {
					tempTramitacoes = append(tempTramitacoes, domain.TramitacaoItem{
						Data:      dataTram,
						Descricao: tram.DescTramitacao,
						Orgao:     tram.SiglaOrgao,
					})
				}
			}
			mu.Lock()
			tramitacoes = tempTramitacoes
			mu.Unlock()
		}
	}()

	// Temas (pode ser carregado depois se necessário)
	wg.Add(1)
	go func() {
		defer wg.Done()
		temasURL := fmt.Sprintf("%s/proposicoes/%d/temas", BaseURL, prop.ID)
		var temasResp TemasResponse
		if err := s.client.Get(temasURL, &temasResp); err == nil {
			var tempTemas []string
			for _, tema := range temasResp.Dados {
				if tema.Nome != "" {
					tempTemas = append(tempTemas, tema.Nome)
				}
			}
			mu.Lock()
			temas = tempTemas
			mu.Unlock()
		}
	}()

	// Aguardar tramitações e temas (timeout de 3 segundos para não travar)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Dados carregados
	case <-time.After(3 * time.Second):
		// Timeout - continuar sem tramitações/temas (podem ser carregados depois)
	}

	// Mapear situação
	situacao := domain.SituacaoEmTramitacao
	if d.StatusProposicao != nil {
		descSituacao := strings.ToUpper(d.StatusProposicao.DescSituacao)
		if strings.Contains(descSituacao, "APROVADA") || strings.Contains(descSituacao, "SANCIONADA") {
			situacao = domain.SituacaoAprovada
		} else if strings.Contains(descSituacao, "REJEITADA") {
			situacao = domain.SituacaoRejeitada
		} else if strings.Contains(descSituacao, "ARQUIVADA") {
			situacao = domain.SituacaoArquivada
		} else if strings.Contains(descSituacao, "RETIRADA") {
			situacao = domain.SituacaoRetirada
		}
	}

	proposicaoDoc := domain.Proposicao{
		Tipo:         d.SiglaTipo,
		Numero:       fmt.Sprintf("%d", d.Numero),
		Ano:          d.Ano,
		Ementa:       d.Ementa,
		AutorID:      autorID,
		CoautoresIDs: coautoresIDs,
		Situacao:     situacao,
		Tema:         temas,
		Tramitacao:   tramitacoes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	filter := bson.M{
		"tipo":   proposicaoDoc.Tipo,
		"numero": proposicaoDoc.Numero,
		"ano":    proposicaoDoc.Ano,
	}

	update := bson.M{
		"$set": bson.M{
			"ementa":        proposicaoDoc.Ementa,
			"autor_id":      proposicaoDoc.AutorID,
			"coautores_ids": proposicaoDoc.CoautoresIDs,
			"situacao":      proposicaoDoc.Situacao,
			"tema":          proposicaoDoc.Tema,
			"tramitacao":    proposicaoDoc.Tramitacao,
			"updated_at":    time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":        primitive.NewObjectID(),
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := proposicoesCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("⚠️  Erro ao salvar proposição: %v", err)
	}
}

// SyncPresencas sincroniza presenças dos deputados em eventos/sessões
func (s *CamaraSync) SyncPresencas(ctx context.Context, ano int) error {
	log.Printf("📥 Buscando presenças em eventos do ano %d...", ano)

	// Buscar eventos do ano (sessões plenárias, reuniões de comissões, etc.)
	url := fmt.Sprintf("%s/eventos?ano=%d&itens=100&ordem=ASC&ordenarPor=dataHoraInicio", BaseURL, ano)
	presencasCollection := s.db.Collection("presencas")

	var allEventos []Evento
	for url != "" {
		var resp EventosResponse
		if err := s.client.Get(url, &resp); err != nil {
			return fmt.Errorf("erro ao buscar eventos: %w", err)
		}

		allEventos = append(allEventos, resp.Dados...)
		log.Printf("   Carregados %d eventos...", len(allEventos))

		// Próxima página
		url = ""
		for _, link := range resp.Links {
			if link.Rel == "next" {
				url = link.Href
				break
			}
		}
	}

	log.Printf("📊 Total: %d eventos encontrados", len(allEventos))

	// Processar eventos em paralelo
	const numWorkers = 5
	eventosChan := make(chan Evento, numWorkers)
	var wg syncpkg.WaitGroup
	var mu syncpkg.Mutex
	processed := 0

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for evento := range eventosChan {
				s.processarEvento(ctx, evento, presencasCollection)
				mu.Lock()
				processed++
				if processed%50 == 0 {
					log.Printf("   Processando eventos: %d/%d", processed, len(allEventos))
				}
				mu.Unlock()
			}
		}()
	}

	for _, evento := range allEventos {
		eventosChan <- evento
	}
	close(eventosChan)
	wg.Wait()

	log.Println("✅ Sincronização de presenças concluída!")
	return nil
}

// processarEvento processa um evento individual (usado em goroutines)
func (s *CamaraSync) processarEvento(ctx context.Context, evento Evento, presencasCollection *mongo.Collection) {
	// Buscar presenças do evento
	presencasURL := fmt.Sprintf("%s/eventos/%d/presencas", BaseURL, evento.ID)
	var presencasResp PresencasEventoResponse
	if err := s.client.Get(presencasURL, &presencasResp); err != nil {
		// Alguns eventos podem não ter presenças registradas
		return
	}

	dataEvento := ParseDate(evento.DataHoraInicio)
	if dataEvento.IsZero() {
		return
	}

	tipoSessao := evento.DescricaoTipo
	if tipoSessao == "" {
		tipoSessao = "Evento"
	}

	// Salvar cada presença
	for _, pres := range presencasResp.Dados {
		var politico domain.Politico
		politicoFilter := bson.M{"id_externo_camara": pres.Deputado.ID}
		err := s.db.Collection("politicos").FindOne(ctx, politicoFilter).Decode(&politico)
		if err != nil {
			continue
		}

		filter := bson.M{
			"politico_id": politico.ID,
			"data":        dataEvento,
			"tipo_sessao": tipoSessao,
		}

		update := bson.M{
			"$set": bson.M{
				"presente": true,
			},
			"$setOnInsert": bson.M{
				"_id": primitive.NewObjectID(),
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err = presencasCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("⚠️  Erro ao salvar presença: %v", err)
		}
	}
}

// mapTipoVoto converte o tipo de voto da API para nosso modelo
func mapTipoVoto(tipo string) domain.TipoVoto {
	switch strings.ToUpper(tipo) {
	case "SIM":
		return domain.VotoSim
	case "NÃO", "NAO":
		return domain.VotoNao
	case "ABSTENÇÃO", "ABSTENCAO":
		return domain.VotoAbstencao
	case "AUSENTE":
		return domain.VotoAusente
	case "OBSTRUÇÃO", "OBSTRUCAO":
		return domain.VotoObstrucao
	default:
		return domain.VotoAusente
	}
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
