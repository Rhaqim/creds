package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type KeyValueStr struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// ConvertJSONToKeyValueStr converts JSON data to a slice of KeyValueStr structs.
func ConvertJSONToKeyValueStr(data []byte) ([]KeyValueStr, error) {
	var result []KeyValueStr
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	result = extractKeyValuePairs("", jsonData)
	return result, nil
}

// ConvertYAMLToKeyValueStr converts YAML data to a slice of KeyValueStr structs.
func ConvertYAMLToKeyValueStr(data []byte) ([]KeyValueStr, error) {
	var result []KeyValueStr
	var yamlData map[interface{}]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}
	result = extractKeyValuePairs("", yamlData)
	return result, nil
}

// extractKeyValuePairs recursively extracts key-value pairs from JSON or YAML data.
func extractKeyValuePairs(prefix string, data interface{}) []KeyValueStr {
	var result []KeyValueStr
	switch d := data.(type) {
	case map[interface{}]interface{}:
		for k, v := range d {
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "->"
			}
			result = append(result, extractKeyValuePairs(subPrefix+fmt.Sprintf("%v", k), v)...)
		}
	case map[string]interface{}:
		for k, v := range d {
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "->"
			}
			result = append(result, extractKeyValuePairs(subPrefix+k, v)...)
		}
	case []interface{}:
		for i, v := range d {
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "->"
			}
			result = append(result, extractKeyValuePairs(subPrefix+fmt.Sprintf("%d", i), v)...)
		}
	default:
		result = append(result, KeyValueStr{
			Key:   prefix,
			Value: fmt.Sprintf("%v", d),
		})
	}
	return result
}

// ConvertKeyValueStrToOriginal converts a slice of KeyValueStr structs back to the original data structure.
func ConvertKeyValueStrToOriginal2(keyValues []KeyValueStr) (interface{}, error) {
	dataMap := make(map[string]interface{})
	for _, kv := range keyValues {
		keys := strings.Split(kv.Key, "->")
		lastIndex := len(keys) - 1
		currentMap := dataMap
		for i, key := range keys {
			if i == lastIndex {
				// If it's the last key, assign the value
				currentMap[key] = kv.Value
			} else {
				// If it's not the last key, create a nested map if it doesn't exist
				if _, ok := currentMap[key]; !ok {
					currentMap[key] = make(map[string]interface{})
				}
				// Move to the next level
				currentMap = currentMap[key].(map[string]interface{})
			}
		}
	}
	// Convert map to appropriate type (map, slice, etc.)
	data := convertMapToType2(dataMap)
	return data, nil
}

// convertMapToType converts a map[string]interface{} to the appropriate data structure.
func convertMapToType2(dataMap map[string]interface{}) interface{} {
	result := make(map[string]interface{})
	for key, value := range dataMap {
		if strings.Contains(key, "->") {
			// This is a nested key
			keys := strings.Split(key, "->")
			lastIndex := len(keys) - 1
			currentMap := result
			for i, k := range keys {
				if i == lastIndex {
					currentMap[k] = value
				} else {
					if _, ok := currentMap[k]; !ok {
						currentMap[k] = make(map[string]interface{})
					}
					currentMap = currentMap[k].(map[string]interface{})
				}
			}
		} else {
			// Regular key-value pair
			result[key] = value
		}
	}
	// If there's only one key and it's not nested, return its value directly
	if len(result) == 1 {
		for _, v := range result {
			return v
		}
	}
	return result
}

func Check2() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Construct the absolute path to the test.yaml file
	filePath := filepath.Join(currentDir, "internal", "tests", "data", "file")

	// Read JSON file
	jsonData, err := os.ReadFile(filepath.Join(filePath, "test.json"))
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Convert JSON to KeyValueStr
	keyValuesJSON, err := ConvertJSONToKeyValueStr(jsonData)
	if err != nil {
		log.Fatalf("Error converting JSON to KeyValueStr: %v", err)
	}
	fmt.Println("JSON to KeyValueStr:")
	for _, kv := range keyValuesJSON {
		fmt.Printf("%s: %s\n", kv.Key, kv.Value)
	}

	// Read YAML file
	yamlData, err := os.ReadFile(filepath.Join(filePath, "test.yaml"))
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Convert YAML to KeyValueStr
	keyValuesYAML, err := ConvertYAMLToKeyValueStr(yamlData)
	if err != nil {
		log.Fatalf("Error converting YAML to KeyValueStr: %v", err)
	}
	fmt.Println("\nYAML to KeyValueStr:")
	for _, kv := range keyValuesYAML {
		fmt.Printf("%s: %s\n", kv.Key, kv.Value)
	}

	// Convert KeyValueStr to original data structure
	originalDataJSON, err := ConvertKeyValueStrToOriginal2(keyValuesJSON)
	if err != nil {
		log.Fatalf("Error converting KeyValueStr to original data structure: %v", err)
	}

	fmt.Println("\nKeyValueStr to JSON:")
	fmt.Println(originalDataJSON)

	originalDataYAML, err := ConvertKeyValueStrToOriginal2(keyValuesYAML)
	if err != nil {
		log.Fatalf("Error converting KeyValueStr to original data structure: %v", err)
	}

	fmt.Println("\nKeyValueStr to YAML:")
	fmt.Println(originalDataYAML)

	// write to a file
	// Convert original data structure to JSON
	jsonData, err = json.Marshal(originalDataJSON)
	if err != nil {
		log.Fatalf("Error converting original data structure to JSON: %v", err)
	}
	err = os.WriteFile(filepath.Join(filePath, "test_output.json"), jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}

	// Convert original data structure to YAML
	yamlData, err = yaml.Marshal(originalDataYAML)
	if err != nil {
		log.Fatalf("Error converting original data structure to YAML: %v", err)
	}
	err = os.WriteFile(filepath.Join(filePath, "test_output.yaml"), yamlData, 0644)
	if err != nil {
		log.Fatalf("Error writing YAML file: %v", err)
	}
}
