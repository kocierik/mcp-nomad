package types

// Variable represents a Nomad variable
type Variable struct {
	Path      string `json:"Path"`
	Value     string `json:"Value"`
	Namespace string `json:"Namespace"`
}
