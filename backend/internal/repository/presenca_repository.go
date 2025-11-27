package repository

import (
	"context"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PresencaRepository struct {
	collection *mongo.Collection
}

func NewPresencaRepository(db *mongo.Database) *PresencaRepository {
	return &PresencaRepository{
		collection: db.Collection("presencas"),
	}
}

func (r *PresencaRepository) ListarPorPolitico(ctx context.Context, politicoID string, ano, mes *int, pagina, porPagina int) (*domain.PaginatedResponse[domain.Presenca], error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"politico_id": objectID}

	// TODO: Adicionar filtros de data quando necessário

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 50
	}

	skip := int64((pagina - 1) * porPagina)
	limit := int64(porPagina)

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "data", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var presencas []domain.Presenca
	if err := cursor.All(ctx, &presencas); err != nil {
		return nil, err
	}

	totalPaginas := int(total) / porPagina
	if int(total)%porPagina > 0 {
		totalPaginas++
	}

	return &domain.PaginatedResponse[domain.Presenca]{
		Data:         presencas,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}, nil
}

func (r *PresencaRepository) CalcularPercentual(ctx context.Context, politicoID string) (float64, error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return 0, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"politico_id": objectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":       nil,
			"total":     bson.M{"$sum": 1},
			"presentes": bson.M{"$sum": bson.M{"$cond": []interface{}{"$presente", 1, 0}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total     int `bson:"total"`
		Presentes int `bson:"presentes"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	if result.Total == 0 {
		return 0, nil
	}

	return float64(result.Presentes) / float64(result.Total) * 100, nil
}

