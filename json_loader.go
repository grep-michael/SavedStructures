package savedstructures

import (
	"bytes"
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

type JsonLoader struct {
	Type         LoaderType
	Path         string //either url or file path
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

func (ld *JsonLoader) SetBearerHeader(token string) {
	bearer := fmt.Sprintf("Bearer %s", token)
	ld.Headers.Set("Authorization", bearer)
}
func (ld *JsonLoader) Save(data []byte) error {
	if ld.Type == LOCAL {
		return os.WriteFile(ld.Path, data, 0644)
	} else if ld.Type == REMOTE {
		req, err := http.NewRequest(ld.updateMethod, ld.Path, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}
		req.Header = ld.Headers
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	return fmt.Errorf("Failed to Save json, Loader type not recognized")
}

func (ld *JsonLoader) Load() error {
	req, err := http.NewRequest("GET", ld.Path, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header = ld.Headers

	client := &http.Client{
		Timeout: 30 * 60,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
