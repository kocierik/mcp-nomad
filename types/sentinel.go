package types

// SentinelPolicy represents a Nomad Sentinel policy
type SentinelPolicy struct {
	Name             string `json:"Name"`
	Description      string `json:"Description"`
	Scope            string `json:"Scope"`
	EnforcementLevel string `json:"EnforcementLevel"`
	Policy           string `json:"Policy"`
	Hash             string `json:"Hash,omitempty"`
	CreateIndex      int    `json:"CreateIndex,omitempty"`
	ModifyIndex      int    `json:"ModifyIndex,omitempty"`
}
