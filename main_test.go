package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

// Test ReadFile function.
func TestReadFile(t *testing.T) {
	// Create a temporary file for testing
	tmpDir, err := os.MkdirTemp("", "test_read_file")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello, World!\nThis is a test file."
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test successful read
	input := ReadFileInput{Path: testFile}
	inputJSON, _ := json.Marshal(input)

	result, err := ReadFile(inputJSON)
	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
	}
	if result != testContent {
		t.Errorf("Expected %q, got %q", testContent, result)
	}

	// Test non-existent file
	input = ReadFileInput{Path: "non_existent_file.txt"}
	inputJSON, _ = json.Marshal(input)

	_, err = ReadFile(inputJSON)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

// Test ListFiles function.
func TestListFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "test_list_files")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files and directories
	testFiles := []string{"file1.txt", "file2.go", "subdir/file3.txt"}
	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}

	// Test listing files in the temp directory
	input := ListFilesInput{Path: tmpDir}
	inputJSON, _ := json.Marshal(input)

	result, err := ListFiles(inputJSON)
	if err != nil {
		t.Errorf("ListFiles failed: %v", err)
	}

	var files []string
	err = json.Unmarshal([]byte(result), &files)
	if err != nil {
		t.Errorf("Failed to unmarshal result: %v", err)
	}

	expectedFiles := []string{"file1.txt", "file2.go", "subdir/", "subdir/file3.txt"}
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
	}

	for _, expected := range expectedFiles {
		if !slices.Contains(files, expected) {
			t.Errorf("Expected file %q not found in result", expected)
		}
	}

	// Test listing files in current directory (empty path)
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	input = ListFilesInput{Path: ""}
	inputJSON, _ = json.Marshal(input)

	result, err = ListFiles(inputJSON)
	if err != nil {
		t.Errorf("ListFiles with empty path failed: %v", err)
	}

	err = json.Unmarshal([]byte(result), &files)
	if err != nil {
		t.Errorf("Failed to unmarshal result: %v", err)
	}

	if len(files) == 0 {
		t.Error("Expected files in current directory, got empty list")
	}
}

// Test EditFile function.
func TestEditFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_edit_file")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	originalContent := "Hello, World!\nThis is a test."
	err = os.WriteFile(testFile, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test successful edit
	input := EditFileInput{
		Path:   testFile,
		OldStr: "World",
		NewStr: "Go",
	}
	inputJSON, _ := json.Marshal(input)

	result, err := EditFile(inputJSON)
	if err != nil {
		t.Errorf("EditFile failed: %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %q", result)
	}

	// Verify the file was actually edited
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read edited file: %v", err)
	}
	expectedContent := "Hello, Go!\nThis is a test."
	if string(content) != expectedContent {
		t.Errorf("Expected %q, got %q", expectedContent, string(content))
	}

	// Test creating new file
	newFile := filepath.Join(tmpDir, "new.txt")
	input = EditFileInput{
		Path:   newFile,
		OldStr: "",
		NewStr: "New file content",
	}
	inputJSON, _ = json.Marshal(input)

	result, err = EditFile(inputJSON)
	if err != nil {
		t.Errorf("EditFile for new file failed: %v", err)
	}
	if !strings.Contains(result, "Successfully created") {
		t.Errorf("Expected creation message, got %q", result)
	}

	// Verify new file was created
	content, err = os.ReadFile(newFile)
	if err != nil {
		t.Errorf("Failed to read new file: %v", err)
	}
	if string(content) != "New file content" {
		t.Errorf("Expected 'New file content', got %q", string(content))
	}

	// Test error cases
	input = EditFileInput{
		Path:   testFile,
		OldStr: "nonexistent",
		NewStr: "replacement",
	}
	inputJSON, _ = json.Marshal(input)

	_, err = EditFile(inputJSON)
	if err == nil {
		t.Error("Expected error for non-existent old_str, got nil")
	}

	// Test invalid input (same old_str and new_str)
	input = EditFileInput{
		Path:   testFile,
		OldStr: "same",
		NewStr: "same",
	}
	inputJSON, _ = json.Marshal(input)

	_, err = EditFile(inputJSON)
	if err == nil {
		t.Error("Expected error for same old_str and new_str, got nil")
	}
}

// Test GenerateSchema function.
func TestGenerateSchema(t *testing.T) {
	schema := GenerateSchema[ReadFileInput]()

	if schema.Properties == nil {
		t.Error("Expected properties to be set, got nil")
	}

	// Basic validation that schema was generated
	// Note: Properties is interface{} type, so we can't directly index it
	// This test just ensures the schema generation doesn't panic
}

// Test ToolDefinition structure.
func TestToolDefinitions(t *testing.T) {
	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition, EditFileDefinition}

	expectedNames := []string{"read_file", "list_files", "edit_file"}

	if len(tools) != len(expectedNames) {
		t.Errorf("Expected %d tools, got %d", len(expectedNames), len(tools))
	}

	for i, tool := range tools {
		if tool.Name != expectedNames[i] {
			t.Errorf("Expected tool name %q, got %q", expectedNames[i], tool.Name)
		}

		if tool.Description == "" {
			t.Errorf("Tool %q has empty description", tool.Name)
		}

		if tool.Function == nil {
			t.Errorf("Tool %q has nil function", tool.Name)
		}

		if tool.InputSchema.Properties == nil {
			t.Errorf("Tool %q has nil input schema properties", tool.Name)
		}
	}
}

// Test Agent creation.
func TestNewAgent(t *testing.T) {
	// Create a real client for testing (it won't be used for actual API calls)
	client := anthropic.NewClient()
	getUserMessage := func() (string, bool) { return "test", true }
	tools := []ToolDefinition{ReadFileDefinition}

	agent := NewAgent(&client, getUserMessage, tools)

	if agent == nil {
		t.Fatal("Expected agent to be created, got nil")
	}

	if agent.client == nil {
		t.Error("Agent client not set")
	}

	if len(agent.tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(agent.tools))
	}

	if agent.tools[0].Name != "read_file" {
		t.Errorf("Expected first tool to be 'read_file', got %q", agent.tools[0].Name)
	}
}
