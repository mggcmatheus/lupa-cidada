package services

import (
	"context"
	"strings"

	"github.com/lupa-cidada/backend/internal/domain"
	"github.com/lupa-cidada/backend/internal/mock"
	"github.com/lupa-cidada/backend/internal/repository"
)

type PoliticoService struct {
	debug          bool
	politicoRepo   *repository.PoliticoRepository
	votacaoRepo    *repository.VotacaoRepository
	despesaRepo    *repository.DespesaRepository
	proposicaoRepo *repository.ProposicaoRepository
}

func NewPoliticoService(
	debug bool,
	politicoRepo *repository.PoliticoRepository,
	votacaoRepo *repository.VotacaoRepository,
	despesaRepo *repository.DespesaRepository,
	proposicaoRepo *repository.ProposicaoRepository,
) *PoliticoService {
	return &PoliticoService{
		debug:          debug,
		politicoRepo:   politicoRepo,
		votacaoRepo:    votacaoRepo,
		despesaRepo:    despesaRepo,
		proposicaoRepo: proposicaoRepo,
	}
}

func (s *PoliticoService) Listar(ctx context.Context, filtros domain.FiltrosPoliticos) (*domain.PaginatedResponse[domain.Politico], error) {
	if s.debug {
		return s.listarMock(filtros), nil
	}
	return s.politicoRepo.Listar(ctx, filtros)
}

func (s *PoliticoService) listarMock(filtros domain.FiltrosPoliticos) *domain.PaginatedResponse[domain.Politico] {
	allPoliticos := mock.Politicos()
	var resultado []domain.Politico

	for _, p := range allPoliticos {
		// Filtro por nome
		if filtros.Nome != "" {
			nome := strings.ToLower(p.Nome)
			nomeCivil := strings.ToLower(p.NomeCivil)
			termo := strings.ToLower(filtros.Nome)
			if !strings.Contains(nome, termo) && !strings.Contains(nomeCivil, termo) {
				continue
			}
		}

		// Filtro por partido
		if len(filtros.Partido) > 0 {
			found := false
			for _, partido := range filtros.Partido {
				if p.Partido.Sigla == partido {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por cargo
		if len(filtros.Cargo) > 0 {
			found := false
			for _, cargo := range filtros.Cargo {
				if p.CargoAtual.Tipo == cargo {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por esfera
		if len(filtros.Esfera) > 0 {
			found := false
			for _, esfera := range filtros.Esfera {
				if p.CargoAtual.Esfera == esfera {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por estado
		if len(filtros.Estado) > 0 {
			found := false
			for _, estado := range filtros.Estado {
				if p.CargoAtual.Estado == estado {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por status (em exercício)
		if filtros.EmExercicio != nil {
			if p.CargoAtual.EmExercicio != *filtros.EmExercicio {
				continue
			}
		}

		// Filtro por gênero
		if len(filtros.Genero) > 0 {
			found := false
			for _, genero := range filtros.Genero {
				if p.Genero == genero {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		resultado = append(resultado, p)
	}

	// Paginação
	pagina := filtros.Pagina
	if pagina < 1 {
		pagina = 1
	}
	porPagina := filtros.PorPagina
	if porPagina < 1 || porPagina > 100 {
		porPagina = 12
	}

	total := int64(len(resultado))
	start := (pagina - 1) * porPagina
	end := start + porPagina

	if start > len(resultado) {
		start = len(resultado)
	}
	if end > len(resultado) {
		end = len(resultado)
	}

	paginado := resultado[start:end]

	totalPaginas := int(total) / porPagina
	if int(total)%porPagina > 0 {
		totalPaginas++
	}

	return &domain.PaginatedResponse[domain.Politico]{
		Data:         paginado,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}
}

func (s *PoliticoService) BuscarPorID(ctx context.Context, id string) (*domain.Politico, error) {
	if s.debug {
		p := mock.GetPoliticoByID(id)
		if p == nil {
			return nil, nil
		}
		return p, nil
	}
	return s.politicoRepo.BuscarPorID(ctx, id)
}

func (s *PoliticoService) BuscarPorIDs(ctx context.Context, ids []string) ([]domain.Politico, error) {
	if s.debug {
		var result []domain.Politico
		for _, id := range ids {
			if p := mock.GetPoliticoByID(id); p != nil {
				result = append(result, *p)
			}
		}
		return result, nil
	}
	return s.politicoRepo.BuscarPorIDs(ctx, ids)
}

func (s *PoliticoService) Buscar(ctx context.Context, query string, limite int) ([]domain.Politico, error) {
	if limite <= 0 || limite > 50 {
		limite = 10
	}

	if s.debug {
		var result []domain.Politico
		termo := strings.ToLower(query)
		for _, p := range mock.Politicos() {
			if strings.Contains(strings.ToLower(p.Nome), termo) ||
				strings.Contains(strings.ToLower(p.NomeCivil), termo) {
				result = append(result, p)
				if len(result) >= limite {
					break
				}
			}
		}
		return result, nil
	}

	return s.politicoRepo.Buscar(ctx, query, limite)
}

func (s *PoliticoService) BuscarEstatisticas(ctx context.Context, id string) (*domain.EstatisticasPolitico, error) {
	if s.debug {
		stats := mock.GetEstatisticasByID(id)
		return stats, nil
	}

	// Contar votos
	votosCounts, err := s.votacaoRepo.ContarPorPolitico(ctx, id)
	if err != nil {
		return nil, err
	}

	totalVotacoes := 0
	for _, count := range votosCounts {
		totalVotacoes += count
	}

	// Contar proposições
	totalProposicoes, proposicoesAprovadas, err := s.proposicaoRepo.ContarPorAutor(ctx, id)
	if err != nil {
		return nil, err
	}

	// Calcular despesas
	totalDespesas, err := s.despesaRepo.TotalPorPolitico(ctx, id)
	if err != nil {
		return nil, err
	}

	mediaGastoMensal, err := s.despesaRepo.MediaMensalPorPolitico(ctx, id)
	if err != nil {
		return nil, err
	}

	// Calcular presença (baseado nas votações)
	var percentualPresenca float64
	if totalVotacoes > 0 {
		ausencias := votosCounts[domain.VotoAusente]
		percentualPresenca = float64(totalVotacoes-ausencias) / float64(totalVotacoes) * 100
	}

	return &domain.EstatisticasPolitico{
		TotalVotacoes:        totalVotacoes,
		VotosSim:             votosCounts[domain.VotoSim],
		VotosNao:             votosCounts[domain.VotoNao],
		Abstencoes:           votosCounts[domain.VotoAbstencao],
		Ausencias:            votosCounts[domain.VotoAusente],
		PercentualPresenca:   percentualPresenca,
		TotalProposicoes:     int(totalProposicoes),
		ProposicoesAprovadas: int(proposicoesAprovadas),
		TotalDespesas:        totalDespesas,
		MediaGastoMensal:     mediaGastoMensal,
	}, nil
}

func (s *PoliticoService) ListarVotacoes(ctx context.Context, politicoID string, pagina, porPagina int) (*domain.PaginatedResponse[domain.Votacao], error) {
	if s.debug {
		// Retorna lista vazia em modo debug (votações são complexas de mockar)
		return &domain.PaginatedResponse[domain.Votacao]{
			Data:         []domain.Votacao{},
			Total:        0,
			Pagina:       1,
			PorPagina:    porPagina,
			TotalPaginas: 0,
		}, nil
	}
	return s.votacaoRepo.ListarPorPolitico(ctx, politicoID, pagina, porPagina)
}

func (s *PoliticoService) ListarDespesas(ctx context.Context, politicoID string, ano, mes *int, pagina, porPagina int) (*domain.PaginatedResponse[domain.Despesa], error) {
	if s.debug {
		return &domain.PaginatedResponse[domain.Despesa]{
			Data:         []domain.Despesa{},
			Total:        0,
			Pagina:       1,
			PorPagina:    porPagina,
			TotalPaginas: 0,
		}, nil
	}
	return s.despesaRepo.ListarPorPolitico(ctx, politicoID, ano, mes, pagina, porPagina)
}

func (s *PoliticoService) ListarProposicoes(ctx context.Context, politicoID string, pagina, porPagina int) (*domain.PaginatedResponse[domain.Proposicao], error) {
	if s.debug {
		return &domain.PaginatedResponse[domain.Proposicao]{
			Data:         []domain.Proposicao{},
			Total:        0,
			Pagina:       1,
			PorPagina:    porPagina,
			TotalPaginas: 0,
		}, nil
	}
	return s.proposicaoRepo.ListarPorAutor(ctx, politicoID, pagina, porPagina)
}

func (s *PoliticoService) Comparar(ctx context.Context, ids []string) (map[string]interface{}, error) {
	politicos, err := s.BuscarPorIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	estatisticas := make(map[string]*domain.EstatisticasPolitico)
	for _, p := range politicos {
		stats, err := s.BuscarEstatisticas(ctx, p.ID.Hex())
		if err != nil {
			continue
		}
		estatisticas[p.ID.Hex()] = stats
	}

	return map[string]interface{}{
		"politicos":    politicos,
		"estatisticas": estatisticas,
	}, nil
}

func (s *PoliticoService) ContarPoliticos(ctx context.Context) (int64, error) {
	if s.debug {
		return int64(len(mock.Politicos())), nil
	}
	return s.politicoRepo.Contar(ctx)
}

func (s *PoliticoService) ContarVotacoes(ctx context.Context) (int64, error) {
	if s.debug {
		// Soma todas as votações mockadas
		var total int64
		for _, stats := range mock.Estatisticas() {
			total += int64(stats.TotalVotacoes)
		}
		return total, nil
	}
	return s.votacaoRepo.Contar(ctx)
}

func (s *PoliticoService) ContarProposicoes(ctx context.Context) (int64, error) {
	if s.debug {
		var total int64
		for _, stats := range mock.Estatisticas() {
			total += int64(stats.TotalProposicoes)
		}
		return total, nil
	}
	return s.proposicaoRepo.Contar(ctx)
}

func (s *PoliticoService) TotalDespesas(ctx context.Context) (float64, error) {
	if s.debug {
		var total float64
		for _, stats := range mock.Estatisticas() {
			total += stats.TotalDespesas
		}
		return total, nil
	}
	return s.despesaRepo.TotalGeral(ctx)
}
