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
	Type    LoaderType
	Path    string //either url or file path
	Headers http.Header
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
		req, err := http.NewRequest("PUT", ld.Path, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

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
