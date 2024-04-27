package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ParsePlain() {
	// Open the .env file
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Parse the plain file content and extract key-value pairs
	keyValues := extractKeyValuePairsPlain(scanner)

	// Print the extracted key-value pairs
	for _, kv := range keyValues {
		fmt.Printf("Key: %s, Value: %s\n", kv.Key, kv.Value)
	}
}

// extractKeyValuePairs extracts key-value pairs from a scanner reading a plain file.
func extractKeyValuePairsPlain(scanner *bufio.Scanner) []KeyValue {
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
