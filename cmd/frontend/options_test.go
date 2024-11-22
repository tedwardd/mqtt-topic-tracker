package main

import (
	"testing"
)

// All test scenarios that should fail
var failPaths = []serverOpts{
	{dbFile: "../../testdata/topics.db", tmplPath: "foo/bar/baz"},
	{dbFile: "foo/bar/baz", tmplPath: "../../web/templates"},
	{dbFile: "foo/bar/baz", tmplPath: "foo/bar/baz"},
	{dbFile: "../../testdata", tmplPath: "foo/bar/baz"},
	{dbFile: "foo/bar/baz", tmplPath: "../../testdata/topics.db"},
	{dbFile: "../../testdata/topics.db", tmplPath: "../../testdata/topics.db"},
	{dbFile: "../../testdata", tmplPath: "../../testdata"},
}

// All test scenarios that should pass
var validPaths = []serverOpts{
	{dbFile: "../../testdata/topics.db", tmplPath: "../../web/templates"},
}

// Invalid paths should return an error
func TestValidate(t *testing.T) {
	for k := range failPaths {
		err := failPaths[k].Validate()
		if err == nil {
			t.Error("invalid values produced nil from Validate()")
		}
	}

	for k := range validPaths {
		err := validPaths[k].Validate()
		if err != nil {
			t.Error("valid values produced non-nil from Validate()")
		}
	}
}
