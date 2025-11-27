package repository

import (
	"context"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DespesaRepository struct {
	collection *mongo.Collection
}

func NewDespesaRepository(db *mongo.Database) *DespesaRepository {
	return &DespesaRepository{
		collection: db.Collection("despesas"),
	}
}

func (r *DespesaRepository) ListarPorPolitico(ctx context.Context, politicoID string, ano, mes *int, pagina, porPagina int) (*domain.PaginatedResponse[domain.Despesa], error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"politico_id": objectID}

	if ano != nil {
		filter["ano_referencia"] = *ano
	}
	if mes != nil {
		filter["mes_referencia"] = *mes
	}

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 20
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

	var despesas []domain.Despesa
	if err := cursor.All(ctx, &despesas); err != nil {
		return nil, err
	}

	totalPaginas := int(total) / porPagina
	if int(total)%porPagina > 0 {
		totalPaginas++
	}

	return &domain.PaginatedResponse[domain.Despesa]{
		Data:         despesas,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}, nil
}

func (r *DespesaRepository) TotalPorPolitico(ctx context.Context, politicoID string) (float64, error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return 0, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"politico_id": objectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$valor"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total float64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.Total, nil
}

func (r *DespesaRepository) MediaMensalPorPolitico(ctx context.Context, politicoID string) (float64, error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return 0, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"politico_id": objectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   bson.M{"ano": "$ano_referencia", "mes": "$mes_referencia"},
			"total": bson.M{"$sum": "$valor"},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"media": bson.M{"$avg": "$total"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Media float64 `bson:"media"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.Media, nil
}

func (r *DespesaRepository) TotalGeral(ctx context.Context) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$valor"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total float64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.Total, nil
}

