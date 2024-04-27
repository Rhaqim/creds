package models

// SaveFormat represents the format of the saved data.
type CredentialSaveFormat string

const (
	JSON  CredentialSaveFormat = "json"
	YAML  CredentialSaveFormat = "yaml"
	Plain CredentialSaveFormat = "plain"
)
