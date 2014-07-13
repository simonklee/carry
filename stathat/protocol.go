package stathat

type ProtocolPackage struct {
	EZKey string         `json:"ezkey"`
	Data  []ProtocolStat `json:"data"`
}

type ProtocolStat struct {
	Stat      string   `json:"stat"`
	Value     *float64 `json:"value,omitempty"`
	Count     *int64   `json:"count,omitempty"`
	Timestamp int64    `json:"t"`
}
