package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SituacaoProposicao representa a situação de uma proposição
type SituacaoProposicao string

const (
	SituacaoEmTramitacao SituacaoProposicao = "EM_TRAMITACAO"
	SituacaoAprovada     SituacaoProposicao = "APROVADA"
	SituacaoRejeitada    SituacaoProposicao = "REJEITADA"
	SituacaoArquivada    SituacaoProposicao = "ARQUIVADA"
	SituacaoRetirada     SituacaoProposicao = "RETIRADA"
)

// TramitacaoItem representa um item de tramitação
type TramitacaoItem struct {
	Data      time.Time `json:"data" bson:"data"`
	Descricao string    `json:"descricao" bson:"descricao"`
	Orgao     string    `json:"orgao" bson:"orgao"`
}

// Proposicao representa uma proposição legislativa
type Proposicao struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Tipo         string               `json:"tipo" bson:"tipo"`
	Numero       string               `json:"numero" bson:"numero"`
	Ano          int                  `json:"ano" bson:"ano"`
	Ementa       string               `json:"ementa" bson:"ementa"`
	AutorID      primitive.ObjectID   `json:"autorId" bson:"autor_id"`
	CoautoresIDs []primitive.ObjectID `json:"coautoresIds" bson:"coautores_ids"`
	Situacao     SituacaoProposicao   `json:"situacao" bson:"situacao"`
	Tema         []string             `json:"tema" bson:"tema"`
	Tramitacao   []TramitacaoItem     `json:"tramitacao" bson:"tramitacao"`
	CreatedAt    time.Time            `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time            `json:"updatedAt" bson:"updated_at"`
}

// ProposicaoComAutor inclui os dados do autor na proposição
type ProposicaoComAutor struct {
	Proposicao
	Autor *Politico `json:"autor,omitempty"`
}

