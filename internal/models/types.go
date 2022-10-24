package models

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Calculator struct {
	NumberOne    float64 `json:"number_first"`
	NumberSecond float64 `json:"number_second"`
	Operation    string  `json:"operation"`
	Result       float64 `json:"result"`
}

const (
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	MULTIPLY = "MULTIPLY"
	DIVIDE   = "DIVIDE"
)
