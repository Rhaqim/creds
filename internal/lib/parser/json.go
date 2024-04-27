package lib

import (
	"encoding/json"
	"fmt"
	"log"
)

// KeyValue represents a key-value pair extracted from JSON.

func ParseJSON() {
	// Read the JSON file
	data := ReadFile("data.json")

	// Unmarshal the JSON data into a map
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Extract key-value pairs
	keyValues := extractKeyValuePairs(jsonData)

	// Print the extracted key-value pairs
	for _, kv := range keyValues {
		fmt.Printf("Key: %s, Value: %v\n", kv.Key, kv.Value)
	}
}

// extractKeyValuePairs extracts key-value pairs recursively from a JSON object.
// func extractKeyValuePairs(data map[string]interface{}) []KeyValue {
//     var keyValues []KeyValue
//     for key, value := range data {
//         switch v := value.(type) {
//         case map[string]interface{}:
//             // Handle nested objects recursively
//             nestedKeyValues := extractKeyValuePairs(v)
//             keyValues = append(keyValues, nestedKeyValues...)
//         default:
//             // Add non-nested key-value pairs
//             keyValues = append(keyValues, KeyValue{Key: key, Value: value})
//         }
//     }
//     return keyValues
// }
