package testdata

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type Corp struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

func (c *Corp) Set(i interface{}) {
	v, ok := i.(string)
	if !ok {
		return
	}
	c.Name = v
}

func (c *Corp) Less(kv interface{}) bool {
	switch v := kv.(type) {
	case Corp:
		if c.ID <= v.ID {
			return true
		}
	case *Corp:
		if c.ID <= (*v).ID {
			return true
		}
	case int:
		if c.ID <= v {
			return true
		}
	case *int:
		if c.ID <= *v {
			return true
		}
	}
	return false
}

func (c *Corp) Find(key interface{}) bool {
	switch v := key.(type) {
	// make no sense, used for test purpose only
	case Corp:
		if c.ID == v.ID {
			return true
		}
	// comment same as above
	case *Corp:
		if c.ID == v.ID {
			return true
		}
	case int:
		if c.ID == v {
			return true
		}
	case *int:
		if c.ID == *v {
			return true
		}
	}
	return false
}

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
