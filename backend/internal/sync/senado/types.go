package senado

import "time"

// API do Senado Federal
// Documentação: https://www12.senado.leg.br/dados-abertos

// SenadoresResponse representa a resposta da API de senadores
type SenadoresResponse struct {
	ListaParlamentarEmExercicio ListaParlamentar `json:"ListaParlamentarEmExercicio"`
}

// ListaParlamentar contém a lista de parlamentares
type ListaParlamentar struct {
	Parlamentares Parlamentares `json:"Parlamentares"`
}

// Parlamentares contém o array de parlamentares
type Parlamentares struct {
	Parlamentar []Parlamentar `json:"Parlamentar"`
}

// Parlamentar representa um senador
type Parlamentar struct {
	IdentificacaoParlamentar IdentificacaoParlamentar `json:"IdentificacaoParlamentar"`
	Mandato                  *Mandato                 `json:"Mandato,omitempty"`
}

// IdentificacaoParlamentar contém dados de identificação
type IdentificacaoParlamentar struct {
	CodigoParlamentar       string `json:"CodigoParlamentar"`
	NomeParlamentar         string `json:"NomeParlamentar"`
	NomeCompletoParlamentar string `json:"NomeCompletoParlamentar"`
	SexoParlamentar         string `json:"SexoParlamentar"`
	FormaTratamento         string `json:"FormaTratamento"`
	URLFotoParlamentar      string `json:"UrlFotoParlamentar"`
	URLPaginaParlamentar    string `json:"UrlPaginaParlamentar"`
	EmailParlamentar        string `json:"EmailParlamentar"`
	SiglaPartidoParlamentar string `json:"SiglaPartidoParlamentar"`
	UfParlamentar           string `json:"UfParlamentar"`
}

// Mandato contém informações do mandato
type Mandato struct {
	CodigoMandato         string       `json:"CodigoMandato"`
	UfParlamentar         string       `json:"UfParlamentar"`
	PrimeiraLegislatura   *Legislatura `json:"PrimeiraLegislaturaDoMandato"`
	SegundaLegislatura    *Legislatura `json:"SegundaLegislaturaDoMandato"`
	DescricaoParticipacao string       `json:"DescricaoParticipacao"`
}

// Legislatura contém dados da legislatura
type Legislatura struct {
	NumeroLegislatura string `json:"NumeroLegislatura"`
	DataInicio        string `json:"DataInicio"`
	DataFim           string `json:"DataFim"`
}

// SenadorDetalheResponse representa detalhes de um senador
type SenadorDetalheResponse struct {
	DetalheParlamentar DetalheParlamentar `json:"DetalheParlamentar"`
}

// DetalheParlamentar contém detalhes do parlamentar
type DetalheParlamentar struct {
	Parlamentar ParlamentarDetalhe `json:"Parlamentar"`
}

// ParlamentarDetalhe contém dados detalhados do parlamentar
type ParlamentarDetalhe struct {
	IdentificacaoParlamentar IdentificacaoParlamentar `json:"IdentificacaoParlamentar"`
	DadosBasicosParlamentar  DadosBasicos             `json:"DadosBasicosParlamentar"`
	Telefones                *Telefones               `json:"Telefones,omitempty"`
}

// DadosBasicos contém dados básicos do parlamentar
type DadosBasicos struct {
	DataNascimento      string `json:"DataNascimento"`
	Naturalidade        string `json:"Naturalidade"`
	UfNaturalidade      string `json:"UfNaturalidade"`
	EnderecoParlamentar string `json:"EnderecoParlamentar"`
}

// Telefones contém telefones de contato
type Telefones struct {
	Telefone []Telefone `json:"Telefone"`
}

// Telefone representa um telefone
type Telefone struct {
	NumeroTelefone  string `json:"NumeroTelefone"`
	OrdemPublicacao string `json:"OrdemPublicacao"`
}

// VotacoesSenadoResponse representa votações do Senado
type VotacoesSenadoResponse struct {
	ListaVotacoes ListaVotacoes `json:"ListaVotacoes"`
}

// ListaVotacoes contém a lista de votações
type ListaVotacoes struct {
	Votacoes Votacoes `json:"Votacoes"`
}

// Votacoes contém array de votações
type Votacoes struct {
	Votacao []VotacaoSenado `json:"Votacao"`
}

// VotacaoSenado representa uma votação
type VotacaoSenado struct {
	CodigoSessao                  string `json:"CodigoSessao"`
	SiglaCasa                     string `json:"SiglaCasa"`
	CodigoSessaoVotacao           string `json:"CodigoSessaoVotacao"`
	DataSessao                    string `json:"DataSessao"`
	HoraInicio                    string `json:"HoraInicio"`
	DescricaoIdentificacaoMateria string `json:"DescricaoIdentificacaoMateria"`
	DescricaoVotacao              string `json:"DescricaoVotacao"`
	Resultado                     string `json:"Resultado"`
}

// ParseDate converte string de data para time.Time
func ParseDate(dateStr string) time.Time {
	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}
