package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Despesa representa uma despesa de um político
type Despesa struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PoliticoID      primitive.ObjectID `json:"politicoId" bson:"politico_id"`
	Tipo            string             `json:"tipo" bson:"tipo"`
	Descricao       string             `json:"descricao" bson:"descricao"`
	Fornecedor      string             `json:"fornecedor" bson:"fornecedor"`
	CNPJFornecedor  string             `json:"cnpjFornecedor" bson:"cnpj_fornecedor"`
	Valor           float64            `json:"valor" bson:"valor"`
	Data            time.Time          `json:"data" bson:"data"`
	MesReferencia   int                `json:"mesReferencia" bson:"mes_referencia"`
	AnoReferencia   int                `json:"anoReferencia" bson:"ano_referencia"`
	DocumentoURL    string             `json:"documentoUrl,omitempty" bson:"documento_url,omitempty"`
}

// Presenca representa a presença de um político em uma sessão
type Presenca struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PoliticoID primitive.ObjectID `json:"politicoId" bson:"politico_id"`
	Data       time.Time          `json:"data" bson:"data"`
	TipoSessao string             `json:"tipoSessao" bson:"tipo_sessao"`
	Presente   bool               `json:"presente" bson:"presente"`
}

