package services

import (
	"context"

	"github.com/lupa-cidada/backend/internal/domain"
	"github.com/lupa-cidada/backend/internal/repository"
)

type PoliticoService struct {
	politicoRepo   *repository.PoliticoRepository
	votacaoRepo    *repository.VotacaoRepository
	despesaRepo    *repository.DespesaRepository
	proposicaoRepo *repository.ProposicaoRepository
}

func NewPoliticoService(
	politicoRepo *repository.PoliticoRepository,
	votacaoRepo *repository.VotacaoRepository,
	despesaRepo *repository.DespesaRepository,
	proposicaoRepo *repository.ProposicaoRepository,
) *PoliticoService {
	return &PoliticoService{
		politicoRepo:   politicoRepo,
		votacaoRepo:    votacaoRepo,
		despesaRepo:    despesaRepo,
		proposicaoRepo: proposicaoRepo,
	}
}

func (s *PoliticoService) Listar(ctx context.Context, filtros domain.FiltrosPoliticos) (*domain.PaginatedResponse[domain.Politico], error) {
	return s.politicoRepo.Listar(ctx, filtros)
}

func (s *PoliticoService) BuscarPorID(ctx context.Context, id string) (*domain.Politico, error) {
	return s.politicoRepo.BuscarPorID(ctx, id)
}

func (s *PoliticoService) BuscarPorIDs(ctx context.Context, ids []string) ([]domain.Politico, error) {
	return s.politicoRepo.BuscarPorIDs(ctx, ids)
}

func (s *PoliticoService) Buscar(ctx context.Context, query string, limite int) ([]domain.Politico, error) {
	if limite <= 0 || limite > 50 {
		limite = 10
	}
	return s.politicoRepo.Buscar(ctx, query, limite)
}

func (s *PoliticoService) BuscarEstatisticas(ctx context.Context, id string) (*domain.EstatisticasPolitico, error) {
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
	return s.votacaoRepo.ListarPorPolitico(ctx, politicoID, pagina, porPagina)
}

func (s *PoliticoService) ListarDespesas(ctx context.Context, politicoID string, ano, mes *int, pagina, porPagina int) (*domain.PaginatedResponse[domain.Despesa], error) {
	return s.despesaRepo.ListarPorPolitico(ctx, politicoID, ano, mes, pagina, porPagina)
}

func (s *PoliticoService) ListarProposicoes(ctx context.Context, politicoID string, pagina, porPagina int) (*domain.PaginatedResponse[domain.Proposicao], error) {
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
	return s.politicoRepo.Contar(ctx)
}

func (s *PoliticoService) ContarVotacoes(ctx context.Context) (int64, error) {
	return s.votacaoRepo.Contar(ctx)
}

func (s *PoliticoService) ContarProposicoes(ctx context.Context) (int64, error) {
	return s.proposicaoRepo.Contar(ctx)
}

func (s *PoliticoService) TotalDespesas(ctx context.Context) (float64, error) {
	return s.despesaRepo.TotalGeral(ctx)
}

