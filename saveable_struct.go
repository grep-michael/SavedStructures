package savedstructures

import (
	"encoding/json"
	//"log"
	"fmt"
	//"os"
	"sync"
)

type Saveable struct {
	//do not export any variables here otherwise we could end up with infinite recursion while trying to
	mu       sync.RWMutex
	filepath string
	loader   *JsonLoader
	data     interface{}
}

func NewSaveable() *Saveable {
	return &Saveable{}
}

func (p *Saveable) InitSaveable(loader *JsonLoader, data interface{}) {
	p.loader = loader
	p.data = data
}

func (p *Saveable) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	jsonData, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return err
	}
	return p.loader.Save(jsonData)
	//return os.WriteFile(p.filepath, jsonData, 0644)
}
func (p *Saveable) SetBackupPath(path string) {
	p.loader.BackupPath = path
}
func (p *Saveable) SetBearerHeader(token string) {
	bearer := fmt.Sprintf("Bearer %s", token)
	p.loader.Headers.Set("Authorization", bearer)
}

func (p *Saveable) UsePostForUpdate() {
	p.loader.updateMethod = "POST"
}
func (p *Saveable) UsePutForUpdate() {
	p.loader.updateMethod = "PUT"
}
func (p *Saveable) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.loader.Load(p.data)
	return err
	/*
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

		return json.Unmarshal(file, p.data)*/
}
