package types

// Variable represents a Nomad variable
type Variable struct {
	Path        string            `json:"Path"`
	Value       string            `json:"Value,omitempty"`
	Items       map[string]string `json:"Items,omitempty"`
	Namespace   string            `json:"Namespace,omitempty"`
	CreateIndex int               `json:"CreateIndex,omitempty"`
	ModifyIndex int               `json:"ModifyIndex,omitempty"`
	CreateTime  int64             `json:"CreateTime,omitempty"`
	ModifyTime  int64             `json:"ModifyTime,omitempty"`
	CAS         *uint64           `json:"-"`
}
