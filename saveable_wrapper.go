package savedstructures

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type SaveableWrapper[T any] struct {
	mu       sync.RWMutex
	Data     T
	filepath string
}

func NewSaveableWrapper[T any](filepath string, defaults ...T) (*SaveableWrapper[T], error) {
	p := &SaveableWrapper[T]{
		filepath: filepath,
	}

	if len(defaults) > 0 {
		p.Data = defaults[0]
	}

	load_err := p.Load()
	return p, load_err
}

// Load reads the JSON file into the struct
func (p *SaveableWrapper[T]) Load() error {
	file, err := os.ReadFile(p.filepath)
	if err == nil { //file found, marshel it
		return json.Unmarshal(file, &p.Data)
	}
	//sanity check
	if !os.IsNotExist(err) {
		log.Printf("Failed to read log file: %v\n", err)
		return err
	}
	//file not found, save it
	if saveErr := p.Save(); saveErr != nil {
		log.Printf("Failed save file during loading: %v\n", err)
		return saveErr
	}
	return nil
}

// Save writes the struct to the JSON file
func (p *SaveableWrapper[T]) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := json.MarshalIndent(p.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filepath, data, 0644)
}

// Update allows safe modification with auto-save
func (p *SaveableWrapper[T]) Update(fn func(*T)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	fn(&p.Data)

	data, err := json.MarshalIndent(p.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filepath, data, 0644)
}

// Read allows safe read-only access
func (p *SaveableWrapper[T]) Read(fn func(*T)) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	fn(&p.Data)
}
