package app

// SaveFormat represents the format of the saved data.
type CredsSaveFormat string

const (
	JSON  CredsSaveFormat = "json"
	YAML  CredsSaveFormat = "yaml"
	Plain CredsSaveFormat = "plain"
)
