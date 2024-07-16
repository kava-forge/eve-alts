package esi

type CorporationData struct {
	AllianceID int64  `json:"alliance_id"`
	Name       string `json:"name"`
	Ticker     string `json:"ticker"`
}

type CorporationIcons struct {
	Large  string `json:"px256x256"`
	Medium string `json:"px128x128"`
	Small  string `json:"px64x64"`
}
