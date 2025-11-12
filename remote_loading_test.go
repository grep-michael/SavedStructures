package savedstructures

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type TestRemoteLoadingStruct struct {
	Saveable
	Created string     `json:"created"`
	Data    []DataItem `json:"data"`
}
type DataItem struct {
	Data  string   `json:"data,omitempty"`
	Data2 []string `json:"data2,omitempty"`
}

var TestData1 = []DataItem{
	{
		Data:  "Dell OptiPlex 7090",
		Data2: []string{"Intel i7-10700", "16GB RAM", "512GB SSD"},
	},
	{
		Data:  "HP EliteBook 840",
		Data2: []string{"Intel i5-10310U", "8GB RAM", "256GB SSD", "Windows 10"},
	},
	{
		Data:  "Lenovo ThinkPad T14",
		Data2: nil, // Omitted field
	},
}

func TestRemoteLoading(t *testing.T) {
	loader := NewJsonLoader(REMOTE, "http://127.0.0.1:8000/file")
	test := &TestRemoteLoadingStruct{
		Data: TestData1,
	}
	test.InitSaveable(loader, test)

	// Test saving to a remote location
	err := test.Save()
	fmt.Println(err)

	err = testJson("data.0.data2.1", "16GB RAM")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	test.Data[0].Data2[1] = "Test Change"

	err = test.Save()
	fmt.Println(err)

	err = testJson("data.0.data2.1", "Test Change")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	test.UsePostForUpdate()

	test.Data[0].Data2[1] = "Test Change via post"

	err = test.Save()
	fmt.Println(err)

	err = testJson("data.0.data2.1", "Test Change via post")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	//test loading from a remote location
	SetJSONValue("Testing/data.json", "data.0.data2.1", "Changed Manually")
	time.Sleep(1)
	test.Load()

	if test.Data[0].Data2[1] != "Changed Manually" {
		t.Errorf("%v\n", err)
	}
}
func TestRemoteFallBacks(t *testing.T) {
	loader := NewJsonLoader(REMOTE, "http://127.0.0.1:8000/notrealfile")
	test := &TestRemoteLoadingStruct{
		Data: TestData1,
	}
	test.InitSaveable(loader, test)
	test.SetBackupPath("./Testing/fallback.json")
	// Test saving to a remote location that isnt real
	err := test.Save()
	fmt.Println(err)
	/*
		err = testJson("data.0.data2.1", "16GB RAM")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		test.Data[0].Data2[1] = "Test Change"

		err = test.Save()
		fmt.Println(err)

		err = testJson("data.0.data2.1", "Test Change")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		test.UsePostForUpdate()

		test.Data[0].Data2[1] = "Test Change via post"

		err = test.Save()
		fmt.Println(err)

		err = testJson("data.0.data2.1", "Test Change via post")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		//test loading from a remote location
		SetJSONValue("Testing/data.json", "data.0.data2.1", "Changed Manually")
		time.Sleep(1)
		test.Load()
		marshaled, err := json.MarshalIndent(test, "", "  ") // Indent with 2 spaces
		fmt.Println(string(marshaled))*/

}

func testJson(path string, expected interface{}) error {
	v, err := GetJSONValue("Testing/data.json", path)
	if err != nil {
		return err
	}

	// Normalize JSON numbers (float64) so comparisons work
	switch expectedTyped := expected.(type) {

	case int:
		// JSON numbers come as float64
		jsonNum, ok := v.(float64)
		if !ok || int(jsonNum) != expectedTyped {
			return fmt.Errorf("expected int %d, got %#v", expectedTyped, v)
		}

	case float64:
		jsonNum, ok := v.(float64)
		if !ok || jsonNum != expectedTyped {
			return fmt.Errorf("expected float %f, got %#v", expectedTyped, v)
		}

	case string:
		jsonStr, ok := v.(string)
		if !ok || jsonStr != expectedTyped {
			return fmt.Errorf("expected string %q, got %#v", expectedTyped, v)
		}

	case bool:
		jsonBool, ok := v.(bool)
		if !ok || jsonBool != expectedTyped {
			return fmt.Errorf("expected bool %v, got %#v", expectedTyped, v)
		}

	default:
		return fmt.Errorf("unsupported expected type %T", expected)
	}

	return nil
}

func GetJSONValue(filename, path string) (interface{}, error) {
	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal into generic structure
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	// Traverse using dot notation
	tokens := strings.Split(path, ".")
	current := jsonData

	for _, token := range tokens {
		switch node := current.(type) {
		case map[string]interface{}:
			// object key
			v, ok := node[token]
			if !ok {
				return nil, fmt.Errorf("key '%s' not found", token)
			}
			current = v

		case []interface{}:
			// array index
			idx, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("expected index but got '%s'", token)
			}
			if idx < 0 || idx >= len(node) {
				return nil, fmt.Errorf("index out of range: %d", idx)
			}
			current = node[idx]

		default:
			return nil, fmt.Errorf("cannot descend into non-object/non-array at '%s'", token)
		}
	}

	return current, nil
}

func SetJSONValue(filename, path string, value interface{}) error {
	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal into interface{}
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	// Ensure top-level is a map
	root, ok := jsonData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("top-level JSON is not an object")
	}

	tokens := strings.Split(path, ".")
	current := interface{}(root)

	// Traverse, creating objects as needed
	for i, token := range tokens {
		last := i == len(tokens)-1

		switch node := current.(type) {

		case map[string]interface{}:
			// Last step -> set value
			if last {
				node[token] = value
				break
			}

			// Intermediate: ensure key exists
			next, exists := node[token]
			if !exists {
				// Create new object by default
				newObj := make(map[string]interface{})
				node[token] = newObj
				next = newObj
			}

			current = next

		case []interface{}:
			idx, err := strconv.Atoi(token)
			if err != nil {
				return fmt.Errorf("expected array index but got '%s'", token)
			}

			if idx < 0 || idx >= len(node) {
				return fmt.Errorf("index out of range %d", idx)
			}

			if last {
				node[idx] = value
				break
			}

			current = node[idx]

		default:
			return fmt.Errorf("cannot descend into %T at '%s'", node, token)
		}
	}

	// Marshal back to JSON
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filename, out, 0644)
}
