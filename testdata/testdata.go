package testdata

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

// Corp is a type used for testing
type Corp struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

// Set implements interface Setter in package container
func (c *Corp) Set(i interface{}) {
	v, ok := i.(string)
	if !ok {
		return
	}
	c.Name = v
}

// Less implements interface Lesser in package container
func (c *Corp) Less(kv interface{}) bool {
	switch v := kv.(type) {
	case Corp:
		return c.ID <= v.ID
	case *Corp:
		return c.ID <= (*v).ID
	case int:
		return c.ID <= v
	case *int:
		return c.ID <= *v
	}
	return false
}

// Find implements interface Finder in package container
func (c *Corp) Find(key interface{}) bool {
	switch v := key.(type) {
	// make no sense, used for test purpose only
	case Corp:
		return c.ID == v.ID
	// comment same as above
	case *Corp:
		return c.ID == v.ID
	case int:
		return c.ID == v
	case *int:
		return c.ID == *v
	}
	return false
}

// TestCases contain some test data loaded from testdata.json, which use lately by package to test
var TestCases []Corp

const testFile = "testdata.json"

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}

	data, err := ioutil.ReadFile(filepath.Join(filepath.Dir(file), testFile))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &TestCases); err != nil {
		panic(err)
	}
}
