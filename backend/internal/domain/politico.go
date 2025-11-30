package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tipos enum
type Cargo string

const (
	CargoDeputadoFederal   Cargo = "DEPUTADO_FEDERAL"
	CargoSenador           Cargo = "SENADOR"
	CargoDeputadoEstadual  Cargo = "DEPUTADO_ESTADUAL"
	CargoDeputadoDistrital Cargo = "DEPUTADO_DISTRITAL"
	CargoVereador          Cargo = "VEREADOR"
	CargoPrefeito          Cargo = "PREFEITO"
	CargoGovernador        Cargo = "GOVERNADOR"
	CargoPresidente        Cargo = "PRESIDENTE"
)

type Esfera string

const (
	EsferaFederal   Esfera = "FEDERAL"
	EsferaEstadual  Esfera = "ESTADUAL"
	EsferaMunicipal Esfera = "MUNICIPAL"
)

type Genero string

const (
	GeneroMasculino Genero = "M"
	GeneroFeminino  Genero = "F"
	GeneroOutro     Genero = "OUTRO"
)

// Partido representa um partido político
type Partido struct {
	Sigla string `json:"sigla" bson:"sigla"`
	Nome  string `json:"nome" bson:"nome"`
	Cor   string `json:"cor" bson:"cor"`
}

// CargoAtual representa o cargo atual do político
type CargoAtual struct {
	Tipo        Cargo     `json:"tipo" bson:"tipo"`
	Esfera      Esfera    `json:"esfera" bson:"esfera"`
	Estado      string    `json:"estado" bson:"estado"`
	Municipio   string    `json:"municipio,omitempty" bson:"municipio,omitempty"`
	DataInicio  time.Time `json:"dataInicio" bson:"data_inicio"`
	DataFim     time.Time `json:"dataFim,omitempty" bson:"data_fim,omitempty"`
	EmExercicio bool      `json:"emExercicio" bson:"em_exercicio"`
}

// Contato representa as informações de contato
type Contato struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Telefone string `json:"telefone,omitempty" bson:"telefone,omitempty"`
	Gabinete string `json:"gabinete,omitempty" bson:"gabinete,omitempty"`
}

// RedesSociais representa os links das redes sociais
type RedesSociais struct {
	Twitter   string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Instagram string `json:"instagram,omitempty" bson:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty" bson:"facebook,omitempty"`
}

// Politico representa um político no sistema
type Politico struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CPF                 string             `json:"cpf,omitempty" bson:"cpf,omitempty"`
	Nome                string             `json:"nome" bson:"nome"`
	NomeCivil           string             `json:"nomeCivil" bson:"nome_civil"`
	NomeEleitoral       string             `json:"nomeEleitoral,omitempty" bson:"nome_eleitoral,omitempty"`
	FotoURL             string             `json:"fotoUrl" bson:"foto_url"`
	DataNascimento      time.Time          `json:"dataNascimento" bson:"data_nascimento"`
	Genero              Genero             `json:"genero" bson:"genero"`
	Partido             Partido            `json:"partido" bson:"partido"`
	CargoAtual          CargoAtual         `json:"cargoAtual" bson:"cargo_atual"`
	HistoricoCargos     []CargoAtual       `json:"historicoCargos" bson:"historico_cargos"`
	Contato             Contato            `json:"contato" bson:"contato"`
	RedesSociais        RedesSociais       `json:"redesSociais" bson:"redes_sociais"`
	SalarioBruto        float64            `json:"salarioBruto" bson:"salario_bruto"`
	SalarioLiquido      float64            `json:"salarioLiquido" bson:"salario_liquido"`
	Escolaridade        string             `json:"escolaridade,omitempty" bson:"escolaridade,omitempty"`
	MunicipioNascimento string             `json:"municipioNascimento,omitempty" bson:"municipio_nascimento,omitempty"`
	UFNascimento        string             `json:"ufNascimento,omitempty" bson:"uf_nascimento,omitempty"`
	Website             string             `json:"website,omitempty" bson:"website,omitempty"`
	IDExternoCamara     int                `json:"idExternoCamara,omitempty" bson:"id_externo_camara,omitempty"` // ID da API da Câmara
	CreatedAt           time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updated_at"`
}

// EstatisticasPolitico representa as estatísticas agregadas de um político
type EstatisticasPolitico struct {
	TotalVotacoes        int     `json:"totalVotacoes"`
	VotosSim             int     `json:"votosSim"`
	VotosNao             int     `json:"votosNao"`
	Abstencoes           int     `json:"abstencoes"`
	Ausencias            int     `json:"ausencias"`
	PercentualPresenca   float64 `json:"percentualPresenca"`
	TotalProposicoes     int     `json:"totalProposicoes"`
	ProposicoesAprovadas int     `json:"proposicoesAprovadas"`
	TotalDespesas        float64 `json:"totalDespesas"`
	MediaGastoMensal     float64 `json:"mediaGastoMensal"`
}

// FiltrosPoliticos representa os filtros disponíveis para busca
type FiltrosPoliticos struct {
	Nome              string   `query:"nome"`
	Partido           []string `query:"partido"`
	Cargo             []Cargo  `query:"cargo"`
	Esfera            []Esfera `query:"esfera"`
	Estado            []string `query:"estado"`
	Municipio         []string `query:"municipio"`
	EmExercicio       *bool    `query:"emExercicio"`
	Genero            []Genero `query:"genero"`
	IdadeMinima       *int     `query:"idadeMinima"`
	IdadeMaxima       *int     `query:"idadeMaxima"`
	PresencaMinima    *float64 `query:"presencaMinima"`
	ProposicoesMinima *int     `query:"proposicoesMinima"`
	OrdenarPor        string   `query:"ordenarPor"`
	Ordem             string   `query:"ordem"`
	Pagina            int      `query:"pagina"`
	PorPagina         int      `query:"porPagina"`
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse[T any] struct {
	Data         []T   `json:"data"`
	Total        int64 `json:"total"`
	Pagina       int   `json:"pagina"`
	PorPagina    int   `json:"porPagina"`
	TotalPaginas int   `json:"totalPaginas"`
}
