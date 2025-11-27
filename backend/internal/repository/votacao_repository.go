package repository

import (
	"context"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VotacaoRepository struct {
	collection *mongo.Collection
}

func NewVotacaoRepository(db *mongo.Database) *VotacaoRepository {
	return &VotacaoRepository{
		collection: db.Collection("votacoes"),
	}
}

func (r *VotacaoRepository) ListarPorPolitico(ctx context.Context, politicoID string, pagina, porPagina int) (*domain.PaginatedResponse[domain.Votacao], error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"politico_id": objectID}

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

	var votacoes []domain.Votacao
	if err := cursor.All(ctx, &votacoes); err != nil {
		return nil, err
	}

	totalPaginas := int(total) / porPagina
	if int(total)%porPagina > 0 {
		totalPaginas++
	}

	return &domain.PaginatedResponse[domain.Votacao]{
		Data:         votacoes,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}, nil
}

func (r *VotacaoRepository) ContarPorPolitico(ctx context.Context, politicoID string) (map[domain.TipoVoto]int, error) {
	objectID, err := primitive.ObjectIDFromHex(politicoID)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"politico_id": objectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$voto",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[domain.TipoVoto]int)
	for cursor.Next(ctx) {
		var item struct {
			ID    domain.TipoVoto `bson:"_id"`
			Count int             `bson:"count"`
		}
		if err := cursor.Decode(&item); err != nil {
			continue
		}
		result[item.ID] = item.Count
	}

	return result, nil
}

func (r *VotacaoRepository) Contar(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

