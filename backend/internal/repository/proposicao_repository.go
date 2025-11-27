package repository

import (
	"context"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProposicaoRepository struct {
	collection *mongo.Collection
}

func NewProposicaoRepository(db *mongo.Database) *ProposicaoRepository {
	return &ProposicaoRepository{
		collection: db.Collection("proposicoes"),
	}
}

func (r *ProposicaoRepository) ListarPorAutor(ctx context.Context, autorID string, pagina, porPagina int) (*domain.PaginatedResponse[domain.Proposicao], error) {
	objectID, err := primitive.ObjectIDFromHex(autorID)
	if err != nil {
		return nil, err
	}

	// Busca por autor ou coautor
	filter := bson.M{
		"$or": []bson.M{
			{"autor_id": objectID},
			{"coautores_ids": objectID},
		},
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
		SetSort(bson.D{{Key: "ano", Value: -1}, {Key: "numero", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var proposicoes []domain.Proposicao
	if err := cursor.All(ctx, &proposicoes); err != nil {
		return nil, err
	}

	totalPaginas := int(total) / porPagina
	if int(total)%porPagina > 0 {
		totalPaginas++
	}

	return &domain.PaginatedResponse[domain.Proposicao]{
		Data:         proposicoes,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}, nil
}

func (r *ProposicaoRepository) ContarPorAutor(ctx context.Context, autorID string) (total, aprovadas int64, err error) {
	objectID, err := primitive.ObjectIDFromHex(autorID)
	if err != nil {
		return 0, 0, err
	}

	filter := bson.M{
		"$or": []bson.M{
			{"autor_id": objectID},
			{"coautores_ids": objectID},
		},
	}

	total, err = r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, 0, err
	}

	filterAprovadas := bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"autor_id": objectID},
				{"coautores_ids": objectID},
			}},
			{"situacao": domain.SituacaoAprovada},
		},
	}

	aprovadas, err = r.collection.CountDocuments(ctx, filterAprovadas)
	if err != nil {
		return 0, 0, err
	}

	return total, aprovadas, nil
}

func (r *ProposicaoRepository) Contar(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *ProposicaoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Proposicao, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var proposicao domain.Proposicao
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&proposicao)
	if err != nil {
		return nil, err
	}

	return &proposicao, nil
}

