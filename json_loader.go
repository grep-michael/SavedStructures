package savedstructures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"net/url"
	"os"
)

type LoaderType int

const (
	REMOTE LoaderType = iota
	LOCAL  LoaderType = iota
)

func (lt LoaderType) Toggle() LoaderType {
	switch lt {
	case REMOTE:
		return LOCAL
	case LOCAL:
		return REMOTE
	default:
		return lt // Return unchanged for unknown values
	}
}

type JsonLoader struct {
	Type         LoaderType
	Path         string //either url or file path
	BackupPath   string //either url or file path, must be a differnt type than Path
	Headers      http.Header
	updateMethod string
}

func NewJsonLoader(typ LoaderType, path string) *JsonLoader {
	loader := &JsonLoader{
		Type: typ,
	}
	if typ == REMOTE {
		if _, err := url.ParseRequestURI(path); err != nil {
			log.Panicf("Tried to pass none url to Remote json loader Path:%s\n", path)
		}
		loader.Path = path
		loader.Headers = http.Header{
			"Content-Type": []string{"application/json"},
			"User-Agent":   []string{"Golang SaveableStructures"},
		}
		loader.updateMethod = "PUT"
	} else {
		loader.Path = path
	}
	return loader
}

func (ld *JsonLoader) withBackupRetry(operation func(LoaderType) error) error {
	err := operation(ld.Type)
	if err != nil {
		if ld.BackupPath == "" {
			return fmt.Errorf("Failed to perform operation using backup method, No BackupPath set")
		}

		// Swap backup path and current path
		ld.Path, ld.BackupPath = ld.BackupPath, ld.Path
		defer func() {
			ld.Path, ld.BackupPath = ld.BackupPath, ld.Path
		}()

		fmt.Println(err)
		err = operation(ld.Type.Toggle())
	}
	return err
}

// Simplified Load function
func (ld *JsonLoader) Load(data interface{}) error {
	return ld.withBackupRetry(func(loadType LoaderType) error {
		return ld.load(data, loadType)
	})
}

// Simplified Save function
func (ld *JsonLoader) Save(data []byte) error {
	return ld.withBackupRetry(func(loadType LoaderType) error {
		return ld.save(data, loadType)
	})
}

func (ld *JsonLoader) Save2(data []byte) error {
	err := ld.save(data, ld.Type)
	if err != nil {
		if ld.BackupPath == "" {
			return fmt.Errorf("Failed to save using backup method, No BackupPath set")
		}
		//swap backup path and current path
		*&ld.Path, *&ld.BackupPath = *&ld.BackupPath, *&ld.Path
		defer func() {
			*&ld.Path, *&ld.BackupPath = *&ld.BackupPath, *&ld.Path
		}()

		err = ld.save(data, ld.Type.Toggle())
	}
	return err
}
func (ld *JsonLoader) save(data []byte, loaderType LoaderType) error {
	if loaderType == LOCAL {
		return os.WriteFile(ld.Path, data, 0644)
	} else if loaderType == REMOTE {
		req, err := http.NewRequest(ld.updateMethod, ld.Path, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("error creating save request: %w", err)
		}
		req.Header = ld.Headers
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending save request: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	return fmt.Errorf("Failed to Save json, Loader type not recognized")
}
func (ld *JsonLoader) Load2(data interface{}) error {
	err := ld.load(data, ld.Type)
	if err != nil {
		if ld.BackupPath == "" {
			return fmt.Errorf("Failed to load using backup method, No BackupPath set")
		}
		//swap backup path and current path
		*&ld.Path, *&ld.BackupPath = *&ld.BackupPath, *&ld.Path
		defer func() {
			*&ld.Path, *&ld.BackupPath = *&ld.BackupPath, *&ld.Path
		}()

		fmt.Println(err)
		err = ld.load(data, ld.Type.Toggle())
	}
	return err
}
func (ld *JsonLoader) load(data interface{}, loaderType LoaderType) error {
	if loaderType == REMOTE {
		//build request
		req, err := http.NewRequest("GET", ld.Path, nil)
		if err != nil {
			return fmt.Errorf("error creating load request: %w", err)
		}
		req.Header = ld.Headers
		//make http client
		client := &http.Client{}
		//do request
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending load request: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		//load json
		return json.NewDecoder(resp.Body).Decode(data)
	} else if loaderType == LOCAL {
		file, err := os.ReadFile(ld.Path)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Failed to read file: %v\n", err)
				return err
			}

			// File not found, save defaults
			jsonData, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(ld.Path, jsonData, 0644); err != nil {
				log.Printf("Failed to save file during loading: %v\n", err)
				return err
			}
			return nil // Data already contains defaults
		}

		return json.Unmarshal(file, data)
	}

	return fmt.Errorf("Failed to load json, Loader type not recognized")
}
