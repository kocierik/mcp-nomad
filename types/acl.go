package types

// ACLToken represents a Nomad ACL token
type ACLToken struct {
	AccessorID  string   `json:"AccessorID"`
	SecretID    string   `json:"SecretID"`
	Name        string   `json:"Name"`
	Type        string   `json:"Type"`
	Policies    []string `json:"Policies"`
	Global      bool     `json:"Global"`
	CreateIndex int      `json:"CreateIndex"`
	ModifyIndex int      `json:"ModifyIndex"`
}

// ACLPolicy represents a Nomad ACL policy
type ACLPolicy struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
	CreateIndex int    `json:"create_index"`
	ModifyIndex int    `json:"modify_index"`
}

// ACLRole represents a Nomad ACL role
type ACLRole struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Policies    []string `json:"policies"`
	CreateIndex int      `json:"create_index"`
	ModifyIndex int      `json:"modify_index"`
}

// ACLTokenList represents a list of ACL tokens
type ACLTokenList struct {
	Tokens []ACLToken `json:"tokens"`
}

// ACLPolicyList represents a list of ACL policies
type ACLPolicyList struct {
	Policies []ACLPolicy `json:"policies"`
}

// ACLRoleList represents a list of ACL roles
type ACLRoleList struct {
	Roles []ACLRole `json:"roles"`
}
