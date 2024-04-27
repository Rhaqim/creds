package utils

import "fmt"

type KeyValue struct {
	Key   string      `json:"key" yaml:"key"`
	Value interface{} `json:"value" yaml:"value"`
}

// extractKeyValuePairs extracts key-value pairs recursively from a YAML map.
func extractKeyValuePairs(data map[string]interface{}) []KeyValue {
	var keyValues []KeyValue
	for key, value := range data {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			// Handle nested maps recursively
			nestedData := make(map[string]interface{})
			for k, val := range v {
				nestedData[fmt.Sprintf("%v", k)] = val
			}
			nestedKeyValues := extractKeyValuePairs(nestedData)
			keyValues = append(keyValues, nestedKeyValues...)
		case []interface{}:
			// Handle nested arrays recursively
			for i, item := range v {
				nestedData := map[string]interface{}{fmt.Sprintf("%d", i): item}
				nestedKeyValues := extractKeyValuePairs(nestedData)
				keyValues = append(keyValues, nestedKeyValues...)
			}
		default:
			// Add non-nested key-value pairs
			keyValues = append(keyValues, KeyValue{Key: key, Value: value})
		}
	}
	return keyValues
}
