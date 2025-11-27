package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lupa-cidada/backend/internal/services"
)

type EstatisticasHandler struct {
	service *services.PoliticoService
}

func NewEstatisticasHandler(service *services.PoliticoService) *EstatisticasHandler {
	return &EstatisticasHandler{service: service}
}

func (h *EstatisticasHandler) Geral(c echo.Context) error {
	ctx := c.Request().Context()

	totalPoliticos, err := h.service.ContarPoliticos(ctx)
	if err != nil {
		totalPoliticos = 0
	}

	totalVotacoes, err := h.service.ContarVotacoes(ctx)
	if err != nil {
		totalVotacoes = 0
	}

	totalProposicoes, err := h.service.ContarProposicoes(ctx)
	if err != nil {
		totalProposicoes = 0
	}

	totalDespesas, err := h.service.TotalDespesas(ctx)
	if err != nil {
		totalDespesas = 0
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalPoliticos":   totalPoliticos,
		"totalVotacoes":    totalVotacoes,
		"totalProposicoes": totalProposicoes,
		"totalDespesas":    totalDespesas,
	})
}

func (h *EstatisticasHandler) Ranking(c echo.Context) error {
	// TODO: Implementar ranking quando houver dados suficientes
	return c.JSON(http.StatusOK, []interface{}{})
}

