package tests

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Rhaqim/creds/internal/lib"
)

func TestFileParser_extractKeyValuePairs(t *testing.T) {
	parser := lib.FileParser{}

	// Test case 1: Nested map
	data1 := map[string]interface{}{
		"key1": "value1",
		"key2": map[interface{}]interface{}{
			"nestedKey1": "nestedValue1",
			"nestedKey2": "nestedValue2",
		},
	}
	expected1 := []lib.KeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "nestedKey1", Value: "nestedValue1"},
		{Key: "nestedKey2", Value: "nestedValue2"},
	}
	result1 := parser.ExtractKeyValuePairs(data1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected: %v, got: %v", expected1, result1)
	}

	// Test case 2: Nested array
	data2 := map[string]interface{}{
		"key1": "value1",
		"key2": []interface{}{
			"item1",
			"item2",
		},
	}
	expected2 := []lib.KeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "0", Value: "item1"},
		{Key: "1", Value: "item2"},
	}
	result2 := parser.ExtractKeyValuePairs(data2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed. Expected: %v, got: %v", expected2, result2)
	}

	// Test case 3: Non-nested key-value pairs
	data3 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	expected3 := []lib.KeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	result3 := parser.ExtractKeyValuePairs(data3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected: %v, got: %v", expected3, result3)
	}
}
func TestFileParser_Parse(t *testing.T) {
	// read file from the file system data/file/test.yaml and data/file/test.json and data/file/test.txt
	// Get the current directory of the parse_file_test.go file
	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}

	// Construct the absolute path to the test.yaml file
	yamlFilePath := filepath.Join(currentDir, "..", "data", "file")

	data, err := os.ReadFile(filepath.Join(yamlFilePath, "test.yaml"))
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	// Test case 1: YAML file format
	parser := lib.FileParser{
		FileFormat: "yaml",
		FileData:   data,
	}

	// Test case 1: YAML file format
	expected1 := []lib.KeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	result1 := parser.Parse()
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected: %v, got: %v", expected1, result1)
	}

	data, err = os.ReadFile(filepath.Join(yamlFilePath, "test.json"))
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	parser.FileFormat = "json"
	parser.FileData = data

	// Test case 2: JSON file format
	expected2 := []lib.KeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	result2 := parser.Parse()
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed. Expected: %v, got: %v", expected2, result2)
	}

	data, err = os.ReadFile(filepath.Join(yamlFilePath, "test.txt"))
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	parser.FileFormat = "plain"
	parser.FileData = data

	// Test case 3: Plain file format
	expected3 := []lib.KeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	result3 := parser.Parse()
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected: %v, got: %v", expected3, result3)
	}

	parser.FileFormat = "invalid"

	// Test case 4: Unsupported file format
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Test case 4 failed. Expected panic, but no panic occurred")
		}
	}()
	parser.Parse()
}
