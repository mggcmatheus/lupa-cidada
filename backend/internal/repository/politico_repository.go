package repository

import (
	"context"
	"time"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PoliticoRepository struct {
	collection *mongo.Collection
}

func NewPoliticoRepository(db *mongo.Database) *PoliticoRepository {
	return &PoliticoRepository{
		collection: db.Collection("politicos"),
	}
}

func (r *PoliticoRepository) Listar(ctx context.Context, filtros domain.FiltrosPoliticos) (*domain.PaginatedResponse[domain.Politico], error) {
	filter := bson.M{}

	// Aplicar filtros
	if filtros.Nome != "" {
		filter["$text"] = bson.M{"$search": filtros.Nome}
	}

	if len(filtros.Partido) > 0 {
		filter["partido.sigla"] = bson.M{"$in": filtros.Partido}
	}

	if len(filtros.Cargo) > 0 {
		filter["cargo_atual.tipo"] = bson.M{"$in": filtros.Cargo}
	}

	if len(filtros.Esfera) > 0 {
		filter["cargo_atual.esfera"] = bson.M{"$in": filtros.Esfera}
	}

	if len(filtros.Estado) > 0 {
		filter["cargo_atual.estado"] = bson.M{"$in": filtros.Estado}
	}

	if len(filtros.Municipio) > 0 {
		filter["cargo_atual.municipio"] = bson.M{"$in": filtros.Municipio}
	}

	if filtros.EmExercicio != nil {
		filter["cargo_atual.em_exercicio"] = *filtros.EmExercicio
	}

	if len(filtros.Genero) > 0 {
		filter["genero"] = bson.M{"$in": filtros.Genero}
	}

	// Configurar paginação
	pagina := filtros.Pagina
	if pagina < 1 {
		pagina = 1
	}

	porPagina := filtros.PorPagina
	if porPagina < 1 || porPagina > 100 {
		porPagina = 12
	}

	skip := int64((pagina - 1) * porPagina)
	limit := int64(porPagina)

	// Configurar ordenação
	sort := bson.D{{Key: "nome", Value: 1}}
	if filtros.OrdenarPor != "" {
		ordem := 1
		if filtros.Ordem == "desc" {
			ordem = -1
		}

		sortField := "nome"
		switch filtros.OrdenarPor {
		case "partido":
			sortField = "partido.sigla"
		case "presenca":
			sortField = "presenca_percentual"
		case "proposicoes":
			sortField = "total_proposicoes"
		case "gastos":
			sortField = "gasto_mensal"
		}
		sort = bson.D{{Key: sortField, Value: ordem}}
	}

	// Contar total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Buscar documentos
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(sort)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var politicos []domain.Politico
	if err := cursor.All(ctx, &politicos); err != nil {
		return nil, err
	}

	totalPaginas := int(total) / porPagina
	if int(total)%porPagina > 0 {
		totalPaginas++
	}

	return &domain.PaginatedResponse[domain.Politico]{
		Data:         politicos,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}, nil
}

func (r *PoliticoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Politico, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var politico domain.Politico
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&politico)
	if err != nil {
		return nil, err
	}

	return &politico, nil
}

func (r *PoliticoRepository) BuscarPorIDs(ctx context.Context, ids []string) ([]domain.Politico, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var politicos []domain.Politico
	if err := cursor.All(ctx, &politicos); err != nil {
		return nil, err
	}

	return politicos, nil
}

func (r *PoliticoRepository) Buscar(ctx context.Context, query string, limite int) ([]domain.Politico, error) {
	filter := bson.M{
		"$text": bson.M{"$search": query},
	}

	opts := options.Find().
		SetLimit(int64(limite)).
		SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}}).
		SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var politicos []domain.Politico
	if err := cursor.All(ctx, &politicos); err != nil {
		return nil, err
	}

	return politicos, nil
}

func (r *PoliticoRepository) Criar(ctx context.Context, politico *domain.Politico) error {
	politico.CreatedAt = time.Now()
	politico.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, politico)
	if err != nil {
		return err
	}

	politico.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PoliticoRepository) Atualizar(ctx context.Context, politico *domain.Politico) error {
	politico.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": politico.ID},
		politico,
	)
	return err
}

func (r *PoliticoRepository) Contar(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

