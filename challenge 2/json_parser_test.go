package parser_test

import (
	"os"
	"parser"
	"testing"
)

func TestParseStep1(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step1/valid.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err != nil {
			t.Error(err)

		}
		t.Log("{} is a valid json")
	})

	t.Run("invalid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step1/invalid.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err == nil {
			t.Error("invalid json is not detected")
		}
		t.Log("empty is not a valid json")
	})
}