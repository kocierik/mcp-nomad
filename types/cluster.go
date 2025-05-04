package types

type RaftOperator struct {
	Address      string `json:"Address"`
	ID           string `json:"ID"`
	Leader       bool   `json:"Leader"`
	Node         string `json:"Node"`
	RaftProtocol string `json:"RaftProtocol"`
	Voter        bool   `json:"Voter"`
}
