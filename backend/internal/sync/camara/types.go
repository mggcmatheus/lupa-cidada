package camara

import "time"

// API da Câmara dos Deputados
// Documentação: https://dadosabertos.camara.leg.br/swagger/api.html

// DeputadoResponse representa a resposta da API de deputados
type DeputadoResponse struct {
	Dados []DeputadoResumo `json:"dados"`
	Links []Link           `json:"links"`
}

// DeputadoResumo representa um deputado na listagem
type DeputadoResumo struct {
	ID            int    `json:"id"`
	URI           string `json:"uri"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"siglaPartido"`
	URIPartido    string `json:"uriPartido"`
	SiglaUF       string `json:"siglaUf"`
	IDLegislatura int    `json:"idLegislatura"`
	URLFoto       string `json:"urlFoto"`
	Email         string `json:"email"`
}

// DeputadoDetalheResponse representa a resposta de detalhes de um deputado
type DeputadoDetalheResponse struct {
	Dados DeputadoDetalhe `json:"dados"`
}

// DeputadoDetalhe representa os detalhes completos de um deputado
type DeputadoDetalhe struct {
	ID                  int          `json:"id"`
	URI                 string       `json:"uri"`
	NomeCivil           string       `json:"nomeCivil"`
	UltimoStatus        UltimoStatus `json:"ultimoStatus"`
	CPF                 string       `json:"cpf"`
	Sexo                string       `json:"sexo"`
	URLWebsite          string       `json:"urlWebsite"`
	RedeSocial          []string     `json:"redeSocial"`
	DataNascimento      string       `json:"dataNascimento"`
	DataFalecimento     string       `json:"dataFalecimento"`
	UfNascimento        string       `json:"ufNascimento"`
	MunicipioNascimento string       `json:"municipioNascimento"`
	Escolaridade        string       `json:"escolaridade"`
}

// UltimoStatus representa o status atual do deputado
type UltimoStatus struct {
	ID                int       `json:"id"`
	URI               string    `json:"uri"`
	Nome              string    `json:"nome"`
	SiglaPartido      string    `json:"siglaPartido"`
	URIPartido        string    `json:"uriPartido"`
	SiglaUF           string    `json:"siglaUf"`
	IDLegislatura     int       `json:"idLegislatura"`
	URLFoto           string    `json:"urlFoto"`
	Email             string    `json:"email"`
	Data              string    `json:"data"`
	NomeEleitoral     string    `json:"nomeEleitoral"`
	Gabinete          *Gabinete `json:"gabinete"`
	Situacao          string    `json:"situacao"`
	CondicaoEleitoral string    `json:"condicaoEleitoral"`
	Descricao         string    `json:"descricaoStatus"`
}

// Gabinete representa informações do gabinete
type Gabinete struct {
	Nome     string `json:"nome"`
	Predio   string `json:"predio"`
	Sala     string `json:"sala"`
	Andar    string `json:"andar"`
	Telefone string `json:"telefone"`
	Email    string `json:"email"`
}

// DespesasResponse representa a resposta de despesas
type DespesasResponse struct {
	Dados []Despesa `json:"dados"`
	Links []Link    `json:"links"`
}

// Despesa representa uma despesa de um deputado
type Despesa struct {
	Ano               int     `json:"ano"`
	Mes               int     `json:"mes"`
	TipoDespesa       string  `json:"tipoDespesa"`
	CodDocumento      int     `json:"codDocumento"`
	TipoDocumento     string  `json:"tipoDocumento"`
	CodTipoDocumento  int     `json:"codTipoDocumento"`
	DataDocumento     string  `json:"dataDocumento"`
	NumDocumento      string  `json:"numDocumento"`
	ValorDocumento    float64 `json:"valorDocumento"`
	URLDocumento      string  `json:"urlDocumento"`
	NomeFornecedor    string  `json:"nomeFornecedor"`
	CNPJCPFFornecedor string  `json:"cnpjCpfFornecedor"`
	ValorLiquido      float64 `json:"valorLiquido"`
	ValorGlosa        float64 `json:"valorGlosa"`
	NumRessarcimento  string  `json:"numRessarcimento"`
	CodLote           int     `json:"codLote"`
	Parcela           int     `json:"parcela"`
}

// VotacoesResponse representa a resposta de votações
type VotacoesResponse struct {
	Dados []Votacao `json:"dados"`
	Links []Link    `json:"links"`
}

// Votacao representa uma votação
type Votacao struct {
	ID               string             `json:"id"`
	URI              string             `json:"uri"`
	Data             string             `json:"data"`
	DataHoraRegistro string             `json:"dataHoraRegistro"`
	SiglaOrgao       string             `json:"siglaOrgao"`
	URIOrgao         string             `json:"uriOrgao"`
	URIEvento        string             `json:"uriEvento"`
	Proposicao       *ProposicaoVotacao `json:"proposicaoObjeto"`
	URIProposicao    string             `json:"uriProposicaoObjeto"`
	Descricao        string             `json:"descricao"`
	Aprovacao        int                `json:"aprovacao"`
}

// ProposicaoVotacao representa a proposição votada
type ProposicaoVotacao struct {
	ID     int    `json:"id"`
	URI    string `json:"uri"`
	Siglum string `json:"siglaTipo"`
	Numero int    `json:"numero"`
	Ano    int    `json:"ano"`
	Ementa string `json:"ementa"`
}

// VotoDeputadoResponse representa os votos de um deputado em uma votação
type VotoDeputadoResponse struct {
	Dados []VotoDeputado `json:"dados"`
}

// VotoDeputado representa o voto de um deputado
type VotoDeputado struct {
	TipoVoto         string       `json:"tipoVoto"`
	DataRegistroVoto string       `json:"dataRegistroVoto"`
	Deputado         DeputadoVoto `json:"deputado_"`
}

// DeputadoVoto representa o deputado que votou
type DeputadoVoto struct {
	ID           int    `json:"id"`
	URI          string `json:"uri"`
	Nome         string `json:"nome"`
	SiglaPartido string `json:"siglaPartido"`
	URIPartido   string `json:"uriPartido"`
	SiglaUF      string `json:"siglaUf"`
	URLFoto      string `json:"urlFoto"`
}

// ProposicoesResponse representa a resposta de proposições
type ProposicoesResponse struct {
	Dados []Proposicao `json:"dados"`
	Links []Link       `json:"links"`
}

// Proposicao representa uma proposição
type Proposicao struct {
	ID        int    `json:"id"`
	URI       string `json:"uri"`
	SiglaTipo string `json:"siglaTipo"`
	CodTipo   int    `json:"codTipo"`
	Numero    int    `json:"numero"`
	Ano       int    `json:"ano"`
	Ementa    string `json:"ementa"`
}

// ProposicaoDetalheResponse representa detalhes de uma proposição
type ProposicaoDetalheResponse struct {
	Dados ProposicaoDetalhe `json:"dados"`
}

// ProposicaoDetalhe representa detalhes completos de uma proposição
type ProposicaoDetalhe struct {
	ID                int               `json:"id"`
	URI               string            `json:"uri"`
	SiglaTipo         string            `json:"siglaTipo"`
	CodTipo           int               `json:"codTipo"`
	Numero            int               `json:"numero"`
	Ano               int               `json:"ano"`
	Ementa            string            `json:"ementa"`
	DataApresentacao  string            `json:"dataApresentacao"`
	URIOrgaoNumerador string            `json:"uriOrgaoNumerador"`
	StatusProposicao  *StatusProposicao `json:"statusProposicao"`
	URIAutores        string            `json:"uriAutores"`
	DescricaoTipo     string            `json:"descricaoTipo"`
	EmentaDetalhada   string            `json:"ementaDetalhada"`
	Keywords          string            `json:"keywords"`
	URIPropPrincipal  string            `json:"uriPropPrincipal"`
	URIPropAnterior   string            `json:"uriPropAnterior"`
	URIPropPosterior  string            `json:"uriPropPosterior"`
	URLInteiroTeor    string            `json:"urlInteiroTeor"`
	URNFinal          string            `json:"urnFinal"`
	Texto             string            `json:"texto"`
	Justificativa     string            `json:"justificativa"`
}

// StatusProposicao representa o status atual da proposição
type StatusProposicao struct {
	DataHora          string `json:"dataHora"`
	Sequencia         int    `json:"sequencia"`
	SiglaOrgao        string `json:"siglaOrgao"`
	URIOrgao          string `json:"uriOrgao"`
	Regime            string `json:"regime"`
	DescTramitacao    string `json:"descricaoTramitacao"`
	CodTipoTramitacao string `json:"codTipoTramitacao"`
	DescSituacao      string `json:"descricaoSituacao"`
	CodSituacao       int    `json:"codSituacao"`
	Despacho          string `json:"despacho"`
	URL               string `json:"url"`
	Ambito            string `json:"ambito"`
}

// Link representa um link de paginação
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// ParseDate converte string de data para time.Time
func ParseDate(dateStr string) time.Time {
	layouts := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}
