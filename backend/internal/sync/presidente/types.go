package presidente

import "time"

// Dados do Presidente da República
// Fonte: Portal da Transparência, TSE, Presidência da República

// PresidenteData representa os dados do presidente atual
type PresidenteData struct {
	Nome           string
	NomeCivil      string
	CPF            string
	DataNascimento time.Time
	Genero         string
	Partido        string
	Estado         string
	DataInicio     time.Time
	DataFim        time.Time
	EmExercicio    bool
	FotoURL        string
	Email          string
	Telefone       string
}

// ParseDate converte string de data para time.Time
func ParseDate(dateStr string) time.Time {
	layouts := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"02/01/2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}

