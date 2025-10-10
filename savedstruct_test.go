package savedstructures

import (
	"testing"
)

type TestStruct struct {
	Field1 string `json:"Field1"`
	Field2 int    `json:"Field2"`
	Field3 bool   `json:"Field3"`
}

/*
testing wrapper
*/
func TestCreatingNonexistantWithDefault(t *testing.T) {
	defaults := TestStruct{}
	defaults.Field1 = "DefaultString"
	defaults.Field2 = 10
	defaults.Field3 = true

	wrapper, err := NewSaveableWrapper[TestStruct](genTmpPath("Test1.json"), defaults)
	if err != nil {
		t.Logf("Error creating wrapper: %v\n", err)
	}

	if wrapper.Data.Field1 != defaults.Field1 {
		t.Errorf("Field1 differs: %s (e) : %s (a)\n", wrapper.Data.Field1, defaults.Field1)
	}
	if wrapper.Data.Field2 != defaults.Field2 {
		t.Errorf("Field2 differs: %d (e) : %d (a)\n", wrapper.Data.Field2, defaults.Field2)
	}
	if wrapper.Data.Field3 != defaults.Field3 {
		t.Errorf("Field3 differs: %t (e) : %t  (a)\n", wrapper.Data.Field3, defaults.Field3)
	}
}

func TestLoadingPreExisting(t *testing.T) {
	defaults := TestStruct{}
	defaults.Field1 = "premadeString"
	defaults.Field2 = -1
	defaults.Field3 = true //all default values start at true

	wrapper, err := NewSaveableWrapper[TestStruct](genTmpPath("premade.json"))
	if err != nil {
		t.Logf("Error creating wrapper: %v\n", err)
	}

	if wrapper.Data.Field1 != defaults.Field1 {
		t.Errorf("Field1 differs from premade: %s (e) : %s (a)\n", wrapper.Data.Field1, defaults.Field1)
	}
	if wrapper.Data.Field2 != defaults.Field2 {
		t.Errorf("Field2 differs from premade: %d (e) : %d (a)\n", wrapper.Data.Field2, defaults.Field2)
	}
	if wrapper.Data.Field3 != defaults.Field3 {
		t.Errorf("Field3 differs from premade: %t (e) : %t  (a)\n", wrapper.Data.Field3, defaults.Field3)
	}

}

func TestCreatingNonexistantWithoutDefaults(t *testing.T) {

	wrapper, err := NewSaveableWrapper[TestStruct](genTmpPath("Test2.json"))
	if err != nil {
		t.Logf("Error creating wrapper: %v\n", err)
	}

	var intT int
	var boolT bool
	var stringT string

	if wrapper.Data.Field1 != stringT {
		t.Errorf("Field1 differs from golang default: %s (e) : %s (a)\n", wrapper.Data.Field1, stringT)
	}
	if wrapper.Data.Field2 != intT {
		t.Errorf("Field2 differs from golang default: %d (e) : %d (a)\n", wrapper.Data.Field2, intT)
	}
	if wrapper.Data.Field3 != boolT {
		t.Errorf("Field3 differs from golang default: %t (e) : %t  (a)\n", wrapper.Data.Field3, boolT)
	}

}
