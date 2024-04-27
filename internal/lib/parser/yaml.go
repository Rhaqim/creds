package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseYaml() {
	// Read the YAML file
	data := ReadFile("data.yaml")

	// Unmarshal the YAML data into a map
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	// Extract key-value pairs
	keyValues := extractKeyValuePairs(yamlData)

	// Print the extracted key-value pairs
	for _, kv := range keyValues {
		fmt.Printf("Key: %s, Value: %v\n", kv.Key, kv.Value)
	}
}

func ReadFile(filePath string) []byte {
	// Read the YAML file
	var buffer bytes.Buffer
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(&buffer, file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return buffer.Bytes()
}
