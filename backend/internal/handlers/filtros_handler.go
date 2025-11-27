package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FiltrosHandler struct {
	db *mongo.Database
}

func NewFiltrosHandler(db *mongo.Database) *FiltrosHandler {
	return &FiltrosHandler{db: db}
}

func (h *FiltrosHandler) ListarPartidos(c echo.Context) error {
	cursor, err := h.db.Collection("partidos").Find(context.Background(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao listar partidos",
		})
	}
	defer cursor.Close(context.Background())

	var partidos []bson.M
	if err := cursor.All(context.Background(), &partidos); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao processar partidos",
		})
	}

	return c.JSON(http.StatusOK, partidos)
}

func (h *FiltrosHandler) ListarEstados(c echo.Context) error {
	estados := []map[string]string{
		{"sigla": "AC", "nome": "Acre"},
		{"sigla": "AL", "nome": "Alagoas"},
		{"sigla": "AP", "nome": "Amapá"},
		{"sigla": "AM", "nome": "Amazonas"},
		{"sigla": "BA", "nome": "Bahia"},
		{"sigla": "CE", "nome": "Ceará"},
		{"sigla": "DF", "nome": "Distrito Federal"},
		{"sigla": "ES", "nome": "Espírito Santo"},
		{"sigla": "GO", "nome": "Goiás"},
		{"sigla": "MA", "nome": "Maranhão"},
		{"sigla": "MT", "nome": "Mato Grosso"},
		{"sigla": "MS", "nome": "Mato Grosso do Sul"},
		{"sigla": "MG", "nome": "Minas Gerais"},
		{"sigla": "PA", "nome": "Pará"},
		{"sigla": "PB", "nome": "Paraíba"},
		{"sigla": "PR", "nome": "Paraná"},
		{"sigla": "PE", "nome": "Pernambuco"},
		{"sigla": "PI", "nome": "Piauí"},
		{"sigla": "RJ", "nome": "Rio de Janeiro"},
		{"sigla": "RN", "nome": "Rio Grande do Norte"},
		{"sigla": "RS", "nome": "Rio Grande do Sul"},
		{"sigla": "RO", "nome": "Rondônia"},
		{"sigla": "RR", "nome": "Roraima"},
		{"sigla": "SC", "nome": "Santa Catarina"},
		{"sigla": "SP", "nome": "São Paulo"},
		{"sigla": "SE", "nome": "Sergipe"},
		{"sigla": "TO", "nome": "Tocantins"},
	}

	return c.JSON(http.StatusOK, estados)
}

func (h *FiltrosHandler) ListarCargos(c echo.Context) error {
	cargos := []map[string]string{
		{"valor": "DEPUTADO_FEDERAL", "label": "Deputado Federal"},
		{"valor": "SENADOR", "label": "Senador"},
		{"valor": "DEPUTADO_ESTADUAL", "label": "Deputado Estadual"},
		{"valor": "DEPUTADO_DISTRITAL", "label": "Deputado Distrital"},
		{"valor": "VEREADOR", "label": "Vereador"},
		{"valor": "PREFEITO", "label": "Prefeito"},
		{"valor": "GOVERNADOR", "label": "Governador"},
		{"valor": "PRESIDENTE", "label": "Presidente"},
	}

	return c.JSON(http.StatusOK, cargos)
}

