package entity

type IpLocation struct {
	Ip             string `json:"ip"`
	City           string `json:"city"`
	Timezone       string `json:"timezone"`
	AccuracyRadius uint16 `json:"accuracy_radius"`
	PostalCode     string `json:"postal_code"`
}
