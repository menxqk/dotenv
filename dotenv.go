package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func LoadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return setEnvVars(f)
}

func setEnvVars(r io.Reader) error {
	if r == nil {
		return errors.New("nil reader")
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, err := parseLine(line)
		if err != nil {
			return err
		}

		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseLine(s string) (string, string, error) {
	s = strings.TrimSpace(s)
	ss := strings.Split(s, "=")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("line does not have key/value assignment: %s", s)
	}

	ss[0] = cleanString(ss[0])
	ss[1] = cleanString(ss[1])

	return ss[0], ss[1], nil
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.TrimSpace(s)
	return s
}
