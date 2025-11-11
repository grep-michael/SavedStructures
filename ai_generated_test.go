package savedstructures

/*
entierly ai generated testing, human reviewed of course but i am stupid :p


import (
	"encoding/json"
	"os"
	"testing"
)

type TestPersistentStruct struct {
	Saveable
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
	Field3 bool   `json:"field3"`
}

type TestEncryptedStruct struct {
	EncryptedSaveable
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
	Field3 bool   `json:"field3"`
}

// Helper to create TestPersistentStruct with defaults
func newTestPersistent(filepath string, field1 string, field2 int, field3 bool) *TestPersistentStruct {
	s := &TestPersistentStruct{}
	s.Field1 = field1
	s.Field2 = field2
	s.Field3 = field3
	s.InitSaveable(filepath, s)
	return s
}

// Helper to create TestEncryptedStruct with defaults
func newTestEncrypted(filepath string, key []byte, field1 string, field2 int, field3 bool) *TestEncryptedStruct {
	s := &TestEncryptedStruct{}
	s.EncryptedSaveable = *NewEncryptedSaveable(key)
	s.Field1 = field1
	s.Field2 = field2
	s.Field3 = field3
	s.InitPersistent(filepath, s)
	return s
}

// ============ Persistent Tests ============

func TestPersistentCreateWithDefaults(t *testing.T) {
	filepath := genTmpPath("persistent_defaults.json")

	s := newTestPersistent(filepath, "default", 42, true)

	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// Verify content
	if s.Field1 != "default" {
		t.Errorf("Field1 differs: expected 'default', got '%s'", s.Field1)
	}
	if s.Field2 != 42 {
		t.Errorf("Field2 differs: expected 42, got %d", s.Field2)
	}
	if s.Field3 != true {
		t.Errorf("Field3 differs: expected true, got %t", s.Field3)
	}
}

func TestPersistentLoadExisting(t *testing.T) {
	filepath := genTmpPath("persistent_existing.json")

	// Create initial data
	initial := newTestPersistent(filepath, "initial", 100, false)
	if err := initial.Save(); err != nil {
		t.Fatalf("Failed to save initial data: %v", err)
	}

	// Load in new instance
	loaded := newTestPersistent(filepath, "wrong", 0, true)
	if err := loaded.Load(); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	// Verify loaded values match saved, not defaults
	if loaded.Field1 != "initial" {
		t.Errorf("Field1 differs: expected 'initial', got '%s'", loaded.Field1)
	}
	if loaded.Field2 != 100 {
		t.Errorf("Field2 differs: expected 100, got %d", loaded.Field2)
	}
	if loaded.Field3 != false {
		t.Errorf("Field3 differs: expected false, got %t", loaded.Field3)
	}
}

func TestPersistentLoadNonexistent(t *testing.T) {
	filepath := genTmpPath("persistent_nonexistent.json")

	s := newTestPersistent(filepath, "default", 99, true)

	// Load should create file with defaults
	if err := s.Load(); err != nil {
		t.Fatalf("Failed to load nonexistent file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Error("File was not created during load")
	}

	// Verify defaults were preserved
	if s.Field1 != "default" {
		t.Errorf("Field1 differs: expected 'default', got '%s'", s.Field1)
	}
}

func TestPersistentUpdate(t *testing.T) {
	filepath := genTmpPath("persistent_update.json")

	s := newTestPersistent(filepath, "before", 1, false)
	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Update with auto-save
	s.Field1 = "after"
	s.Field2 = 2
	s.Field3 = true
	s.Save()
	// Load in new instance to verify persistence
	loaded := newTestPersistent(filepath, "", 0, false)
	if err := loaded.Load(); err != nil {
		t.Fatalf("Failed to load after update: %v", err)
	}

	if loaded.Field1 != "after" {
		t.Errorf("Field1 not updated: expected 'after', got '%s'", loaded.Field1)
	}
	if loaded.Field2 != 2 {
		t.Errorf("Field2 not updated: expected 2, got %d", loaded.Field2)
	}
	if loaded.Field3 != true {
		t.Errorf("Field3 not updated: expected true, got %t", loaded.Field3)
	}
}

func TestPersistentJSONFormat(t *testing.T) {
	filepath := genTmpPath("persistent_json.json")

	s := newTestPersistent(filepath, "test", 123, true)
	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Read raw file and verify it's valid JSON
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("File is not valid JSON: %v", err)
	}

	// Verify Persistent fields are NOT in JSON
	if _, exists := parsed["mu"]; exists {
		t.Error("Unexported field 'mu' should not be in JSON")
	}
	if _, exists := parsed["filepath"]; exists {
		t.Error("Unexported field 'filepath' should not be in JSON")
	}
	if _, exists := parsed["data"]; exists {
		t.Error("Unexported field 'data' should not be in JSON")
	}
}

// ============ EncryptedPersistent Tests ============

func TestEncryptedCreateWithDefaults(t *testing.T) {
	filepath := genTmpPath("encrypted_defaults.json")

	key := []byte("12345678901234567890123456789012")
	s := newTestEncrypted(filepath, key, "secret", 42, true)

	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save encrypted: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// Verify content
	if s.Field1 != "secret" {
		t.Errorf("Field1 differs: expected 'secret', got '%s'", s.Field1)
	}
}

func TestEncryptedLoadExisting(t *testing.T) {
	filepath := genTmpPath("encrypted_existing.json")

	key := []byte("12345678901234567890123456789012")

	// Create and save
	initial := newTestEncrypted(filepath, key, "encrypted_data", 999, false)
	if err := initial.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Load in new instance
	loaded := newTestEncrypted(filepath, key, "wrong", 0, true)
	if err := loaded.Load(); err != nil {
		t.Fatalf("Failed to load encrypted: %v", err)
	}

	// Verify decryption worked
	if loaded.Field1 != "encrypted_data" {
		t.Errorf("Field1 differs: expected 'encrypted_data', got '%s'", loaded.Field1)
	}
	if loaded.Field2 != 999 {
		t.Errorf("Field2 differs: expected 999, got %d", loaded.Field2)
	}
}

func TestEncryptedNotReadable(t *testing.T) {
	filepath := genTmpPath("encrypted_not_readable.json")

	key := []byte("12345678901234567890123456789012")
	s := newTestEncrypted(filepath, key, "secret_data", 123, true)

	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Read raw file - should NOT be valid JSON
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err == nil {
		t.Error("Encrypted file should NOT be valid JSON")
	}

	// Should not contain plaintext
	if string(data) == "secret_data" {
		t.Error("Encrypted file contains plaintext data")
	}
}

func TestEncryptedWrongKey(t *testing.T) {
	filepath := genTmpPath("encrypted_wrong_key.json")

	key1 := []byte("12345678901234567890123456789012")
	key2 := []byte("99999999999999999999999999999999")

	// Save with key1
	s := newTestEncrypted(filepath, key1, "secret", 42, true)
	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Try to load with key2
	loaded := newTestEncrypted(filepath, key2, "", 0, false)
	if err := loaded.Load(); err == nil {
		t.Error("Should fail to decrypt with wrong key")
	}
}

func TestEncryptedUpdate(t *testing.T) {
	filepath := genTmpPath("encrypted_update.json")

	key := []byte("12345678901234567890123456789012")
	s := newTestEncrypted(filepath, key, "before", 1, false)
	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Update with auto-save
	s.Field1 = "after"
	s.Field2 = 2
	s.Field3 = true
	s.Save()

	// Load and verify
	loaded := newTestEncrypted(filepath, key, "", 0, false)
	if err := loaded.Load(); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if loaded.Field1 != "after" {
		t.Errorf("Field1 not updated: expected 'after', got '%s'", loaded.Field1)
	}
	if loaded.Field2 != 2 {
		t.Errorf("Field2 not updated: expected 2, got %d", loaded.Field2)
	}
}

func TestEncryptedLoadNonexistent(t *testing.T) {
	filepath := genTmpPath("encrypted_nonexistent.json")

	key := []byte("12345678901234567890123456789012")
	s := newTestEncrypted(filepath, key, "default", 77, false)

	// Load should create encrypted file with defaults
	if err := s.Load(); err != nil {
		t.Fatalf("Failed to load nonexistent encrypted file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Error("File was not created during load")
	}

	// Verify defaults were preserved
	if s.Field1 != "default" {
		t.Errorf("Field1 differs: expected 'default', got '%s'", s.Field1)
	}
}
*/
