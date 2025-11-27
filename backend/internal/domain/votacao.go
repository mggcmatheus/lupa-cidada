package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TipoVoto representa o tipo de voto
type TipoVoto string

const (
	VotoSim       TipoVoto = "SIM"
	VotoNao       TipoVoto = "NAO"
	VotoAbstencao TipoVoto = "ABSTENCAO"
	VotoAusente   TipoVoto = "AUSENTE"
	VotoObstrucao TipoVoto = "OBSTRUCAO"
)

// Votacao representa uma votação de um político
type Votacao struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PoliticoID   primitive.ObjectID `json:"politicoId" bson:"politico_id"`
	ProposicaoID primitive.ObjectID `json:"proposicaoId" bson:"proposicao_id"`
	Voto         TipoVoto           `json:"voto" bson:"voto"`
	Data         time.Time          `json:"data" bson:"data"`
	Sessao       string             `json:"sessao" bson:"sessao"`
}

// VotacaoComProposicao inclui os dados da proposição na votação
type VotacaoComProposicao struct {
	Votacao
	Proposicao *Proposicao `json:"proposicao,omitempty"`
}

