package mock

import (
	"time"

	"github.com/lupa-cidada/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	id1 = primitive.NewObjectID()
	id2 = primitive.NewObjectID()
	id3 = primitive.NewObjectID()
	id4 = primitive.NewObjectID()
	id5 = primitive.NewObjectID()
	id6 = primitive.NewObjectID()
)

// Politicos retorna a lista de políticos mockados
func Politicos() []domain.Politico {
	return []domain.Politico{
		{
			ID:             id1,
			Nome:           "João Silva",
			NomeCivil:      "João Pedro da Silva Santos",
			FotoURL:        "",
			DataNascimento: time.Date(1965, 3, 15, 0, 0, 0, 0, time.UTC),
			Genero:         domain.GeneroMasculino,
			Partido: domain.Partido{
				Sigla: "PT",
				Nome:  "Partido dos Trabalhadores",
				Cor:   "#CC0000",
			},
			CargoAtual: domain.CargoAtual{
				Tipo:        domain.CargoDeputadoFederal,
				Esfera:      domain.EsferaFederal,
				Estado:      "SP",
				DataInicio:  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EmExercicio: true,
			},
			Contato: domain.Contato{
				Email:    "joao.silva@camara.leg.br",
				Telefone: "(61) 3215-5000",
				Gabinete: "Anexo IV, Gabinete 300",
			},
			RedesSociais: domain.RedesSociais{
				Twitter:   "@joaosilva",
				Instagram: "@dep.joaosilva",
			},
			SalarioBruto:   33763.00,
			SalarioLiquido: 25000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             id2,
			Nome:           "Maria Santos",
			NomeCivil:      "Maria das Graças Santos",
			FotoURL:        "",
			DataNascimento: time.Date(1970, 7, 22, 0, 0, 0, 0, time.UTC),
			Genero:         domain.GeneroFeminino,
			Partido: domain.Partido{
				Sigla: "PL",
				Nome:  "Partido Liberal",
				Cor:   "#003366",
			},
			CargoAtual: domain.CargoAtual{
				Tipo:        domain.CargoSenador,
				Esfera:      domain.EsferaFederal,
				Estado:      "RJ",
				DataInicio:  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EmExercicio: true,
			},
			Contato: domain.Contato{
				Email:    "maria.santos@senado.leg.br",
				Telefone: "(61) 3303-1000",
				Gabinete: "Anexo II, Gabinete 15",
			},
			RedesSociais: domain.RedesSociais{
				Twitter:   "@mariasantos",
				Instagram: "@sen.mariasantos",
			},
			SalarioBruto:   41650.92,
			SalarioLiquido: 30000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             id3,
			Nome:           "Pedro Costa",
			NomeCivil:      "Pedro Henrique Costa Lima",
			FotoURL:        "",
			DataNascimento: time.Date(1980, 11, 8, 0, 0, 0, 0, time.UTC),
			Genero:         domain.GeneroMasculino,
			Partido: domain.Partido{
				Sigla: "MDB",
				Nome:  "Movimento Democrático Brasileiro",
				Cor:   "#00AA00",
			},
			CargoAtual: domain.CargoAtual{
				Tipo:        domain.CargoDeputadoEstadual,
				Esfera:      domain.EsferaEstadual,
				Estado:      "MG",
				DataInicio:  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EmExercicio: true,
			},
			Contato: domain.Contato{
				Email: "pedro.costa@almg.gov.br",
			},
			SalarioBruto:   25322.25,
			SalarioLiquido: 18000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             id4,
			Nome:           "Ana Oliveira",
			NomeCivil:      "Ana Paula Oliveira Souza",
			FotoURL:        "",
			DataNascimento: time.Date(1975, 4, 30, 0, 0, 0, 0, time.UTC),
			Genero:         domain.GeneroFeminino,
			Partido: domain.Partido{
				Sigla: "PSOL",
				Nome:  "Partido Socialismo e Liberdade",
				Cor:   "#FFD700",
			},
			CargoAtual: domain.CargoAtual{
				Tipo:        domain.CargoVereador,
				Esfera:      domain.EsferaMunicipal,
				Estado:      "SP",
				Municipio:   "São Paulo",
				DataInicio:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				EmExercicio: true,
			},
			Contato: domain.Contato{
				Email: "ana.oliveira@saopaulo.sp.leg.br",
			},
			RedesSociais: domain.RedesSociais{
				Twitter: "@anaoliveira",
			},
			SalarioBruto:   18991.68,
			SalarioLiquido: 14000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             id5,
			Nome:           "Carlos Ferreira",
			NomeCivil:      "Carlos Alberto Ferreira",
			FotoURL:        "",
			DataNascimento: time.Date(1960, 9, 12, 0, 0, 0, 0, time.UTC),
			Genero:         domain.GeneroMasculino,
			Partido: domain.Partido{
				Sigla: "UNIÃO",
				Nome:  "União Brasil",
				Cor:   "#2E3092",
			},
			CargoAtual: domain.CargoAtual{
				Tipo:        domain.CargoGovernador,
				Esfera:      domain.EsferaEstadual,
				Estado:      "BA",
				DataInicio:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EmExercicio: true,
			},
			Contato: domain.Contato{
				Email: "gabinete@ba.gov.br",
			},
			SalarioBruto:   35462.22,
			SalarioLiquido: 26000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             id6,
			Nome:           "Fernanda Lima",
			NomeCivil:      "Fernanda Cristina Lima",
			FotoURL:        "",
			DataNascimento: time.Date(1985, 2, 28, 0, 0, 0, 0, time.UTC),
			Genero:         domain.GeneroFeminino,
			Partido: domain.Partido{
				Sigla: "NOVO",
				Nome:  "Partido Novo",
				Cor:   "#FF6600",
			},
			CargoAtual: domain.CargoAtual{
				Tipo:        domain.CargoDeputadoFederal,
				Esfera:      domain.EsferaFederal,
				Estado:      "RS",
				DataInicio:  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EmExercicio: true,
			},
			Contato: domain.Contato{
				Email:    "fernanda.lima@camara.leg.br",
				Telefone: "(61) 3215-5001",
			},
			RedesSociais: domain.RedesSociais{
				Instagram: "@fernanda.lima",
			},
			SalarioBruto:   33763.00,
			SalarioLiquido: 25000.00,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// Estatisticas retorna estatísticas mockadas para um político
func Estatisticas() map[string]domain.EstatisticasPolitico {
	return map[string]domain.EstatisticasPolitico{
		id1.Hex(): {
			TotalVotacoes:        245,
			VotosSim:             180,
			VotosNao:             45,
			Abstencoes:           10,
			Ausencias:            10,
			PercentualPresenca:   95.8,
			TotalProposicoes:     32,
			ProposicoesAprovadas: 8,
			TotalDespesas:        156000.00,
			MediaGastoMensal:     13000.00,
		},
		id2.Hex(): {
			TotalVotacoes:        189,
			VotosSim:             120,
			VotosNao:             50,
			Abstencoes:           15,
			Ausencias:            4,
			PercentualPresenca:   97.8,
			TotalProposicoes:     45,
			ProposicoesAprovadas: 12,
			TotalDespesas:        198000.00,
			MediaGastoMensal:     16500.00,
		},
		id3.Hex(): {
			TotalVotacoes:        156,
			VotosSim:             100,
			VotosNao:             40,
			Abstencoes:           8,
			Ausencias:            8,
			PercentualPresenca:   94.9,
			TotalProposicoes:     18,
			ProposicoesAprovadas: 5,
			TotalDespesas:        89000.00,
			MediaGastoMensal:     7416.67,
		},
		id4.Hex(): {
			TotalVotacoes:        312,
			VotosSim:             200,
			VotosNao:             80,
			Abstencoes:           20,
			Ausencias:            12,
			PercentualPresenca:   96.2,
			TotalProposicoes:     67,
			ProposicoesAprovadas: 23,
			TotalDespesas:        45000.00,
			MediaGastoMensal:     3750.00,
		},
		id5.Hex(): {
			TotalVotacoes:        0,
			VotosSim:             0,
			VotosNao:             0,
			Abstencoes:           0,
			Ausencias:            0,
			PercentualPresenca:   100,
			TotalProposicoes:     156,
			ProposicoesAprovadas: 89,
			TotalDespesas:        0,
			MediaGastoMensal:     0,
		},
		id6.Hex(): {
			TotalVotacoes:        245,
			VotosSim:             150,
			VotosNao:             70,
			Abstencoes:           15,
			Ausencias:            10,
			PercentualPresenca:   95.9,
			TotalProposicoes:     28,
			ProposicoesAprovadas: 6,
			TotalDespesas:        120000.00,
			MediaGastoMensal:     10000.00,
		},
	}
}

// GetPoliticoByID retorna um político pelo ID
func GetPoliticoByID(id string) *domain.Politico {
	for _, p := range Politicos() {
		if p.ID.Hex() == id {
			return &p
		}
	}
	return nil
}

// GetEstatisticasByID retorna estatísticas de um político pelo ID
func GetEstatisticasByID(id string) *domain.EstatisticasPolitico {
	stats := Estatisticas()
	if s, ok := stats[id]; ok {
		return &s
	}
	// Retorna estatísticas padrão se não encontrar
	return &domain.EstatisticasPolitico{
		TotalVotacoes:        100,
		VotosSim:             60,
		VotosNao:             30,
		Abstencoes:           5,
		Ausencias:            5,
		PercentualPresenca:   95.0,
		TotalProposicoes:     10,
		ProposicoesAprovadas: 3,
		TotalDespesas:        50000.00,
		MediaGastoMensal:     5000.00,
	}
}
