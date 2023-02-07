package dotenv

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestLoadFileInvalid(t *testing.T) {
	err := LoadFile(".invalid-env-file")

	if err == nil {
		t.Fatal("should return error when loading invalid .env file")
	}
}

func TestLoadFileValid(t *testing.T) {
	err := LoadFile(".")

	if err != nil {
		t.Error("shoud not return error when loading a valid .env file")
	}
}

func TestSetEnvVarsReaderNil(t *testing.T) {
	err := setEnvVars(nil)

	if err == nil {
		t.Fatal("should return error when nil reader is passed")
	}
}

func TestSetEnvVarsReaderOK(t *testing.T) {
	err := setEnvVars(strings.NewReader(""))

	if err != nil {
		t.Error("should not return error when a valid reader is passed")
	}
}

func TestSetEnvVarsError(t *testing.T) {
	envVars := []struct {
		Key   string
		Op    string
		Value string
	}{
		{` VAR1`, "", ` var1`},
		{`"" `, "=", "var2 "},
		{`VAR3	`, "", "var3	"},
		{"", "=", " "},
	}

	for _, ev := range envVars {
		line := fmt.Sprintf("%s%s%s", ev.Key, ev.Op, ev.Value)
		err := setEnvVars(strings.NewReader(line))
		if err == nil {
			t.Errorf("should not return nil error for invalid line: %s", line)
		}
	}

	for _, ev := range envVars {
		value := os.Getenv(cleanString(ev.Key))
		if value != "" {
			t.Errorf("expected empty string value for envrironment variable %s, got %s", cleanString(ev.Key), value)
		}
	}
}

func TestSetEnvVarsOK(t *testing.T) {
	envVars := []struct {
		Key   string
		Op    string
		Value string
	}{
		{` VAR1`, "=", ` var1`},
		{`"VAR2" `, "=", "var2 "},
		{`VAR3	`, "=", "var3	"},
		{"dsn", "=", "host=localhost user=user"},
		{"a", "=", ""},
		{"# a comment", "", ""},
	}

	for _, ev := range envVars {
		line := fmt.Sprintf("%s%s%s", ev.Key, ev.Op, ev.Value)
		err := setEnvVars(strings.NewReader(line))
		if err != nil {
			t.Errorf("should not return error when parsing valid line: %v", err)
		}
	}

	for _, ev := range envVars {
		value := os.Getenv(cleanString(ev.Key))
		if value != cleanString(ev.Value) {
			t.Errorf("expected value %s for envrironment variable %s, got %s", cleanString(ev.Value), cleanString(ev.Key), value)
		}
	}

	for _, ev := range envVars {
		err := os.Unsetenv(cleanString(ev.Key))
		if err != nil {
			t.Errorf("error when unsetting envrironment variable %s", cleanString(ev.Key))
		}
	}
}
