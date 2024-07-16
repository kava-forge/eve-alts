package esi

type AllianceData struct {
	Name   string `json:"name"`
	Ticker string `json:"ticker"`
}

type AllianceIcons struct {
	Medium string `json:"px128x128"`
	Small  string `json:"px64x64"`
}
