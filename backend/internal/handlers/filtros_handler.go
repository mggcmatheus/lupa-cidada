package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FiltrosHandler struct {
	db    *mongo.Database
	debug bool
}

func NewFiltrosHandler(db *mongo.Database, debug bool) *FiltrosHandler {
	return &FiltrosHandler{db: db, debug: debug}
}

func (h *FiltrosHandler) ListarPartidos(c echo.Context) error {
	if h.debug || h.db == nil {
		// Retorna partidos mockados
		partidos := []map[string]string{
			{"sigla": "PT", "nome": "Partido dos Trabalhadores", "cor": "#CC0000"},
			{"sigla": "PL", "nome": "Partido Liberal", "cor": "#003366"},
			{"sigla": "UNIÃO", "nome": "União Brasil", "cor": "#2E3092"},
			{"sigla": "PP", "nome": "Progressistas", "cor": "#0066CC"},
			{"sigla": "MDB", "nome": "Movimento Democrático Brasileiro", "cor": "#00AA00"},
			{"sigla": "PSD", "nome": "Partido Social Democrático", "cor": "#FF6600"},
			{"sigla": "REPUBLICANOS", "nome": "Republicanos", "cor": "#0033CC"},
			{"sigla": "PDT", "nome": "Partido Democrático Trabalhista", "cor": "#FF0000"},
			{"sigla": "PSDB", "nome": "Partido da Social Democracia Brasileira", "cor": "#003399"},
			{"sigla": "PSOL", "nome": "Partido Socialismo e Liberdade", "cor": "#FFD700"},
			{"sigla": "PSB", "nome": "Partido Socialista Brasileiro", "cor": "#FF6347"},
			{"sigla": "PODE", "nome": "Podemos", "cor": "#00CED1"},
			{"sigla": "CIDADANIA", "nome": "Cidadania", "cor": "#9932CC"},
			{"sigla": "AVANTE", "nome": "Avante", "cor": "#FF8C00"},
			{"sigla": "SOLIDARIEDADE", "nome": "Solidariedade", "cor": "#FF4500"},
			{"sigla": "PCdoB", "nome": "Partido Comunista do Brasil", "cor": "#8B0000"},
			{"sigla": "PV", "nome": "Partido Verde", "cor": "#228B22"},
			{"sigla": "NOVO", "nome": "Partido Novo", "cor": "#FF6600"},
			{"sigla": "REDE", "nome": "Rede Sustentabilidade", "cor": "#00AA66"},
		}
		return c.JSON(http.StatusOK, partidos)
	}

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
