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

type KeyValue struct {
	Key   string      `json:"key" yaml:"key"`
	Value interface{} `json:"value" yaml:"value"`
}

func (O *FileParser) Parse() []KeyValue {
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

func (O *FileParser) ParseYAML() []KeyValue {
	// Unmarshal the YAML data into a map
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(O.FileData, &yamlData); err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	return O.extractKeyValuePairs(yamlData)
}

func (O *FileParser) ParseJSON() []KeyValue {
	// Unmarshal the JSON data into a map
	var jsonData map[string]interface{}
	if err := json.Unmarshal(O.FileData, &jsonData); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return O.extractKeyValuePairs(jsonData)
}

func (O *FileParser) ParsePlain() []KeyValue {
	scanner := bufio.NewScanner(strings.NewReader(string(O.FileData)))
	return O.extractKeyValuePairsPlain(scanner)
}

// extractKeyValuePairs extracts key-value pairs recursively from a YAML map.
func (O *FileParser) extractKeyValuePairs(data map[string]interface{}) []KeyValue {
	var keyValues []KeyValue
	for key, value := range data {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			// Handle nested maps recursively
			nestedData := make(map[string]interface{})
			for k, val := range v {
				nestedData[fmt.Sprintf("%v", k)] = val
			}
			nestedKeyValues := O.extractKeyValuePairs(nestedData)
			keyValues = append(keyValues, nestedKeyValues...)
		case []interface{}:
			// Handle nested arrays recursively
			for i, item := range v {
				nestedData := map[string]interface{}{fmt.Sprintf("%d", i): item}
				nestedKeyValues := O.extractKeyValuePairs(nestedData)
				keyValues = append(keyValues, nestedKeyValues...)
			}
		default:
			// Add non-nested key-value pairs
			keyValues = append(keyValues, KeyValue{Key: key, Value: value})
		}
	}
	return keyValues
}

// extractKeyValuePairs extracts key-value pairs from a scanner reading a plain file.
func (O *FileParser) extractKeyValuePairsPlain(scanner *bufio.Scanner) []KeyValue {
	var keyValues []KeyValue
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			// Skip empty lines or lines starting with #
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Invalid line: %s", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		keyValues = append(keyValues, KeyValue{Key: key, Value: value})
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return keyValues
}
