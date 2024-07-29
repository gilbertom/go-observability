package entity

// Cep represents a CEP entity.
type Cep struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

type ResponseServiceB struct {
	City   string  `json:"city"`
	TempC  float64 `json:"tempC"`
	TempF  float64 `json:"tempF"`
	TempK  float64 `json:"tempK"`
}

// CEPRequest represents a request for a CEP (postal code).
type CEPRequest struct {
    CEP string `json:"cep"`
}
