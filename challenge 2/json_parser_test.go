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

func TestParseStep2(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step2/valid.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err != nil {
			t.Error(err)

		}
		t.Log(string(jsonString) + " is a valid json")
	})
	t.Run("valid 2", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step2/valid2.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err != nil {
			t.Error(err)

		}
		t.Log(string(jsonString) + " is a valid json")
	})

	t.Run("invalid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step2/invalid.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err == nil {
			t.Error("invalid json is not detected")
		}
		t.Log(string(jsonString) + " is not a valid json")
	})
	t.Run("invalid 2", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step2/invalid2.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err == nil {
			t.Error("invalid json is not detected")
		}
		t.Log(string(jsonString) + " is not a valid json")
	})
}

func TestParseStep3(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step3/valid.json")
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
		jsonString, err := os.ReadFile("tests/step3/invalid.json")
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

func TestParseSte4(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step4/valid.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err != nil {
			t.Error(err)

		}
		t.Log(string(jsonString) + " is a valid json")
	})
	t.Run("valid 2", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step4/valid2.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err != nil {
			t.Error(err)

		}
		t.Log(string(jsonString) + " is a valid json")
	})

	t.Run("invalid", func(t *testing.T) {
		jsonString, err := os.ReadFile("tests/step4/invalid.json")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = parser.Parse(string(jsonString))
		if err == nil {
			t.Error("invalid json is not detected")
		}
		t.Log(string(jsonString) + " is not a valid json")
	})
}
