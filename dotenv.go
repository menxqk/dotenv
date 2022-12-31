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
	idx := strings.Index(s, "=")
	if idx == -1 || idx == 0 {
		return "", "", fmt.Errorf("line does not have key/value assignment: %s", s)
	}

	ss := []string{}
	ss = append(ss, s[0:idx])
	idx = idx + 1
	if idx < len(s) {
		ss = append(ss, s[idx:])
	} else {
		ss = append(ss, "")
	}

	ss[0] = cleanString(ss[0])
	ss[1] = cleanString(ss[1])

	return ss[0], ss[1], nil
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\"", "")
	return s
}
