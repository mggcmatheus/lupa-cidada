package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lupa-cidada/backend/internal/domain"
	"github.com/lupa-cidada/backend/internal/services"
)

type PoliticoHandler struct {
	service *services.PoliticoService
}

func NewPoliticoHandler(service *services.PoliticoService) *PoliticoHandler {
	return &PoliticoHandler{service: service}
}

func (h *PoliticoHandler) Listar(c echo.Context) error {
	var filtros domain.FiltrosPoliticos

	// Parse manual dos filtros
	filtros.Nome = c.QueryParam("nome")
	filtros.Pagina, _ = strconv.Atoi(c.QueryParam("pagina"))
	filtros.PorPagina, _ = strconv.Atoi(c.QueryParam("porPagina"))
	filtros.OrdenarPor = c.QueryParam("ordenarPor")
	filtros.Ordem = c.QueryParam("ordem")

	if partido := c.QueryParam("partido"); partido != "" {
		filtros.Partido = strings.Split(partido, ",")
	}

	if cargo := c.QueryParam("cargo"); cargo != "" {
		for _, c := range strings.Split(cargo, ",") {
			filtros.Cargo = append(filtros.Cargo, domain.Cargo(c))
		}
	}

	if esfera := c.QueryParam("esfera"); esfera != "" {
		for _, e := range strings.Split(esfera, ",") {
			filtros.Esfera = append(filtros.Esfera, domain.Esfera(e))
		}
	}

	if estado := c.QueryParam("estado"); estado != "" {
		filtros.Estado = strings.Split(estado, ",")
	}

	if emExercicio := c.QueryParam("emExercicio"); emExercicio != "" {
		val := emExercicio == "true"
		filtros.EmExercicio = &val
	}

	if genero := c.QueryParam("genero"); genero != "" {
		for _, g := range strings.Split(genero, ",") {
			filtros.Genero = append(filtros.Genero, domain.Genero(g))
		}
	}

	result, err := h.service.Listar(c.Request().Context(), filtros)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao listar políticos",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PoliticoHandler) BuscarPorID(c echo.Context) error {
	id := c.Param("id")

	politico, err := h.service.BuscarPorID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Político não encontrado",
		})
	}

	return c.JSON(http.StatusOK, politico)
}

func (h *PoliticoHandler) BuscarEstatisticas(c echo.Context) error {
	id := c.Param("id")

	estatisticas, err := h.service.BuscarEstatisticas(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao buscar estatísticas",
		})
	}

	return c.JSON(http.StatusOK, estatisticas)
}

func (h *PoliticoHandler) ListarVotacoes(c echo.Context) error {
	id := c.Param("id")
	pagina, _ := strconv.Atoi(c.QueryParam("pagina"))
	porPagina, _ := strconv.Atoi(c.QueryParam("porPagina"))

	result, err := h.service.ListarVotacoes(c.Request().Context(), id, pagina, porPagina)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao listar votações",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PoliticoHandler) ListarDespesas(c echo.Context) error {
	id := c.Param("id")
	pagina, _ := strconv.Atoi(c.QueryParam("pagina"))
	porPagina, _ := strconv.Atoi(c.QueryParam("porPagina"))

	var ano, mes *int
	if anoStr := c.QueryParam("ano"); anoStr != "" {
		a, _ := strconv.Atoi(anoStr)
		ano = &a
	}
	if mesStr := c.QueryParam("mes"); mesStr != "" {
		m, _ := strconv.Atoi(mesStr)
		mes = &m
	}

	result, err := h.service.ListarDespesas(c.Request().Context(), id, ano, mes, pagina, porPagina)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao listar despesas",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PoliticoHandler) ListarProposicoes(c echo.Context) error {
	id := c.Param("id")
	pagina, _ := strconv.Atoi(c.QueryParam("pagina"))
	porPagina, _ := strconv.Atoi(c.QueryParam("porPagina"))

	result, err := h.service.ListarProposicoes(c.Request().Context(), id, pagina, porPagina)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao listar proposições",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PoliticoHandler) ListarPresencas(c echo.Context) error {
	// TODO: Implementar quando o repositório de presenças for usado
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":         []interface{}{},
		"total":        0,
		"pagina":       1,
		"porPagina":    50,
		"totalPaginas": 0,
	})
}

func (h *PoliticoHandler) Comparar(c echo.Context) error {
	idsParam := c.QueryParam("ids")
	if idsParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IDs são obrigatórios",
		})
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Mínimo de 2 políticos para comparação",
		})
	}

	if len(ids) > 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Máximo de 4 políticos para comparação",
		})
	}

	result, err := h.service.Comparar(c.Request().Context(), ids)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao comparar políticos",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PoliticoHandler) Buscar(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Query de busca é obrigatória",
		})
	}

	limite, _ := strconv.Atoi(c.QueryParam("limite"))

	result, err := h.service.Buscar(c.Request().Context(), query, limite)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao buscar políticos",
		})
	}

	return c.JSON(http.StatusOK, result)
}

