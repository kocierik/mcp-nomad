// File: types/namespaces.go
package types

// Namespace represents a Nomad namespace
type Namespace struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
