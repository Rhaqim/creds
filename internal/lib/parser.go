package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type FileParser struct {
	// FileFormat is the format of the file to be parsed.
	FileFormat string `json:"file_format" yaml:"file_format"`
	// FileData is the raw data of the file to be parsed.
	FileData []byte `json:"file_data" yaml:"file_data"`
}

type KeyValueStr struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

func (O *FileParser) Parse() []KeyValueStr {
	switch strings.ToLower(O.FileFormat) {
	case "yaml":
		return O.ParseYAML()
	case "json":
		return O.ParseJSON()
	case "plain":
		return O.ParsePlain()
	default:
		log.Fatalf("Unsupported file format: %s", O.FileFormat)
		return nil
	}
}

func (O *FileParser) ParseYAML() []KeyValueStr {
	parsed, err := O.ConvertYAMLToKeyValueStr(O.FileData)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	return parsed
}

func (O *FileParser) ParseJSON() []KeyValueStr {
	parsed, err := O.ConvertJSONToKeyValueStr(O.FileData)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return parsed
}

func (O *FileParser) ParsePlain() []KeyValueStr {
	_ = bufio.NewScanner(strings.NewReader(string(O.FileData)))
	return []KeyValueStr{}
	// return O.ExtractKeyValuePairsPlain(scanner)
}

// ConvertJSONToKeyValueStr converts JSON data to a slice of KeyValueStr structs.
func (O *FileParser) ConvertJSONToKeyValueStr(data []byte) ([]KeyValueStr, error) {
	var result []KeyValueStr
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	result = O.ExtractKeyValuePairs("", jsonData)
	return result, nil
}

// ConvertYAMLToKeyValueStr converts YAML data to a slice of KeyValueStr structs.
func (O *FileParser) ConvertYAMLToKeyValueStr(data []byte) ([]KeyValueStr, error) {
	var result []KeyValueStr
	var yamlData map[interface{}]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}
	result = O.ExtractKeyValuePairs("", yamlData)
	return result, nil
}

// extractKeyValuePairs recursively extracts key-value pairs from JSON or YAML data.
func (O *FileParser) ExtractKeyValuePairs(prefix string, data interface{}) []KeyValueStr {
	var result []KeyValueStr
	switch d := data.(type) {
	case map[interface{}]interface{}:
		for k, v := range d {
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "->"
			}
			result = append(result, O.ExtractKeyValuePairs(subPrefix+fmt.Sprintf("%v", k), v)...)
		}
	case map[string]interface{}:
		for k, v := range d {
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "->"
			}
			result = append(result, O.ExtractKeyValuePairs(subPrefix+k, v)...)
		}
	case []interface{}:
		for i, v := range d {
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "->"
			}
			result = append(result, O.ExtractKeyValuePairs(subPrefix+fmt.Sprintf("%d", i), v)...)
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
func (O *FileParser) ConvertKeyValueStrToOriginal(keyValues []KeyValueStr) (interface{}, error) {
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
	data := O.convertMapToType(dataMap)
	return data, nil
}

// convertMapToType converts a map[string]interface{} to the appropriate data structure.
func (O *FileParser) convertMapToType(dataMap map[string]interface{}) interface{} {
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

// // ExtractKeyValuePairs extracts key-value pairs recursively from a YAML map.
// func (O *FileParser) ExtractKeyValuePairs(data map[string]interface{}) []KeyValue { //FIXME: This function only works for yaml files
// 	var keyValues []KeyValue
// 	for key, value := range data {
// 		switch v := value.(type) {
// 		case map[interface{}]interface{}:
// 			// Handle nested maps recursively
// 			nestedData := make(map[string]interface{})
// 			for k, val := range v {
// 				nestedData[fmt.Sprintf("%v", k)] = val
// 			}
// 			nestedKeyValues := O.ExtractKeyValuePairs(nestedData)
// 			keyValues = append(keyValues, nestedKeyValues...)
// 		case []interface{}:
// 			// Handle nested arrays recursively
// 			for i, item := range v {
// 				nestedData := map[string]interface{}{fmt.Sprintf("%d", i): item}
// 				nestedKeyValues := O.ExtractKeyValuePairs(nestedData)
// 				keyValues = append(keyValues, nestedKeyValues...)
// 			}
// 		default:
// 			// Add non-nested key-value pairs
// 			keyValues = append(keyValues, KeyValue{Key: key, Value: value})
// 		}
// 	}
// 	return keyValues
// }

// // ExtractKeyValuePairs extracts key-value pairs from a scanner reading a plain file.
// func (O *FileParser) ExtractKeyValuePairsPlain(scanner *bufio.Scanner) []KeyValue {
// 	var keyValues []KeyValue
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		line = strings.TrimSpace(line)
// 		if line == "" || strings.HasPrefix(line, "#") {
// 			// Skip empty lines or lines starting with #
// 			continue
// 		}
// 		parts := strings.SplitN(line, "=", 2)
// 		if len(parts) != 2 {
// 			log.Printf("Invalid line: %s", line)
// 			continue
// 		}
// 		key := strings.TrimSpace(parts[0])
// 		value := strings.TrimSpace(parts[1])
// 		keyValues = append(keyValues, KeyValue{Key: key, Value: value})
// 	}
// 	if err := scanner.Err(); err != nil {
// 		log.Fatalf("Error reading file: %v", err)
// 	}
// 	return keyValues
// }
