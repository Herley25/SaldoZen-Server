package models

type Summary struct {
	TotalDespesas float64 `json:"total_despesas"`
	TotalPagas    float64 `json:"total_pagas"`
	Pendentes     float64 `json:"pendentes"`
	TotalVencidas float64 `json:"total_vencidas"`
	Saldo         float64 `json:"saldo"`
	Mes           int     `json:"mes"`
	Ano           int     `json:"ano"`
}
