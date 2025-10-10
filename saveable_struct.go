package savedstructures

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Saveable struct {
	//do not export any variables here otherwise we could end up with infinite recursion while trying to
	mu       sync.RWMutex
	filepath string
	data     interface{}
}

func NewSaveable() *Saveable {
	return &Saveable{}
}

func (p *Saveable) InitSaveable(filepath string, data interface{}) {
	p.filepath = filepath
	p.data = data
}

func (p *Saveable) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	jsonData, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filepath, jsonData, 0644)
}

func (p *Saveable) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.ReadFile(p.filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to read file: %v\n", err)
			return err
		}

		// File not found, save defaults
		jsonData, err := json.MarshalIndent(p.data, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(p.filepath, jsonData, 0644); err != nil {
			log.Printf("Failed to save file during loading: %v\n", err)
			return err
		}
		return nil // Data already contains defaults
	}

	return json.Unmarshal(file, p.data)
}
